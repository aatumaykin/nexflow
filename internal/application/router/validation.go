package router

import (
	"context"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// MessageValidator validates incoming messages
type MessageValidator struct {
	config *Config
	logger logging.Logger
}

// NewMessageValidator creates a new message validator
func NewMessageValidator(config *Config, logger logging.Logger) *MessageValidator {
	return &MessageValidator{
		config: config,
		logger: logger,
	}
}

// Validate validates a message and returns an error if validation fails
func (v *MessageValidator) Validate(msg *channels.Message) error {
	if msg == nil {
		return fmt.Errorf("message cannot be nil")
	}

	if v.config.ValidationEnabled {
		return v.validateWithConfig(msg)
	}

	return nil
}

// validateWithConfig validates a message using configuration
func (v *MessageValidator) validateWithConfig(msg *channels.Message) error {
	// Validate UserID
	if msg.UserID == "" {
		v.logger.Warn("message validation failed: empty user ID")
		return NewValidationError("user ID cannot be empty")
	}

	// Validate Content
	if msg.Content == "" {
		v.logger.Warn("message validation failed: empty content", "user_id", msg.UserID)
		return NewValidationError("message content cannot be empty")
	}

	// Validate message length
	if len(msg.Content) > v.config.MaxMessageLength {
		v.logger.Warn("message validation failed: message too long",
			"user_id", msg.UserID,
			"length", len(msg.Content),
			"max_length", v.config.MaxMessageLength,
		)
		return NewValidationError(fmt.Sprintf("message content exceeds maximum length of %d characters", v.config.MaxMessageLength))
	}

	return nil
}

// RetryHandler handles retry logic for operations that may fail
type RetryHandler struct {
	config RetryConfig
	logger logging.Logger
}

// NewRetryHandler creates a new retry handler
func NewRetryHandler(config RetryConfig, logger logging.Logger) *RetryHandler {
	return &RetryHandler{
		config: config,
		logger: logger,
	}
}

// Do executes a function with retry logic
// Returns the last error if all attempts fail
func (r *RetryHandler) Do(ctx context.Context, operation string, fn func() error) error {
	if r.config.MaxAttempts <= 1 {
		// No retry, execute once
		return fn()
	}

	var lastErr error
	delay := r.config.InitialDelay

	for attempt := 1; attempt <= r.config.MaxAttempts; attempt++ {
		err := fn()
		if err == nil {
			// Success
			if attempt > 1 {
				r.logger.Info("operation succeeded after retry",
					"operation", operation,
					"attempt", attempt,
				)
			}
			return nil
		}

		lastErr = err

		// Check if we should retry
		if attempt < r.config.MaxAttempts {
			r.logger.Warn("operation failed, will retry",
				"operation", operation,
				"attempt", attempt,
				"max_attempts", r.config.MaxAttempts,
				"error", err,
				"retry_delay", delay,
			)

			// Wait before retry
			select {
			case <-time.After(delay):
				// Delay elapsed
			case <-ctx.Done():
				// Context cancelled
				return ctx.Err()
			}

			// Calculate next delay with exponential backoff
			delay = time.Duration(float64(delay) * r.config.BackoffMultiplier)
			if delay > r.config.MaxDelay {
				delay = r.config.MaxDelay
			}
		} else {
			// Last attempt failed
			r.logger.Error("operation failed after all retry attempts",
				"operation", operation,
				"attempts", attempt,
				"max_attempts", r.config.MaxAttempts,
				"error", lastErr,
			)
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", r.config.MaxAttempts, lastErr)
}

// IsRetryableError checks if an error is retryable
// This can be extended to check specific error types
func (r *RetryHandler) IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Context cancellation is not retryable
	if err == context.Canceled || err == context.DeadlineExceeded {
		return false
	}

	// Add more specific error type checks here if needed
	// For now, all other errors are considered retryable

	return true
}
