package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTelegramConfig_Validation tests Telegram configuration validation
func TestTelegramConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      ChannelsConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid telegram config with allowed_users",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:      true,
					BotToken:     "test-token",
					AllowedUsers: []string{"123456789"},
				},
			},
			expectError: false,
		},
		{
			name: "valid telegram config with allowed_chats",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:      true,
					BotToken:     "test-token",
					AllowedChats: []string{"-1001234567890"},
				},
			},
			expectError: false,
		},
		{
			name: "valid telegram config with both allowed_users and allowed_chats",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:      true,
					BotToken:     "test-token",
					AllowedUsers: []string{"123456789"},
					AllowedChats: []string{"-1001234567890"},
				},
			},
			expectError: false,
		},
		{
			name: "telegram enabled but bot_token missing",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:      true,
					AllowedUsers: []string{"123456789"},
				},
			},
			expectError: true,
			errorMsg:    "bot_token is required",
		},
		{
			name: "telegram enabled but no allowed users or chats",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:  true,
					BotToken: "test-token",
				},
			},
			expectError: true,
			errorMsg:    "at least one of allowed_users or allowed_chats must be specified",
		},
		{
			name: "telegram disabled - no validation",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
		{
			name: "telegram with webhook_url",
			config: ChannelsConfig{
				Telegram: TelegramConfig{
					Enabled:      true,
					BotToken:     "test-token",
					AllowedUsers: []string{"123456789"},
					WebhookURL:   "https://example.com/webhook/telegram",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				require.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				require.NoError(t, err)
			}
		})
	}
}

// TestTelegramConfig_FullConfiguration tests complete Telegram configuration
func TestTelegramConfig_FullConfiguration(t *testing.T) {
	config := ChannelsConfig{
		Telegram: TelegramConfig{
			Enabled:      true,
			BotToken:     "${TELEGRAM_BOT_TOKEN}",
			AllowedUsers: []string{"123456789", "987654321"},
			AllowedChats: []string{"-1001234567890"},
			WebhookURL:   "https://mydomain.com/telegram/webhook",
		},
	}

	err := config.Validate()
	assert.NoError(t, err)

	assert.True(t, config.Telegram.Enabled)
	assert.Equal(t, "${TELEGRAM_BOT_TOKEN}", config.Telegram.BotToken)
	assert.Len(t, config.Telegram.AllowedUsers, 2)
	assert.Len(t, config.Telegram.AllowedChats, 1)
	assert.Equal(t, "https://mydomain.com/telegram/webhook", config.Telegram.WebhookURL)
}

// TestDiscordConfig_Validation tests Discord configuration validation
func TestDiscordConfig_Validation(t *testing.T) {
	tests := []struct {
		name        string
		config      ChannelsConfig
		expectError bool
	}{
		{
			name: "valid discord config",
			config: ChannelsConfig{
				Discord: DiscordConfig{
					Enabled:  true,
					BotToken: "test-token",
				},
			},
			expectError: false,
		},
		{
			name: "discord enabled but bot_token missing",
			config: ChannelsConfig{
				Discord: DiscordConfig{
					Enabled: true,
				},
			},
			expectError: true,
		},
		{
			name: "discord disabled - no validation",
			config: ChannelsConfig{
				Discord: DiscordConfig{
					Enabled: false,
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestWebConfig_NoValidation tests that Web configuration has no validation
func TestWebConfig_NoValidation(t *testing.T) {
	config := ChannelsConfig{
		Web: WebConfig{
			Enabled: true,
		},
	}

	// Web config has no required fields, should always pass validation
	err := config.Validate()
	assert.NoError(t, err)
}
