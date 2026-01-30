package telegram

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Connector implements the channels.Connector interface for Telegram
type Connector struct {
	mu       sync.RWMutex
	config   config.TelegramConfig
	bot      *tgbotapi.BotAPI
	userRepo repository.UserRepository
	logger   *slog.Logger
	running  bool
	incoming chan *channels.Message
	cancel   context.CancelFunc
	updates  <-chan tgbotapi.Update
}

// NewConnector creates a new Telegram connector
func NewConnector(
	cfg config.TelegramConfig,
	userRepo repository.UserRepository,
	logger *slog.Logger,
) *Connector {
	return &Connector{
		config:   cfg,
		userRepo: userRepo,
		logger:   logger,
		incoming: make(chan *channels.Message, 100),
	}
}

// Name returns the name of the channel
func (c *Connector) Name() string {
	return "telegram"
}

// Start initializes and starts the connector
func (c *Connector) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("telegram connector is already running")
	}

	// Create bot instance
	bot, err := tgbotapi.NewBotAPI(c.config.BotToken)
	if err != nil {
		return fmt.Errorf("failed to create telegram bot: %w", err)
	}

	c.bot = bot
	if c.logger != nil {
		c.logger.Info("Telegram bot initialized", "bot_username", bot.Self.UserName)
	}

	// Setup webhook or polling
	if c.config.WebhookURL != "" {
		if err := c.setupWebhook(); err != nil {
			return fmt.Errorf("failed to setup webhook: %w", err)
		}
	} else {
		c.setupPolling()
	}

	// Create context for cancellation
	ctx, cancel := context.WithCancel(ctx)
	c.cancel = cancel
	c.running = true

	// Start processing updates
	go c.processUpdates(ctx)

	if c.logger != nil {
		c.logger.Info("Telegram connector started", "mode", c.getMode())
	}
	return nil
}

// Stop gracefully stops the connector
func (c *Connector) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("telegram connector is not running")
	}

	// Cancel update processing
	if c.cancel != nil {
		c.cancel()
	}

	// Remove webhook if it was set
	if c.config.WebhookURL != "" {
		_, err := c.bot.Request(tgbotapi.DeleteWebhookConfig{})
		if err != nil {
			if c.logger != nil {
				c.logger.Warn("Failed to remove webhook", "error", err)
			}
		}
	}

	c.running = false
	close(c.incoming)

	if c.logger != nil {
		c.logger.Info("Telegram connector stopped")
	}
	return nil
}

// SendResponse sends a response to a user
func (c *Connector) SendResponse(ctx context.Context, userID string, response *channels.Response) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.running {
		return fmt.Errorf("telegram connector is not running")
	}

	// Parse chat ID from userID
	chatID, err := parseChatID(userID)
	if err != nil {
		return fmt.Errorf("invalid user ID format: %w", err)
	}

	// Create and send message
	msg := tgbotapi.NewMessage(chatID, response.Content)
	if _, err := c.bot.Send(msg); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	if c.logger != nil {
		c.logger.Debug("Message sent", "user_id", userID, "chat_id", chatID)
	}
	return nil
}

// Incoming returns a channel for incoming messages
func (c *Connector) Incoming() <-chan *channels.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.incoming
}

// IsRunning returns whether the connector is currently running
func (c *Connector) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// GetUser retrieves a user by channel-specific ID
func (c *Connector) GetUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	return c.userRepo.FindByChannel(ctx, "telegram", channelUserID)
}

// CreateUser creates a new user in the system
func (c *Connector) CreateUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	user := entity.NewUser("telegram", channelUserID)
	if err := c.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	if c.logger != nil {
		c.logger.Info("User created", "channel_id", channelUserID, "user_id", user.ID)
	}
	return user, nil
}

// setupWebhook configures the webhook for incoming updates
func (c *Connector) setupWebhook() error {
	webhook, err := tgbotapi.NewWebhook(c.config.WebhookURL)
	if err != nil {
		return err
	}

	_, err = c.bot.Request(webhook)
	if err != nil {
		return err
	}

	// Get webhook updates
	c.updates = c.bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	if c.logger != nil {
		c.logger.Info("Telegram webhook configured", "url", c.config.WebhookURL)
	}
	return nil
}

// setupPolling configures long polling for incoming updates
func (c *Connector) setupPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	c.updates = c.bot.GetUpdatesChan(u)
	if c.logger != nil {
		c.logger.Info("Telegram polling configured")
	}
}

