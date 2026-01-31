package config

import (
	"fmt"
)

// RouterConfig represents configuration for message router
type RouterConfig struct {
	// MaxMessageLength is the maximum allowed length of a message content
	MaxMessageLength int `yaml:"max_message_length"`

	// ValidationEnabled enables or disables message validation
	ValidationEnabled bool `yaml:"validation_enabled"`

	// RetryMaxAttempts is the maximum number of retry attempts for failed operations
	RetryMaxAttempts int `yaml:"retry_max_attempts"`

	// RetryInitialDelayMs is the initial delay before first retry in milliseconds
	RetryInitialDelayMs int `yaml:"retry_initial_delay_ms"`

	// RetryMaxDelayMs is the maximum delay between retries in milliseconds
	RetryMaxDelayMs int `yaml:"retry_max_delay_ms"`

	// RetryBackoffMultiplier is the multiplier for exponential backoff
	RetryBackoffMultiplier float64 `yaml:"retry_backoff_multiplier"`
}

// Validate validates the router configuration
func (c *RouterConfig) Validate() error {
	if c.MaxMessageLength < 0 {
		return fmt.Errorf("router max_message_length must be positive, got %d", c.MaxMessageLength)
	}

	if c.MaxMessageLength > 100000 {
		return fmt.Errorf("router max_message_length too large, got %d (max 100000)", c.MaxMessageLength)
	}

	if c.RetryMaxAttempts < 0 {
		return fmt.Errorf("router retry_max_attempts must be non-negative, got %d", c.RetryMaxAttempts)
	}

	if c.RetryMaxAttempts == 0 {
		return nil // Skip other validation if retries are disabled
	}

	if c.RetryInitialDelayMs <= 0 {
		return fmt.Errorf("router retry_initial_delay_ms must be positive, got %d", c.RetryInitialDelayMs)
	}

	if c.RetryMaxDelayMs <= 0 {
		return fmt.Errorf("router retry_max_delay_ms must be positive, got %d", c.RetryMaxDelayMs)
	}

	if c.RetryMaxDelayMs < c.RetryInitialDelayMs {
		return fmt.Errorf("router retry_max_delay_ms must be greater than or equal to retry_initial_delay_ms")
	}

	if c.RetryBackoffMultiplier <= 1.0 {
		return fmt.Errorf("router retry_backoff_multiplier must be greater than 1.0, got %f", c.RetryBackoffMultiplier)
	}

	return nil
}

// DefaultRouterConfig returns default router configuration
func DefaultRouterConfig() RouterConfig {
	return RouterConfig{
		MaxMessageLength:       10000,
		ValidationEnabled:      true,
		RetryMaxAttempts:       3,
		RetryInitialDelayMs:    100,
		RetryMaxDelayMs:        5000,
		RetryBackoffMultiplier: 2.0,
	}
}

// ToRouterConfig converts this configuration to router.Config
// This returns a map representation that can be used to create router.Config
func (c *RouterConfig) ToRouterConfig() map[string]interface{} {
	return map[string]interface{}{
		"max_message_length":       c.MaxMessageLength,
		"validation_enabled":       c.ValidationEnabled,
		"retry_max_attempts":       c.RetryMaxAttempts,
		"retry_initial_delay_ms":   c.RetryInitialDelayMs,
		"retry_max_delay_ms":       c.RetryMaxDelayMs,
		"retry_backoff_multiplier": c.RetryBackoffMultiplier,
	}
}
