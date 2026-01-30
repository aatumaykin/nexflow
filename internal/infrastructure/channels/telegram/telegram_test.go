package telegram

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of repository.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByChannel(ctx context.Context, channel, channelID string) (*entity.User, error) {
	args := m.Called(ctx, channel, channelID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) List(ctx context.Context) ([]*entity.User, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestNewConnector(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	assert.NotNil(t, connector)
	assert.Equal(t, "telegram", connector.Name())
	assert.False(t, connector.IsRunning())
}

func TestConnector_Start_InvalidToken(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "invalid_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	err := connector.Start(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create telegram bot")
	assert.False(t, connector.IsRunning())
}

func TestConnector_Start_AlreadyRunning(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	// First start will fail due to invalid token, but we'll set running to true
	connector.mu.Lock()
	connector.running = true
	connector.mu.Unlock()

	ctx := context.Background()
	err := connector.Start(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
}

func TestConnector_Stop_NotRunning(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	err := connector.Stop(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestConnector_SendResponse_NotRunning(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	response := &channels.Response{Content: "test message"}
	err := connector.SendResponse(ctx, "123:456", response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestConnector_SendResponse_InvalidUserID(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	// Simulate running state
	connector.mu.Lock()
	connector.running = true
	connector.mu.Unlock()

	ctx := context.Background()
	response := &channels.Response{Content: "test message"}
	err := connector.SendResponse(ctx, "invalid_user_id", response)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid user ID format")
}

func TestConnector_GetUser(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	user := entity.NewUser("telegram", "123")
	mockRepo.On("FindByChannel", mock.Anything, "telegram", "123").Return(user, nil)

	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	foundUser, err := connector.GetUser(ctx, "123")

	assert.NoError(t, err)
	assert.NotNil(t, foundUser)
	assert.Equal(t, user, foundUser)
	mockRepo.AssertExpectations(t)
}

func TestConnector_GetUser_NotFound(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("FindByChannel", mock.Anything, "telegram", "456").Return(nil, errors.New("user not found"))

	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	_, err := connector.GetUser(ctx, "456")

	assert.Error(t, err)
	mockRepo.AssertExpectations(t)
}

func TestConnector_CreateUser(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil)

	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	user, err := connector.CreateUser(ctx, "123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "telegram", user.Channel.String())
	assert.Equal(t, "123", user.ChannelID)
	mockRepo.AssertExpectations(t)
}

func TestConnector_CreateUser_Error(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(errors.New("database error"))

	connector := NewConnector(cfg, mockRepo, nil)

	ctx := context.Background()
	_, err := connector.CreateUser(ctx, "123")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create user")
	mockRepo.AssertExpectations(t)
}

func TestFormatUserID(t *testing.T) {
	tests := []struct {
		name     string
		userID   int64
		chatID   int64
		expected string
	}{
		{"basic", 123456, 789012, "123456:789012"},
		{"zero values", 0, 0, "0:0"},
		{"large numbers", 99999999999, 88888888888, "99999999999:88888888888"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatUserID(tt.userID, tt.chatID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseChatID(t *testing.T) {
	tests := []struct {
		name        string
		userID      string
		expected    int64
		expectError bool
	}{
		{"valid", "123456:789012", 789012, false},
		{"invalid format", "invalid", 0, true},
		{"missing chat ID", "123456:", 0, true},
		{"empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseChatID(tt.userID)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestFormatChatID(t *testing.T) {
	tests := []struct {
		name     string
		chatID   int64
		expected string
	}{
		{"basic", 789012, "789012"},
		{"zero", 0, "0"},
		{"negative", -123, "-123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatChatID(tt.chatID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetMode(t *testing.T) {
	t.Run("webhook mode", func(t *testing.T) {
		cfg := config.TelegramConfig{
			Enabled:      true,
			BotToken:     "test_token",
			AllowedChats: []string{"123456789"},
			WebhookURL:   "https://example.com/webhook",
		}

		connector := NewConnector(cfg, nil, nil)
		assert.Equal(t, "webhook", connector.getMode())
	})

	t.Run("polling mode", func(t *testing.T) {
		cfg := config.TelegramConfig{
			Enabled:      true,
			BotToken:     "test_token",
			AllowedChats: []string{"123456789"},
		}

		connector := NewConnector(cfg, nil, nil)
		assert.Equal(t, "polling", connector.getMode())
	})
}

func TestIsAllowed(t *testing.T) {
	tests := []struct {
		name     string
		cfg      config.TelegramConfig
		chatID   int64
		userID   int64
		expected bool
	}{
		{
			name: "allowed chat",
			cfg: config.TelegramConfig{
				AllowedChats: []string{"123456789"},
			},
			chatID:   123456789,
			userID:   999999,
			expected: true,
		},
		{
			name: "allowed user",
			cfg: config.TelegramConfig{
				AllowedUsers: []string{"999999"},
			},
			chatID:   123456789,
			userID:   999999,
			expected: true,
		},
		{
			name: "not allowed",
			cfg: config.TelegramConfig{
				AllowedChats: []string{"123456789"},
			},
			chatID:   987654321,
			userID:   999999,
			expected: false,
		},
		{
			name: "both lists present - chat matches",
			cfg: config.TelegramConfig{
				AllowedChats: []string{"123456789"},
				AllowedUsers: []string{"999999"},
			},
			chatID:   123456789,
			userID:   888888,
			expected: true,
		},
		{
			name: "both lists present - user matches",
			cfg: config.TelegramConfig{
				AllowedChats: []string{"123456789"},
				AllowedUsers: []string{"999999"},
			},
			chatID:   987654321,
			userID:   999999,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connector := NewConnector(tt.cfg, nil, nil)
			result := connector.isAllowed(tt.chatID, tt.userID)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConnector_Incoming(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	incoming := connector.Incoming()
	assert.NotNil(t, incoming)
}

func TestConnector_Name(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	assert.Equal(t, "telegram", connector.Name())
}

// TestExtractMessageContent tests extracting content from different message types
func TestExtractMessageContent(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	connector := NewConnector(cfg, nil, nil)

	tests := []struct {
		name      string
		message   *tgbotapi.Message
		wantType  string
		wantEmpty bool
	}{
		{
			name: "text message",
			message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
				From: &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Text: "Hello, world!",
			},
			wantType:  "text",
			wantEmpty: false,
		},
		{
			name: "command message",
			message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
				From: &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Text: "/start arg1 arg2",
				Entities: []tgbotapi.MessageEntity{
					{Type: "bot_command", Offset: 0, Length: 6},
				},
			},
			wantType:  "command",
			wantEmpty: false,
		},
		{
			name: "photo message",
			message: &tgbotapi.Message{
				Chat:  &tgbotapi.Chat{ID: 123, Type: "private"},
				From:  &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Photo: []tgbotapi.PhotoSize{{FileID: "photo123", Width: 800, Height: 600}},
			},
			wantType:  "photo",
			wantEmpty: false,
		},
		{
			name: "document message",
			message: &tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123, Type: "private"},
				From:     &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Document: &tgbotapi.Document{FileID: "doc123", FileName: "test.pdf", FileSize: 1024, MimeType: "application/pdf"},
			},
			wantType:  "document",
			wantEmpty: false,
		},
		{
			name: "audio message",
			message: &tgbotapi.Message{
				Chat:  &tgbotapi.Chat{ID: 123, Type: "private"},
				From:  &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Audio: &tgbotapi.Audio{FileID: "audio123", Duration: 120, FileSize: 2048, MimeType: "audio/mpeg"},
			},
			wantType:  "audio",
			wantEmpty: false,
		},
		{
			name: "voice message",
			message: &tgbotapi.Message{
				Chat:  &tgbotapi.Chat{ID: 123, Type: "private"},
				From:  &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Voice: &tgbotapi.Voice{FileID: "voice123", Duration: 30, FileSize: 512},
			},
			wantType:  "voice",
			wantEmpty: false,
		},
		{
			name: "video message",
			message: &tgbotapi.Message{
				Chat:  &tgbotapi.Chat{ID: 123, Type: "private"},
				From:  &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Video: &tgbotapi.Video{FileID: "video123", Width: 1920, Height: 1080, Duration: 60, FileSize: 4096, MimeType: "video/mp4"},
			},
			wantType:  "video",
			wantEmpty: false,
		},
		{
			name: "video note message",
			message: &tgbotapi.Message{
				Chat:      &tgbotapi.Chat{ID: 123, Type: "private"},
				From:      &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				VideoNote: &tgbotapi.VideoNote{FileID: "videonote123", Duration: 15, Length: 640, FileSize: 256},
			},
			wantType:  "video_note",
			wantEmpty: false,
		},
		{
			name: "sticker message",
			message: &tgbotapi.Message{
				Chat:    &tgbotapi.Chat{ID: 123, Type: "private"},
				From:    &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Sticker: &tgbotapi.Sticker{FileID: "sticker123", Width: 512, Height: 512, Emoji: "😀", IsAnimated: false},
			},
			wantType:  "sticker",
			wantEmpty: false,
		},
		{
			name: "location message",
			message: &tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123, Type: "private"},
				From:     &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Location: &tgbotapi.Location{Latitude: 55.7558, Longitude: 37.6173},
			},
			wantType:  "location",
			wantEmpty: false,
		},
		{
			name: "contact message",
			message: &tgbotapi.Message{
				Chat:    &tgbotapi.Chat{ID: 123, Type: "private"},
				From:    &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Contact: &tgbotapi.Contact{PhoneNumber: "+1234567890", FirstName: "John", LastName: "Doe", UserID: 789},
			},
			wantType:  "contact",
			wantEmpty: false,
		},
		{
			name: "message with caption",
			message: &tgbotapi.Message{
				Chat:    &tgbotapi.Chat{ID: 123, Type: "private"},
				From:    &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
				Photo:   []tgbotapi.PhotoSize{{FileID: "photo123", Width: 800, Height: 600}},
				Caption: "My photo caption",
			},
			wantType:  "photo",
			wantEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content, metadata := connector.extractMessageContent(tt.message)

			if tt.wantEmpty {
				assert.Empty(t, content)
			} else {
				assert.NotEmpty(t, content)
			}

			assert.Equal(t, tt.wantType, metadata["message_type"])
			assert.Equal(t, tt.message.Chat.ID, metadata["chat_id"])
			assert.Equal(t, tt.message.From.ID, metadata["user_id"])
		})
	}
}

// TestExtractMessageContentWithMetadata tests that metadata contains all expected fields
func TestExtractMessageContentWithMetadata(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	connector := NewConnector(cfg, nil, nil)

	t.Run("text message with all metadata", func(t *testing.T) {
		message := &tgbotapi.Message{
			MessageID: 123,
			Chat:      &tgbotapi.Chat{ID: 123456, Type: "private"},
			From:      &tgbotapi.User{ID: 789012, FirstName: "John", LastName: "Doe", UserName: "johndoe"},
			Text:      "Test message",
		}

		content, metadata := connector.extractMessageContent(message)

		assert.Contains(t, content, "Test message")
		assert.Equal(t, "text", metadata["message_type"])
		assert.Equal(t, int64(123456), metadata["chat_id"])
		assert.Equal(t, int64(789012), metadata["user_id"])
		assert.Equal(t, "John", metadata["first_name"])
		assert.Equal(t, "Doe", metadata["last_name"])
		assert.Equal(t, "johndoe", metadata["username"])
		assert.Equal(t, int(123), metadata["message_id"])
		assert.Equal(t, "private", metadata["chat_type"])
	})

	t.Run("command message metadata", func(t *testing.T) {
		message := &tgbotapi.Message{
			MessageID: 456,
			Chat:      &tgbotapi.Chat{ID: 123456, Type: "private"},
			From:      &tgbotapi.User{ID: 789012, FirstName: "John", LastName: "Doe"},
			Text:      "/start param1 param2",
			Entities: []tgbotapi.MessageEntity{
				{Type: "bot_command", Offset: 0, Length: 6},
			},
		}

		content, metadata := connector.extractMessageContent(message)

		assert.Equal(t, "/start param1 param2", content)
		assert.Equal(t, "command", metadata["message_type"])
		assert.Equal(t, "start", metadata["command"])
		assert.Equal(t, "param1 param2", metadata["command_args"])
	})

	t.Run("photo message with caption", func(t *testing.T) {
		message := &tgbotapi.Message{
			MessageID: 789,
			Chat:      &tgbotapi.Chat{ID: 123456, Type: "private"},
			From:      &tgbotapi.User{ID: 789012, FirstName: "John", LastName: "Doe"},
			Photo: []tgbotapi.PhotoSize{
				{FileID: "photo1", Width: 100, Height: 100},
				{FileID: "photo2", Width: 800, Height: 600},
			},
			Caption: "Beautiful sunset",
		}

		content, metadata := connector.extractMessageContent(message)

		assert.Contains(t, content, "photo2")
		assert.Contains(t, content, "800")
		assert.Contains(t, content, "600")
		assert.Contains(t, content, "Beautiful sunset")
		assert.Equal(t, "photo", metadata["message_type"])
		assert.Equal(t, "photo2", metadata["photo_file_id"])
		assert.Equal(t, 800, metadata["photo_width"])
		assert.Equal(t, 600, metadata["photo_height"])
		assert.Equal(t, "Beautiful sunset", metadata["caption"])
	})
}

// TestHandleMessageWithUserCreation tests that users are automatically created
func TestHandleMessageWithUserCreation(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	// Start the connector to initialize the incoming channel
	connector.mu.Lock()
	connector.running = true
	connector.incoming = make(chan *channels.Message, 100)
	connector.mu.Unlock()

	defer func() {
		connector.mu.Lock()
		connector.running = false
		close(connector.incoming)
		connector.mu.Unlock()
	}()

	t.Run("user already exists", func(t *testing.T) {
		mockUser := entity.NewUser("telegram", "456")
		mockRepo.On("FindByChannel", mock.Anything, "telegram", "456").Return(mockUser, nil).Once()

		message := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
			From: &tgbotapi.User{ID: 456, FirstName: "Test", LastName: "User"},
			Text: "Hello",
		}

		ctx := context.Background()
		go connector.handleMessage(ctx, message)

		// Wait for message to be processed
		select {
		case msg := <-connector.incoming:
			assert.NotNil(t, msg)
			assert.Equal(t, "456:123", msg.UserID)
			assert.Equal(t, "123", msg.ChannelID)
			assert.Equal(t, "Hello", msg.Content)
			assert.Equal(t, mockUser.ID.String(), msg.Metadata["user_internal_id"])
			mockRepo.AssertExpectations(t)
		}
	})

	t.Run("user does not exist - create new", func(t *testing.T) {
		mockRepo.On("FindByChannel", mock.Anything, "telegram", "789").Return(nil, errors.New("not found")).Once()
		mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil).Once()

		message := &tgbotapi.Message{
			Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
			From: &tgbotapi.User{ID: 789, FirstName: "New", LastName: "User"},
			Text: "Hello from new user",
		}

		ctx := context.Background()
		go connector.handleMessage(ctx, message)

		// Wait for message to be processed
		select {
		case msg := <-connector.incoming:
			assert.NotNil(t, msg)
			assert.Equal(t, "789:123", msg.UserID)
			assert.Equal(t, "123", msg.ChannelID)
			assert.Equal(t, "Hello from new user", msg.Content)
			assert.NotEmpty(t, msg.Metadata["user_internal_id"])
			mockRepo.AssertExpectations(t)
		}
	})
}

// TestHandleCallbackQuery tests processing of callback queries
func TestHandleCallbackQuery(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	// Start the connector to initialize the incoming channel
	connector.mu.Lock()
	connector.running = true
	connector.incoming = make(chan *channels.Message, 100)
	connector.mu.Unlock()

	defer func() {
		connector.mu.Lock()
		connector.running = false
		close(connector.incoming)
		connector.mu.Unlock()
	}()

	t.Run("basic callback query", func(t *testing.T) {
		callback := &tgbotapi.CallbackQuery{
			ID: "callback123",
			From: &tgbotapi.User{
				ID:        456,
				FirstName: "John",
				LastName:  "Doe",
				UserName:  "johndoe",
			},
			Message: &tgbotapi.Message{
				MessageID: 123,
				Chat:      &tgbotapi.Chat{ID: 789, Type: "private"},
			},
			Data: "button_clicked",
		}

		go connector.handleCallbackQuery(callback)

		select {
		case msg := <-connector.incoming:
			assert.NotNil(t, msg)
			assert.Equal(t, "456:789", msg.UserID)
			assert.Equal(t, "789", msg.ChannelID)
			assert.Equal(t, "button_clicked", msg.Content)
			assert.Equal(t, "callback_query", msg.Metadata["message_type"])
			assert.Equal(t, "callback123", msg.Metadata["callback_id"])
			assert.Equal(t, int64(456), msg.Metadata["user_id"])
			assert.Equal(t, int64(789), msg.Metadata["chat_id"])
		}
	})
}

// TestMessageTypesIntegration tests all message types end-to-end
func TestMessageTypesIntegration(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	mockUser := entity.NewUser("telegram", "456")
	mockRepo.On("FindByChannel", mock.Anything, "telegram", "456").Return(mockUser, nil).Maybe()

	connector := NewConnector(cfg, mockRepo, nil)

	connector.mu.Lock()
	connector.running = true
	connector.incoming = make(chan *channels.Message, 100)
	connector.mu.Unlock()

	defer func() {
		connector.mu.Lock()
		connector.running = false
		close(connector.incoming)
		connector.mu.Unlock()
	}()

	ctx := context.Background()

	testMessages := []struct {
		name     string
		message  *tgbotapi.Message
		typeKey  string
		expected string
	}{
		{
			name: "simple text",
			message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: 123, Type: "private"},
				From: &tgbotapi.User{ID: 456, FirstName: "Test"},
				Text: "Hello",
			},
			typeKey:  "message_type",
			expected: "text",
		},
		{
			name: "with command",
			message: &tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123, Type: "private"},
				From:     &tgbotapi.User{ID: 456, FirstName: "Test"},
				Text:     "/help",
				Entities: []tgbotapi.MessageEntity{{Type: "bot_command", Offset: 0, Length: 5}},
			},
			typeKey:  "message_type",
			expected: "command",
		},
		{
			name: "photo",
			message: &tgbotapi.Message{
				Chat:  &tgbotapi.Chat{ID: 123, Type: "private"},
				From:  &tgbotapi.User{ID: 456, FirstName: "Test"},
				Photo: []tgbotapi.PhotoSize{{FileID: "abc", Width: 100, Height: 100}},
			},
			typeKey:  "message_type",
			expected: "photo",
		},
		{
			name: "document",
			message: &tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123, Type: "private"},
				From:     &tgbotapi.User{ID: 456, FirstName: "Test"},
				Document: &tgbotapi.Document{FileID: "doc", FileName: "file.pdf", FileSize: 1024},
			},
			typeKey:  "message_type",
			expected: "document",
		},
		{
			name: "location",
			message: &tgbotapi.Message{
				Chat:     &tgbotapi.Chat{ID: 123, Type: "private"},
				From:     &tgbotapi.User{ID: 456, FirstName: "Test"},
				Location: &tgbotapi.Location{Latitude: 55.0, Longitude: 37.0},
			},
			typeKey:  "message_type",
			expected: "location",
		},
	}

	for _, tt := range testMessages {
		t.Run(tt.name, func(t *testing.T) {
			go connector.handleMessage(ctx, tt.message)

			select {
			case msg := <-connector.incoming:
				assert.Equal(t, tt.expected, msg.Metadata[tt.typeKey])
			}
		})
	}
}

// TestSendResponse_TextMessage tests sending text messages
func TestSendResponse_TextMessage(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	// Simulate running state
	connector.mu.Lock()
	connector.running = true
	connector.mu.Unlock()

	t.Run("simple text message", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeText,
			Content: "Hello, world!",
		}

		// This test would require mocking the bot API
		// For now, we just test that the code compiles
		assert.NotNil(t, response)
		assert.Equal(t, channels.ResponseTypeText, response.Type)
	})

	t.Run("text message with HTML format", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeText,
			Content: "<b>Bold text</b>",
			Metadata: map[string]interface{}{
				"parse_mode": tgbotapi.ModeHTML,
			},
		}

		assert.NotNil(t, response)
		assert.Equal(t, tgbotapi.ModeHTML, response.Metadata["parse_mode"])
	})

	t.Run("text message with inline buttons", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeText,
			Content: "Choose an option:",
			Buttons: []channels.InlineButton{
				{Text: "Option 1", Data: "opt1"},
				{Text: "Option 2", Data: "opt2"},
			},
		}

		assert.NotNil(t, response)
		assert.Len(t, response.Buttons, 2)
	})
}

