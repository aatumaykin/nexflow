package config

import "fmt"

// ChannelsConfig represents channels configuration
type ChannelsConfig struct {
	Telegram TelegramConfig `json:"telegram" yaml:"telegram"`
	Discord  DiscordConfig  `json:"discord" yaml:"discord"`
	Web      WebConfig      `json:"web" yaml:"web"`
}

// TelegramConfig represents Telegram bot configuration
type TelegramConfig struct {
	Enabled      bool     `json:"enabled" yaml:"enabled"`
	BotToken     string   `json:"bot_token" yaml:"bot_token"`
	AllowedUsers []string `json:"allowed_users" yaml:"allowed_users"`
	AllowedChats []string `json:"allowed_chats" yaml:"allowed_chats"`
	WebhookURL   string   `json:"webhook_url" yaml:"webhook_url"` // Optional: use webhook instead of long polling
}

// DiscordConfig represents Discord bot configuration
type DiscordConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	BotToken string `json:"bot_token" yaml:"bot_token"`
}

// WebConfig represents web interface configuration
type WebConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}

// Validate validates the channels configuration
func (c *ChannelsConfig) Validate() error {
	if c.Telegram.Enabled {
		if c.Telegram.BotToken == "" {
			return fmt.Errorf("telegram bot_token is required when telegram is enabled")
		}
		// At least one of allowed_users or allowed_chats must be specified for security
		if len(c.Telegram.AllowedUsers) == 0 && len(c.Telegram.AllowedChats) == 0 {
			return fmt.Errorf("telegram: at least one of allowed_users or allowed_chats must be specified for security when telegram is enabled")
		}
	}
	if c.Discord.Enabled && c.Discord.BotToken == "" {
		return fmt.Errorf("discord bot_token is required when discord is enabled")
	}
	return nil
}
