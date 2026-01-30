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
	if c.Telegram.Enabled && c.Telegram.BotToken == "" {
		return fmt.Errorf("telegram bot_token is required when telegram is enabled")
	}
	if c.Discord.Enabled && c.Discord.BotToken == "" {
		return fmt.Errorf("discord bot_token is required when discord is enabled")
	}
	return nil
}
