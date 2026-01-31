package config

import (
	"fmt"
)

// EventBusConfig represents configuration for the event bus
type EventBusConfig struct {
	// Enabled enables or disables the event bus
	Enabled bool `yaml:"enabled"`

	// BatchSize is the number of events to batch before processing
	BatchSize int `yaml:"batch_size"`

	// FlushIntervalMs is the maximum time in milliseconds to wait before flushing events
	FlushIntervalMs int `yaml:"flush_interval_ms"`

	// EnableLogging enables logging of all events to the database
	EnableLogging bool `yaml:"enable_logging"`

	// BufferSize is the size of the internal event channel buffer
	BufferSize int `yaml:"buffer_size"`
}

// Validate validates the event bus configuration
func (c *EventBusConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	if c.BatchSize <= 0 {
		return fmt.Errorf("event bus batch_size must be positive, got %d", c.BatchSize)
	}

	if c.BatchSize > 10000 {
		return fmt.Errorf("event bus batch_size too large, got %d (max 10000)", c.BatchSize)
	}

	if c.FlushIntervalMs <= 0 {
		return fmt.Errorf("event bus flush_interval_ms must be positive, got %d", c.FlushIntervalMs)
	}

	if c.FlushIntervalMs > 60000 {
		return fmt.Errorf("event bus flush_interval_ms too large, got %d (max 60000)", c.FlushIntervalMs)
	}

	if c.BufferSize <= 0 {
		return fmt.Errorf("event bus buffer_size must be positive, got %d", c.BufferSize)
	}

	if c.BufferSize > 100000 {
		return fmt.Errorf("event bus buffer_size too large, got %d (max 100000)", c.BufferSize)
	}

	return nil
}

// DefaultEventBusConfig returns default event bus configuration
func DefaultEventBusConfig() EventBusConfig {
	return EventBusConfig{
		Enabled:         true,
		BatchSize:       100,
		FlushIntervalMs: 100,
		EnableLogging:   true,
		BufferSize:      1000,
	}
}