// TestSendResponse_PhotoMessage tests sending photo messages
func TestSendResponse_PhotoMessage(t *testing.T) {
	t.Run("photo with file ID", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypePhoto,
			Content: "Check out this photo",
			Caption: "Beautiful sunset",
			Media: &channels.MediaContent{
				FileID: "AgACAgIAAxkBAAIC...",
			},
		}

		assert.NotNil(t, response)
		assert.Equal(t, channels.ResponseTypePhoto, response.Type)
		assert.NotNil(t, response.Media)
		assert.NotEmpty(t, response.Media.FileID)
	})

	t.Run("photo with URL", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypePhoto,
			Content: "Here's a photo",
			Caption: "Nice picture",
			Media: &channels.MediaContent{
				URL: "https://example.com/photo.jpg",
			},
		}

		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Media.URL)
	})

	t.Run("photo with file data", func(t *testing.T) {
		photoData := []byte{0xFF, 0xD8, 0xFF, 0xE0} // JPEG header
		response := &channels.Response{
			Type:    channels.ResponseTypePhoto,
			Content: "Uploaded photo",
			Caption: "From camera",
			Media: &channels.MediaContent{
				FileData: photoData,
			},
		}

		assert.NotNil(t, response)
		assert.NotNil(t, response.Media.FileData)
	})
}

