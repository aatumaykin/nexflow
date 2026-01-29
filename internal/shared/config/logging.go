package config

import (
	"fmt"
)

const (
	ValidLevels  = "debug, info, warn, error, fatal"
	ValidFormats = "json, text"
)

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