// processUpdates processes incoming updates from Telegram
func (c *Connector) processUpdates(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if c.logger != nil {
				c.logger.Debug("Update processing stopped")
			}
			return
		case update, ok := <-c.updates:
			if !ok {
				if c.logger != nil {
					c.logger.Debug("Updates channel closed")
				}
				return
			}

			// Handle callback queries (inline buttons)
			if update.CallbackQuery != nil {
				c.handleCallbackQuery(update.CallbackQuery)
				continue
			}

			// Handle regular messages
			if update.Message == nil {
				continue
			}

			if !c.isAllowed(update.Message.Chat.ID, update.Message.From.ID) {
				if c.logger != nil {
					c.logger.Warn("Message from unauthorized user/chat ignored",
						"chat_id", update.Message.Chat.ID,
						"user_id", update.Message.From.ID,
						"username", update.Message.From.UserName,
					)
				}
				continue
			}

			c.handleMessage(ctx, update.Message)
		}
	}
}

// handleMessage processes an incoming message from Telegram
func (c *Connector) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	// Ensure user exists in the system
	channelUserID := fmt.Sprintf("%d", message.From.ID)
	user, err := c.GetUser(ctx, channelUserID)
	if err != nil {
		// User doesn't exist, create it
		if c.logger != nil {
			c.logger.Debug("User not found, creating new user", "channel_user_id", channelUserID)
		}
		user, err = c.CreateUser(ctx, channelUserID)
		if err != nil {
			if c.logger != nil {
				c.logger.Error("Failed to create user", "error", err)
			}
			return
		}
	}

	// Extract message content and metadata
	content, metadata := c.extractMessageContent(message)

	// Create message struct
	msg := &channels.Message{
		UserID:    formatUserID(message.From.ID, message.Chat.ID),
		ChannelID: formatChatID(message.Chat.ID),
		Content:   content,
		Metadata:  metadata,
	}

	// Add user ID to metadata
	if user != nil {
		msg.Metadata["user_internal_id"] = user.ID.String()
	}

	// Send message to incoming channel
	select {
	case c.incoming <- msg:
		if c.logger != nil {
			c.logger.Debug("Message received",
				"user_id", msg.UserID,
				"type", msg.Metadata["message_type"],
				"content_length", len(msg.Content),
			)
		}
	default:
		if c.logger != nil {
			c.logger.Warn("Incoming channel is full, message dropped")
		}
	}
}

// handleCallbackQuery processes callback queries from inline buttons
func (c *Connector) handleCallbackQuery(callback *tgbotapi.CallbackQuery) {
	// Create message struct from callback query
	msg := &channels.Message{
		UserID:    formatUserID(callback.From.ID, callback.Message.Chat.ID),
		ChannelID: formatChatID(callback.Message.Chat.ID),
		Content:   callback.Data,
		Metadata: map[string]interface{}{
			"chat_id":        callback.Message.Chat.ID,
			"user_id":        callback.From.ID,
			"username":       callback.From.UserName,
			"first_name":     callback.From.FirstName,
			"last_name":      callback.From.LastName,
			"message_id":     callback.Message.MessageID,
			"callback_id":    callback.ID,
			"message_type":   "callback_query",
			"chat_type":      callback.Message.Chat.Type,
			"inline_message": callback.InlineMessageID != "",
		},
	}

	// Answer the callback query to remove loading state
	if c.bot != nil {
		callbackCfg := tgbotapi.NewCallback(callback.ID, "")
		c.bot.Request(callbackCfg)
	}

	// Send message to incoming channel
	select {
	case c.incoming <- msg:
		if c.logger != nil {
			c.logger.Debug("Callback query received",
				"user_id", msg.UserID,
				"callback_data", callback.Data,
			)
		}
	default:
		if c.logger != nil {
			c.logger.Warn("Incoming channel is full, callback dropped")
		}
	}
}

