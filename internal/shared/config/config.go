package config

import (
	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	LLM      LLMConfig      `yaml:"llm"`
	Channels ChannelsConfig `yaml:"channels"`
	Skills   SkillsConfig   `yaml:"skills"`
	Logging  LoggingConfig  `yaml:"logging"`
}

// Load loads configuration from a YAML file.
// It expands environment variables in the format ${VAR_NAME} and validates the configuration.
// Returns an error if the file cannot be read, parsed, or if the configuration is invalid.
func Load(path string) (*Config, error) {
	// Read file content
	data, err := readConfigFile(path)
	if err != nil {
		return nil, err
	}

	// Parse configuration
	config, err := parseConfig(data, path)
	if err != nil {
		return nil, err
	}

	// Expand environment variables
	if err := expandEnvVars(config); err != nil {
		return nil, err
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
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

// parseConfig parses configuration data from YAML file
func parseConfig(data []byte, path string) (*Config, error) {
	var config Config
	ext := getFileExtension(path)

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
	default:
		return nil, errUnsupportedFormat(ext)
	}

	return &config, nil
}