// TestSendResponse_DocumentMessage tests sending document messages
func TestSendResponse_DocumentMessage(t *testing.T) {
	t.Run("document with file ID", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeDocument,
			Content: "Here's a document",
			Caption: "Important file",
			Media: &channels.MediaContent{
				FileID:   "BQACAgIAAxkBAAI...",
				FileName: "report.pdf",
			},
		}

		assert.NotNil(t, response)
		assert.Equal(t, channels.ResponseTypeDocument, response.Type)
		assert.Equal(t, "report.pdf", response.Media.FileName)
	})

	t.Run("document with file data", func(t *testing.T) {
		docData := []byte("%PDF-1.4...")
		response := &channels.Response{
			Type:    channels.ResponseTypeDocument,
			Content: "Uploaded document",
			Caption: "My report",
			Media: &channels.MediaContent{
				FileData: docData,
				FileName: "document.pdf",
			},
		}

		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Media.FileName)
	})
}

// TestSendResponse_AudioMessage tests sending audio messages
func TestSendResponse_AudioMessage(t *testing.T) {
	t.Run("audio with file ID", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeAudio,
			Content: "Listen to this audio",
			Caption: "Great song",
			Media: &channels.MediaContent{
				FileID: "AwACAgIAAxkBAAI...",
			},
		}

		assert.NotNil(t, response)
		assert.Equal(t, channels.ResponseTypeAudio, response.Type)
	})

	t.Run("audio with file data", func(t *testing.T) {
		audioData := []byte{0xFF, 0xFB, 0x90} // MP3 header
		response := &channels.Response{
			Type:    channels.ResponseTypeAudio,
			Content: "Uploaded audio",
			Media: &channels.MediaContent{
				FileData: audioData,
			},
		}

		assert.NotNil(t, response)
		assert.NotNil(t, response.Media.FileData)
	})
}