// extractMessageContent extracts content and metadata from a Telegram message
func (c *Connector) extractMessageContent(message *tgbotapi.Message) (string, map[string]interface{}) {
	metadata := map[string]interface{}{
		"chat_id":    message.Chat.ID,
		"user_id":    message.From.ID,
		"username":   message.From.UserName,
		"first_name": message.From.FirstName,
		"last_name":  message.From.LastName,
		"message_id": message.MessageID,
		"chat_type":  message.Chat.Type,
	}

	var content string

	// Handle different message types
	if message.IsCommand() {
		// Command message (/start, /help, etc.)
		content = message.Text
		metadata["message_type"] = "command"
		metadata["command"] = message.Command()
		metadata["command_args"] = message.CommandArguments()
	} else if message.Text != "" {
		// Text message
		content = message.Text
		metadata["message_type"] = "text"
	} else if message.Photo != nil && len(message.Photo) > 0 {
		// Photo message
		photo := message.Photo[len(message.Photo)-1] // Get highest resolution photo
		content = fmt.Sprintf("[Photo] FileID: %s, Width: %d, Height: %d",
			photo.FileID, photo.Width, photo.Height)
		metadata["message_type"] = "photo"
		metadata["photo_file_id"] = photo.FileID
		metadata["photo_width"] = photo.Width
		metadata["photo_height"] = photo.Height
		metadata["photo_file_size"] = photo.FileSize
		if message.Caption != "" {
			content = fmt.Sprintf("%s\nCaption: %s", content, message.Caption)
			metadata["caption"] = message.Caption
		}
	} else if message.Document != nil {
		// Document message
		content = fmt.Sprintf("[Document] FileID: %s, FileName: %s, FileSize: %d",
			message.Document.FileID, message.Document.FileName, message.Document.FileSize)
		metadata["message_type"] = "document"
		metadata["document_file_id"] = message.Document.FileID
		metadata["document_file_name"] = message.Document.FileName
		metadata["document_file_size"] = message.Document.FileSize
		metadata["document_mime_type"] = message.Document.MimeType
		if message.Caption != "" {
			content = fmt.Sprintf("%s\nCaption: %s", content, message.Caption)
			metadata["caption"] = message.Caption
		}
	} else if message.Audio != nil {
		// Audio message
		duration := 0
		if message.Audio.Duration != 0 {
			duration = message.Audio.Duration
		}
		content = fmt.Sprintf("[Audio] FileID: %s, Duration: %ds", message.Audio.FileID, duration)
		metadata["message_type"] = "audio"
		metadata["audio_file_id"] = message.Audio.FileID
		metadata["audio_duration"] = duration
		metadata["audio_file_size"] = message.Audio.FileSize
		metadata["audio_mime_type"] = message.Audio.MimeType
		if message.Audio.Performer != "" {
			metadata["audio_performer"] = message.Audio.Performer
		}
		if message.Audio.Title != "" {
			metadata["audio_title"] = message.Audio.Title
		}
	} else if message.Voice != nil {
		// Voice message
		duration := 0
		if message.Voice.Duration != 0 {
			duration = message.Voice.Duration
		}
		content = fmt.Sprintf("[Voice] FileID: %s, Duration: %ds", message.Voice.FileID, duration)
		metadata["message_type"] = "voice"
		metadata["voice_file_id"] = message.Voice.FileID
		metadata["voice_duration"] = duration
		metadata["voice_file_size"] = message.Voice.FileSize
	} else if message.Video != nil {
		// Video message
		duration := 0
		if message.Video.Duration != 0 {
			duration = message.Video.Duration
		}
		content = fmt.Sprintf("[Video] FileID: %s, Width: %d, Height: %d, Duration: %ds",
			message.Video.FileID, message.Video.Width, message.Video.Height, duration)
		metadata["message_type"] = "video"
		metadata["video_file_id"] = message.Video.FileID
		metadata["video_width"] = message.Video.Width
		metadata["video_height"] = message.Video.Height
		metadata["video_duration"] = duration
		metadata["video_file_size"] = message.Video.FileSize
		metadata["video_mime_type"] = message.Video.MimeType
		if message.Caption != "" {
			content = fmt.Sprintf("%s\nCaption: %s", content, message.Caption)
			metadata["caption"] = message.Caption
		}
	} else if message.VideoNote != nil {
		// Video note (round video message)
		duration := 0
		if message.VideoNote.Duration != 0 {
			duration = message.VideoNote.Duration
		}
		content = fmt.Sprintf("[VideoNote] FileID: %s, Duration: %ds, Length: %d",
			message.VideoNote.FileID, duration, message.VideoNote.Length)
		metadata["message_type"] = "video_note"
		metadata["video_note_file_id"] = message.VideoNote.FileID
		metadata["video_note_duration"] = duration
		metadata["video_note_length"] = message.VideoNote.Length
		metadata["video_note_file_size"] = message.VideoNote.FileSize
	} else if message.Sticker != nil {
		// Sticker message
		content = fmt.Sprintf("[Sticker] FileID: %s, Emoji: %s",
			message.Sticker.FileID, message.Sticker.Emoji)
		metadata["message_type"] = "sticker"
		metadata["sticker_file_id"] = message.Sticker.FileID
		metadata["sticker_emoji"] = message.Sticker.Emoji
		metadata["sticker_set_name"] = message.Sticker.SetName
		metadata["sticker_width"] = message.Sticker.Width
		metadata["sticker_height"] = message.Sticker.Height
		metadata["sticker_is_animated"] = message.Sticker.IsAnimated
	} else if message.Location != nil {
		// Location message
		content = fmt.Sprintf("[Location] Latitude: %.6f, Longitude: %.6f",
			message.Location.Latitude, message.Location.Longitude)
		metadata["message_type"] = "location"
		metadata["location_latitude"] = message.Location.Latitude
		metadata["location_longitude"] = message.Location.Longitude
	} else if message.Contact != nil {
		// Contact message
		content = fmt.Sprintf("[Contact] Phone: %s, Name: %s",
			message.Contact.PhoneNumber, message.Contact.FirstName)
		metadata["message_type"] = "contact"
		metadata["contact_phone_number"] = message.Contact.PhoneNumber
		metadata["contact_first_name"] = message.Contact.FirstName
		metadata["contact_last_name"] = message.Contact.LastName
		metadata["contact_user_id"] = message.Contact.UserID
	} else {
		// Unknown message type
		content = "[Unsupported message type]"
		metadata["message_type"] = "unknown"
	}

	// Add reply to information if exists
	if message.ReplyToMessage != nil {
		metadata["reply_to_message_id"] = message.ReplyToMessage.MessageID
		if message.ReplyToMessage.From != nil {
			metadata["reply_to_user_id"] = message.ReplyToMessage.From.ID
			metadata["reply_to_username"] = message.ReplyToMessage.From.UserName
		}
		if message.ReplyToMessage.Text != "" {
			metadata["reply_to_text"] = message.ReplyToMessage.Text
		}
	}

	// Add forward information if exists
	if message.ForwardFrom != nil {
		metadata["forward_from_user_id"] = message.ForwardFrom.ID
		metadata["forward_from_username"] = message.ForwardFrom.UserName
	}
	if message.ForwardDate != 0 {
		metadata["forward_date"] = message.ForwardDate
	}

	// Add edit information if edited
	if message.EditDate != 0 {
		metadata["edit_date"] = message.EditDate
	}

	return content, metadata
}

