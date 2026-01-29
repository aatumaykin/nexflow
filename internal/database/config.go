package database

import (
	"fmt"
	"time"
)

// DBConfig holds database configuration
type DBConfig struct {
	Type            string        // "sqlite" or "postgres"
	Path            string        // connection string or file path
	MigrationsPath  string        // path to migrations directory
	MaxOpenConns    int           // maximum open connections
	MaxIdleConns    int           // maximum idle connections
	ConnMaxLifetime time.Duration // maximum connection lifetime
}

// Validate checks if configuration is valid.
func (c *DBConfig) Validate() error {
	if c.Type == "" {
		return fmt.Errorf("database type is required")
	}
	if c.Path == "" {
		return fmt.Errorf("database path is required")
	}
	if c.Type != "sqlite" && c.Type != "postgres" {
		return fmt.Errorf("unsupported database type: %s, must be 'sqlite' or 'postgres'", c.Type)
	}
	return nil
}