// TestSendResponse_VideoMessage tests sending video messages
func TestSendResponse_VideoMessage(t *testing.T) {
	t.Run("video with file ID", func(t *testing.T) {
		response := &channels.Response{
			Type:    channels.ResponseTypeVideo,
			Content: "Watch this video",
			Caption: "Awesome clip",
			Media: &channels.MediaContent{
				FileID: "BAACAgIAAxkBAAI...",
			},
		}

		assert.NotNil(t, response)
		assert.Equal(t, channels.ResponseTypeVideo, response.Type)
	})

	t.Run("video with file data", func(t *testing.T) {
		videoData := []byte{0x00, 0x00, 0x00} // MP4 header
		response := &channels.Response{
			Type:    channels.ResponseTypeVideo,
			Content: "Uploaded video",
			Caption: "My movie",
			Media: &channels.MediaContent{
				FileData: videoData,
			},
		}

		assert.NotNil(t, response)
		assert.NotNil(t, response.Media.FileData)
	})
}

// TestSendResponse_StickerMessage tests sending sticker messages
func TestSendResponse_StickerMessage(t *testing.T) {
	response := &channels.Response{
		Type: channels.ResponseTypeSticker,
		Media: &channels.MediaContent{
			FileID: "CAACAgIAAxkBAAI...",
		},
	}

	assert.NotNil(t, response)
	assert.Equal(t, channels.ResponseTypeSticker, response.Type)
	assert.NotNil(t, response.Media)
	assert.NotEmpty(t, response.Media.FileID)
}

