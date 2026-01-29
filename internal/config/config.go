package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `json:"server" yaml:"server"`
	Database DatabaseConfig `json:"database" yaml:"database"`
	LLM      LLMConfig      `json:"llm" yaml:"llm"`
	Channels ChannelsConfig `json:"channels" yaml:"channels"`
	Skills   SkillsConfig   `json:"skills" yaml:"skills"`
	Logging  LoggingConfig  `json:"logging" yaml:"logging"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type string `json:"type" yaml:"type"`
	Path string `json:"path" yaml:"path"`
}

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	DefaultProvider string                 `json:"default_provider" yaml:"default_provider"`
	Providers       map[string]LLMProvider `json:"providers" yaml:"providers"`
}

// LLMProvider represents a single LLM provider configuration
type LLMProvider struct {
	APIKey  string `json:"api_key" yaml:"api_key"`
	BaseURL string `json:"base_url" yaml:"base_url"`
	Model   string `json:"model" yaml:"model"`
}

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

// SkillsConfig represents skills configuration
type SkillsConfig struct {
	Directory      string `json:"directory" yaml:"directory"`
	TimeoutSec     int    `json:"timeout_sec" yaml:"timeout_sec"`
	SandboxEnabled bool   `json:"sandbox_enabled" yaml:"sandbox_enabled"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
}

// Load loads configuration from a file (YAML or JSON)
func Load(path string) (*Config, error) {
	// Read file content
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Determine file type and parse accordingly
	var config Config
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config: %w", err)
		}
	case ".json":
		if err := json.Unmarshal(data, &config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	// Expand environment variables
	if err := expandEnvVars(&config); err != nil {
		return nil, fmt.Errorf("failed to expand environment variables: %w", err)
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// expandEnvVars expands environment variable references in the config
func expandEnvVars(config *Config) error {
	// Expand in LLM providers
	for providerName, provider := range config.LLM.Providers {
		updated := provider
		if strings.Contains(provider.APIKey, "${") && strings.Contains(provider.APIKey, "}") {
			updated.APIKey = expandVar(provider.APIKey)
		}
		if strings.Contains(provider.BaseURL, "${") && strings.Contains(provider.BaseURL, "}") {
			updated.BaseURL = expandVar(provider.BaseURL)
		}
		config.LLM.Providers[providerName] = updated
	}

	// Expand in channels
	if strings.Contains(config.Channels.Telegram.BotToken, "${") && strings.Contains(config.Channels.Telegram.BotToken, "}") {
		config.Channels.Telegram.BotToken = expandVar(config.Channels.Telegram.BotToken)
	}

	return nil
}

// expandVar expands a single environment variable reference
func expandVar(s string) string {
	start := strings.Index(s, "${")
	if start == -1 {
		return s
	}
	end := strings.Index(s, "}")
	if end == -1 || end <= start {
		return s
	}

	varName := s[start+2 : end]
	varValue := os.Getenv(varName)

	if varValue == "" {
		// If environment variable is not set, keep the template
		return s
	}

	return s[:start] + varValue + s[end+1:]
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Validate server config
	if c.Server.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	// Validate database config
	if c.Database.Type == "" {
		return fmt.Errorf("database.type is required")
	}
	if c.Database.Path == "" {
		return fmt.Errorf("database.path is required")
	}

	// Validate LLM config
	if c.LLM.DefaultProvider == "" {
		return fmt.Errorf("llm.default_provider is required")
	}
	if len(c.LLM.Providers) == 0 {
		return fmt.Errorf("at least one llm provider is required")
	}
	if _, ok := c.LLM.Providers[c.LLM.DefaultProvider]; !ok {
		return fmt.Errorf("llm.default_provider '%s' not found in providers", c.LLM.DefaultProvider)
	}

	// Validate skills config
	if c.Skills.Directory == "" {
		return fmt.Errorf("skills.directory is required")
	}
	if c.Skills.TimeoutSec <= 0 {
		return fmt.Errorf("skills.timeout_sec must be positive")
	}

	// Validate logging config
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("logging.level must be one of: debug, info, warn, error, fatal")
	}
	validFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validFormats[c.Logging.Format] {
		return fmt.Errorf("logging.format must be one of: json, text")
	}

	return nil
}
