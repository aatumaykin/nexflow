package config

// ChannelsConfig represents channels configuration
type ChannelsConfig struct {
	Telegram TelegramConfig `json:"telegram" yaml:"telegram"`
	Web      WebConfig      `json:"web" yaml:"web"`
}

// TelegramConfig represents Telegram bot configuration
type TelegramConfig struct {
	BotToken     string   `json:"bot_token" yaml:"bot_token"`
	AllowedUsers []string `json:"allowed_users" yaml:"allowed_users"`
}

// WebConfig represents web interface configuration
type WebConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}