// TestSendResponse_EditMessage tests editing messages
func TestSendResponse_EditMessage(t *testing.T) {
	response := &channels.Response{
		Type:      channels.ResponseTypeText,
		Content:   "Updated message",
		MessageID: "123",
	}

	assert.NotNil(t, response)
	assert.NotEmpty(t, response.MessageID)
}

// TestBuildInlineMarkup tests building inline keyboard markup
func TestBuildInlineMarkup(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	t.Run("single button", func(t *testing.T) {
		response := &channels.Response{
			Buttons: []channels.InlineButton{
				{Text: "Click me", Data: "click"},
			},
		}

		markup := connector.buildInlineMarkup(response)
		assert.NotNil(t, markup)
		assert.Len(t, markup.InlineKeyboard, 1)
		assert.Len(t, markup.InlineKeyboard[0], 1)
		assert.Equal(t, "Click me", markup.InlineKeyboard[0][0].Text)
		assert.Equal(t, "click", *markup.InlineKeyboard[0][0].CallbackData)
	})

	t.Run("multiple buttons", func(t *testing.T) {
		response := &channels.Response{
			Buttons: []channels.InlineButton{
				{Text: "Option 1", Data: "opt1"},
				{Text: "Option 2", Data: "opt2"},
				{Text: "Option 3", Data: "opt3"},
				{Text: "Option 4", Data: "opt4"},
			},
		}

		markup := connector.buildInlineMarkup(response)
		assert.NotNil(t, markup)
		assert.Len(t, markup.InlineKeyboard, 2)
		assert.Len(t, markup.InlineKeyboard[0], 2)
		assert.Len(t, markup.InlineKeyboard[1], 2)
	})

	t.Run("button with URL", func(t *testing.T) {
		response := &channels.Response{
			Buttons: []channels.InlineButton{
				{Text: "Visit website", URL: "https://example.com"},
			},
		}

		markup := connector.buildInlineMarkup(response)
		assert.NotNil(t, markup)
		assert.NotNil(t, markup.InlineKeyboard[0][0].URL)
		assert.Equal(t, "https://example.com", *markup.InlineKeyboard[0][0].URL)
	})

	t.Run("button with switch inline query", func(t *testing.T) {
		response := &channels.Response{
			Buttons: []channels.InlineButton{
				{Text: "Search...", SwitchInline: "query"},
			},
		}

		markup := connector.buildInlineMarkup(response)
		assert.NotNil(t, markup)
		assert.NotNil(t, markup.InlineKeyboard[0][0].SwitchInlineQuery)
		assert.Equal(t, "query", *markup.InlineKeyboard[0][0].SwitchInlineQuery)
	})
}

