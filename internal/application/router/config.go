package router

import (
	"fmt"
	"time"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error for %s: %s", e.Field, e.Message)
}

// NewValidationError creates a new ValidationError
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		Message: message,
	}
}

// Config holds configuration for MessageRouter
type Config struct {
	// MaxMessageLength is the maximum allowed length of a message content
	MaxMessageLength int

	// RetryConfig holds retry configuration for failed operations
	RetryConfig RetryConfig

	// ValidationEnabled enables/disables message validation
	ValidationEnabled bool
}

// RetryConfig holds retry configuration
type RetryConfig struct {
	// MaxAttempts is the maximum number of retry attempts
	MaxAttempts int

	// InitialDelay is the initial delay before the first retry
	InitialDelay time.Duration

	// MaxDelay is the maximum delay between retries
	MaxDelay time.Duration

	// BackoffMultiplier is the multiplier for exponential backoff
	BackoffMultiplier float64
}

// DefaultConfig returns the default configuration for MessageRouter
func DefaultConfig() *Config {
	return &Config{
		MaxMessageLength:  10000, // 10,000 characters
		ValidationEnabled: true,
		RetryConfig: RetryConfig{
			MaxAttempts:       3,
			InitialDelay:      100 * time.Millisecond,
			MaxDelay:          5 * time.Second,
			BackoffMultiplier: 2.0,
		},
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.MaxMessageLength <= 0 {
		return NewValidationError("MaxMessageLength must be positive")
	}

	if c.RetryConfig.MaxAttempts < 0 {
		return NewValidationError("MaxAttempts must be non-negative")
	}

	if c.RetryConfig.MaxAttempts == 0 {
		return nil // Skip other validation if retries are disabled
	}

	if c.RetryConfig.InitialDelay <= 0 {
		return NewValidationError("InitialDelay must be positive")
	}

	if c.RetryConfig.MaxDelay <= 0 {
		return NewValidationError("MaxDelay must be positive")
	}

	if c.RetryConfig.MaxDelay < c.RetryConfig.InitialDelay {
		return NewValidationError("MaxDelay must be greater than or equal to InitialDelay")
	}

	if c.RetryConfig.BackoffMultiplier <= 1.0 {
		return NewValidationError("BackoffMultiplier must be greater than 1.0")
	}

	return nil
}
