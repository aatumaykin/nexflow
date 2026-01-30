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

			c.handleMessage(update.Message)
		}
	}
}

// handleMessage processes an incoming message from Telegram
func (c *Connector) handleMessage(message *tgbotapi.Message) {
	// Create message struct
	msg := &channels.Message{
		UserID:    formatUserID(message.From.ID, message.Chat.ID),
		ChannelID: formatChatID(message.Chat.ID),
		Content:   message.Text,
		Metadata: map[string]interface{}{
			"chat_id":    message.Chat.ID,
			"user_id":    message.From.ID,
			"username":   message.From.UserName,
			"first_name": message.From.FirstName,
			"last_name":  message.From.LastName,
			"message_id": message.MessageID,
			"is_command": message.IsCommand(),
			"chat_type":  message.Chat.Type,
		},
	}

	// Send message to incoming channel
	select {
	case c.incoming <- msg:
		if c.logger != nil {
			c.logger.Debug("Message received", "user_id", msg.UserID, "content", msg.Content)
		}
	default:
		if c.logger != nil {
			c.logger.Warn("Incoming channel is full, message dropped")
		}
	}
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