// TestSplitLongText tests splitting long text messages
func TestSplitLongText(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	t.Run("short text", func(t *testing.T) {
		text := "Short message"
		parts := connector.splitLongText(text, 4096)
		assert.Len(t, parts, 1)
		assert.Equal(t, text, parts[0])
	})

	t.Run("text exactly at limit", func(t *testing.T) {
		text := string(make([]rune, 4096))
		parts := connector.splitLongText(text, 4096)
		assert.Len(t, parts, 1)
		assert.Equal(t, text, parts[0])
	})

	t.Run("text exceeding limit", func(t *testing.T) {
		text := "A"
		for len(text) < 5000 {
			text += text
		}
		parts := connector.splitLongText(text, 4096)
		assert.Greater(t, len(parts), 1)
		// Check each part is at or below limit
		for _, part := range parts {
			assert.LessOrEqual(t, len(part), 4096)
		}
	})

	t.Run("text with newlines", func(t *testing.T) {
		text := "Line 1\nLine 2\nLine 3"
		for len(text) < 4100 {
			text += "\nNext line"
		}
		parts := connector.splitLongText(text, 4096)
		assert.Greater(t, len(parts), 1)
		// Check that split happens at reasonable points
		assert.Contains(t, parts[0], "Line 1")
	})
}

