// Package config provides configuration management for Nexflow.
// It supports YAML and JSON configuration files with environment variable expansion.
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

// Constants for validation
const (
	MinPort        = 1
	MaxPort        = 65535
	DefaultTimeout = 30
	ValidLevels    = "debug, info, warn, error, fatal"
	ValidFormats   = "json, text"
	EnvVarPattern  = `\$\{([A-Za-z_][A-Za-z0-9_]*)\}`
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

// Validate validates the server configuration
func (s *ServerConfig) Validate() error {
	if s.Host == "" {
		return fmt.Errorf("server.host is required")
	}
	if s.Port < MinPort || s.Port > MaxPort {
		return fmt.Errorf("server.port must be between %d and %d", MinPort, MaxPort)
	}
	return nil
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type string `json:"type" yaml:"type"`
	Path string `json:"path" yaml:"path"`
}

// Validate validates the database configuration
func (d *DatabaseConfig) Validate() error {
	if d.Type == "" {
		return fmt.Errorf("database.type is required")
	}
	if d.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	return nil
}

// LLMConfig represents LLM provider configuration
type LLMConfig struct {
	DefaultProvider string                 `json:"default_provider" yaml:"default_provider"`
	Providers       map[string]LLMProvider `json:"providers" yaml:"providers"`
}

// Validate validates the LLM configuration
func (l *LLMConfig) Validate() error {
	if l.DefaultProvider == "" {
		return fmt.Errorf("llm.default_provider is required")
	}
	if len(l.Providers) == 0 {
		return fmt.Errorf("at least one llm provider is required")
	}
	if _, ok := l.Providers[l.DefaultProvider]; !ok {
		return fmt.Errorf("llm.default_provider '%s' not found in providers", l.DefaultProvider)
	}
	return nil
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

// Validate validates the skills configuration
func (s *SkillsConfig) Validate() error {
	if s.Directory == "" {
		return fmt.Errorf("skills.directory is required")
	}
	if s.TimeoutSec <= 0 {
		return fmt.Errorf("skills.timeout_sec must be positive")
	}
	return nil
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level" yaml:"level"`
	Format string `json:"format" yaml:"format"`
}

// Validate validates the logging configuration
func (l *LoggingConfig) Validate() error {
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true, "fatal": true,
	}
	if !validLevels[l.Level] {
		return fmt.Errorf("logging.level must be one of: %s", ValidLevels)
	}
	validFormats := map[string]bool{
		"json": true, "text": true,
	}
	if !validFormats[l.Format] {
		return fmt.Errorf("logging.format must be one of: %s", ValidFormats)
	}
	return nil
}

// Load loads configuration from a file (YAML or JSON).
// It expands environment variables in the format ${VAR_NAME} and validates the configuration.
// Returns an error if the file cannot be read, parsed, or if the configuration is invalid.
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
	re := regexp.MustCompile(EnvVarPattern)

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

// Validate validates the configuration by validating all sub-configurations
func (c *Config) Validate() error {
	if err := c.Server.Validate(); err != nil {
		return err
	}
	if err := c.Database.Validate(); err != nil {
		return err
	}
	if err := c.LLM.Validate(); err != nil {
		return err
	}
	if err := c.Skills.Validate(); err != nil {
		return err
	}
	if err := c.Logging.Validate(); err != nil {
		return err
	}
	return nil
}
