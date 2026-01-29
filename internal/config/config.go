package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
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
// This is a universal function that processes all string fields recursively
func expandEnvVars(config *Config) error {
	return expandValue(reflect.ValueOf(config).Elem())
}

// expandValue recursively processes all string fields in a struct or map
func expandValue(v reflect.Value) error {
	// Skip invalid or unexported values
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.String:
		if v.CanSet() {
			v.SetString(expandAllEnvVars(v.String()))
		}
	case reflect.Struct:
		// Skip unexported fields
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// Check if field is exported and can be set
			if field.CanInterface() && field.CanAddr() {
				if err := expandValue(field); err != nil {
					return err
				}
			}
		}
	case reflect.Map:
		// Skip nil maps
		if v.IsNil() {
			return nil
		}
		// Create a new map to hold modified values
		mapType := v.Type()
		elemType := mapType.Elem()
		newMap := reflect.MakeMapWithSize(mapType, v.Len())

		for _, key := range v.MapKeys() {
			value := v.MapIndex(key)
			if !value.IsValid() {
				continue
			}

			// Create a copy of the value to expand
			newValue := reflect.New(elemType).Elem()
			newValue.Set(value)

			// Expand the copy
			if err := expandValue(newValue); err != nil {
				return err
			}

			newMap.SetMapIndex(key, newValue)
		}

		if v.CanSet() {
			v.Set(newMap)
		}
	case reflect.Slice:
		// Skip nil slices
		if v.IsNil() {
			return nil
		}
		// Process each element
		for i := 0; i < v.Len(); i++ {
			element := v.Index(i)
			if element.CanAddr() && element.CanSet() {
				if err := expandValue(element); err != nil {
					return err
				}
			}
		}
	case reflect.Pointer:
		// Skip nil pointers
		if v.IsNil() {
			return nil
		}
		// Dereference and process
		if v.Elem().CanAddr() {
			return expandValue(v.Elem())
		}
	}

	return nil
}

// expandAllEnvVars expands all environment variable references in a string
// Supports multiple ${VAR} patterns in a single string
func expandAllEnvVars(s string) string {
	// Use regex to find all ${VAR_NAME} patterns
	re := regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)\}`)

	// Replace function that looks up each variable in the environment
	result := re.ReplaceAllStringFunc(s, func(match string) string {
		// Extract the variable name (remove ${ and })
		varName := match[2 : len(match)-1]
		varValue := os.Getenv(varName)

		if varValue == "" {
			// If environment variable is not set, keep the template
			return match
		}

		return varValue
	})

	return result
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