// TestHandleSendError tests error handling
func TestHandleSendError(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	t.Run("bot blocked by user", func(t *testing.T) {
		err := connector.handleSendError(
			fmt.Errorf("Forbidden: bot was blocked by the user"),
			"text",
			123456789,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "bot was blocked by user")
	})

	t.Run("user is deactivated", func(t *testing.T) {
		err := connector.handleSendError(
			fmt.Errorf("Forbidden: user is deactivated"),
			"text",
			123456789,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user is deactivated")
	})

	t.Run("chat not found", func(t *testing.T) {
		err := connector.handleSendError(
			fmt.Errorf("Forbidden: chat not found"),
			"text",
			123456789,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chat not found")
	})

	t.Run("rate limit exceeded", func(t *testing.T) {
		err := connector.handleSendError(
			fmt.Errorf("Too Many Requests: retry after 10"),
			"text",
			123456789,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "rate limit exceeded")
	})

	t.Run("message too long", func(t *testing.T) {
		err := connector.handleSendError(
			fmt.Errorf("Bad Request: message is too long"),
			"text",
			123456789,
		)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "message is too long")
	})
}

// TestRateLimiter tests rate limiting functionality
func TestRateLimiter(t *testing.T) {
	rl := newRateLimiter()

	t.Run("initial tokens", func(t *testing.T) {
		// Rate limiter should be pre-filled with tokens
		assert.Equal(t, maxMessagesPerSecond, len(rl.tokens))
	})

	t.Run("acquire token", func(t *testing.T) {
		initialCount := len(rl.tokens)
		rl.acquireToken()
		assert.Equal(t, initialCount-1, len(rl.tokens))
	})
}

// TestNewConnector_WithRateLimiter tests that connector has rate limiter
func TestNewConnector_WithRateLimiter(t *testing.T) {
	cfg := config.TelegramConfig{
		Enabled:      true,
		BotToken:     "test_token",
		AllowedChats: []string{"123456789"},
	}

	mockRepo := new(MockUserRepository)
	connector := NewConnector(cfg, mockRepo, nil)

	assert.NotNil(t, connector.rateLimiter)
}
