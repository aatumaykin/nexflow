package config

import (
	"fmt"
	"time"
)

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Type            string        `json:"type" yaml:"type"`
	Path            string        `json:"path" yaml:"path"`
	MigrationsPath  string        `json:"migrations_path" yaml:"migrations_path"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
}

// Validate validates the database configuration
func (d *DatabaseConfig) Validate() error {
	if d.Type == "" {
		return fmt.Errorf("database.type is required")
	}
	if d.Path == "" {
		return fmt.Errorf("database.path is required")
	}
	if d.MigrationsPath == "" {
		return fmt.Errorf("database.migrations_path is required")
	}
	return nil
}