// isAllowed checks if the user or chat is allowed to interact with the bot
func (c *Connector) isAllowed(chatID int64, userID int64) bool {
	// Check if chat is in allowed chats list
	if len(c.config.AllowedChats) > 0 {
		chatIDStr := fmt.Sprintf("%d", chatID)
		for _, allowedChat := range c.config.AllowedChats {
			if allowedChat == chatIDStr {
				return true
			}
		}
	}

	// Check if user is in allowed users list
	if len(c.config.AllowedUsers) > 0 {
		userIDStr := fmt.Sprintf("%d", userID)
		for _, allowedUser := range c.config.AllowedUsers {
			if allowedUser == userIDStr {
				return true
			}
		}
	}

	return false
}

// getMode returns the current mode of operation (webhook or polling)
func (c *Connector) getMode() string {
	if c.config.WebhookURL != "" {
		return "webhook"
	}
	return "polling"
}

// formatUserID formats a user ID in the format "user_id:chat_id"
func formatUserID(userID, chatID int64) string {
	return fmt.Sprintf("%d:%d", userID, chatID)
}

// parseChatID extracts the chat ID from a formatted user ID
func parseChatID(userID string) (int64, error) {
	// Parse userID in format "user_id:chat_id"
	parts := strings.Split(userID, ":")
	if len(parts) != 2 {
		return 0, fmt.Errorf("invalid user ID format: expected 'user_id:chat_id', got '%s'", userID)
	}

	chatID, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid chat ID: %w", err)
	}

	return chatID, nil
}

// formatChatID formats a chat ID as a string
func formatChatID(chatID int64) string {
	return fmt.Sprintf("%d", chatID)
}
