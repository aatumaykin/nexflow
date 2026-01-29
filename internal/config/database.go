package config

import (
	"fmt"
)

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
