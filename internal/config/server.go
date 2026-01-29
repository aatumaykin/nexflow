package config

import (
	"fmt"
)

// Constants for validation
const (
	MinPort        = 1
	MaxPort        = 65535
	DefaultTimeout = 30
)

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
