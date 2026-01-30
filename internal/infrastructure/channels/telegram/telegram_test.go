package telegram

import (
	"context"
	"errors"
	"testing"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/config"
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
