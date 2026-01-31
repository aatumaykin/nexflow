package router

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// TestMessageValidator_ValidateNilMessage tests validation of nil message
func TestMessageValidator_ValidateNilMessage(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	err := validator.Validate(nil)

	if err == nil {
		t.Error("Expected error for nil message")
	}

	expectedErr := "message cannot be nil"
	if err != nil && err.Error() != expectedErr {
		t.Errorf("Expected error '%s', got '%s'", expectedErr, err.Error())
	}
}

// TestMessageValidator_ValidateEmptyUserID tests validation of empty user ID
func TestMessageValidator_ValidateEmptyUserID(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "",
		Content: "Hello",
	}

	err := validator.Validate(msg)

	if err == nil {
		t.Error("Expected error for empty user ID")
	}

	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

// TestMessageValidator_ValidateEmptyContent tests validation of empty content
func TestMessageValidator_ValidateEmptyContent(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "user123",
		Content: "",
	}

	err := validator.Validate(msg)

	if err == nil {
		t.Error("Expected error for empty content")
	}

	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

// TestMessageValidator_ValidateMessageTooLong tests validation of message exceeding max length
func TestMessageValidator_ValidateMessageTooLong(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.MaxMessageLength = 100
	validator := NewMessageValidator(config, logger)

	// Create a message with 101 characters
	longContent := ""
	for i := 0; i < 101; i++ {
		longContent += "a"
	}

	msg := &channels.Message{
		UserID:  "user123",
		Content: longContent,
	}

	err := validator.Validate(msg)

	if err == nil {
		t.Error("Expected error for message too long")
	}

	if _, ok := err.(*ValidationError); !ok {
		t.Errorf("Expected ValidationError, got %T", err)
	}
}

// TestMessageValidator_ValidateValidMessage tests validation of valid message
func TestMessageValidator_ValidateValidMessage(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "user123",
		Content: "Hello, world!",
	}

	err := validator.Validate(msg)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestMessageValidator_ValidateDisabled tests validation when disabled
func TestMessageValidator_ValidateDisabled(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.ValidationEnabled = false
	validator := NewMessageValidator(config, logger)

	// Invalid message should pass when validation is disabled
	msg := &channels.Message{
		UserID:  "",
		Content: "",
	}

	err := validator.Validate(msg)

	if err != nil {
		t.Errorf("Expected no error when validation disabled, got: %v", err)
	}
}

// TestMessageValidator_ValidateExactMaxLength tests validation of message at exact max length
func TestMessageValidator_ValidateExactMaxLength(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.MaxMessageLength = 100
	validator := NewMessageValidator(config, logger)

	// Create a message with exactly 100 characters
	exactLengthContent := ""
	for i := 0; i < 100; i++ {
		exactLengthContent += "a"
	}

	msg := &channels.Message{
		UserID:  "user123",
		Content: exactLengthContent,
	}

	err := validator.Validate(msg)

	if err != nil {
		t.Errorf("Expected no error for message at exact max length, got: %v", err)
	}
}

// TestRetryHandler_DoSuccess tests successful operation on first attempt
func TestRetryHandler_DoSuccess(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attempts := 0
	err := retryHandler.Do(ctx, "test_operation", func() error {
		attempts++
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

// TestRetryHandler_DoRetrySuccess tests operation succeeds after retries
func TestRetryHandler_DoRetrySuccess(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 3
	config.RetryConfig.InitialDelay = 50 * time.Millisecond
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attempts := 0
	err := retryHandler.Do(ctx, "test_operation", func() error {
		attempts++
		if attempts < 3 {
			return errors.New("temporary error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error after retries, got: %v", err)
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}
}

// TestRetryHandler_DoAllFail tests operation fails after all attempts
func TestRetryHandler_DoAllFail(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 3
	config.RetryConfig.InitialDelay = 50 * time.Millisecond
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attempts := 0
	err := retryHandler.Do(ctx, "test_operation", func() error {
		attempts++
		return errors.New("permanent error")
	})

	if err == nil {
		t.Error("Expected error after all attempts failed")
	}

	if attempts != 3 {
		t.Errorf("Expected 3 attempts, got %d", attempts)
	}

	expectedErrMsg := "operation failed after 3 attempts"
	if err != nil && err.Error()[:len(expectedErrMsg)] != expectedErrMsg {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrMsg, err.Error())
	}
}

// TestRetryHandler_DoNoRetry tests single attempt when MaxAttempts is 1
func TestRetryHandler_DoNoRetry(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 1
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attempts := 0
	err := retryHandler.Do(ctx, "test_operation", func() error {
		attempts++
		return errors.New("error")
	})

	if err == nil {
		t.Error("Expected error")
	}

	if attempts != 1 {
		t.Errorf("Expected 1 attempt, got %d", attempts)
	}
}

// TestRetryHandler_ContextCancellation tests context cancellation during retry
func TestRetryHandler_ContextCancellation(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 10
	config.RetryConfig.InitialDelay = 100 * time.Millisecond
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	// Cancel context after first attempt
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := retryHandler.Do(ctx, "test_operation", func() error {
		attempts++
		return errors.New("error")
	})

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got: %v", err)
	}

	// Should have been cancelled after first attempt
	if attempts > 2 {
		t.Errorf("Expected at most 2 attempts before cancellation, got %d", attempts)
	}
}

// TestRetryHandler_ExponentialBackoff tests exponential backoff calculation
func TestRetryHandler_ExponentialBackoff(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 5
	config.RetryConfig.InitialDelay = 100 * time.Millisecond
	config.RetryConfig.BackoffMultiplier = 2.0
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attemptTimes := make([]time.Time, 0)

	err := retryHandler.Do(ctx, "test_operation", func() error {
		attemptTimes = append(attemptTimes, time.Now())
		if len(attemptTimes) < 5 {
			return errors.New("error")
		}
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify exponential backoff
	if len(attemptTimes) >= 3 {
		delay1 := attemptTimes[1].Sub(attemptTimes[0])
		delay2 := attemptTimes[2].Sub(attemptTimes[1])

		// Second delay should be approximately 2x the first delay
		ratio := float64(delay2) / float64(delay1)
		if ratio < 1.8 || ratio > 2.2 {
			t.Errorf("Expected exponential backoff ratio around 2.0, got %.2f", ratio)
		}
	}
}

// TestRetryHandler_MaxDelayCaps tests that delay is capped at MaxDelay
func TestRetryHandler_MaxDelayCaps(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	config.RetryConfig.MaxAttempts = 10
	config.RetryConfig.InitialDelay = 10 * time.Millisecond
	config.RetryConfig.MaxDelay = 50 * time.Millisecond
	config.RetryConfig.BackoffMultiplier = 10.0
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	ctx := context.Background()
	attemptTimes := make([]time.Time, 0)

	err := retryHandler.Do(ctx, "test_operation", func() error {
		attemptTimes = append(attemptTimes, time.Now())
		if len(attemptTimes) < 10 {
			return errors.New("error")
		}
		return nil
	})

	// Verify no delay exceeds MaxDelay
	for i := 1; i < len(attemptTimes); i++ {
		delay := attemptTimes[i].Sub(attemptTimes[i-1])
		if delay > config.RetryConfig.MaxDelay+10*time.Millisecond {
			t.Errorf("Delay %d (%v) exceeds MaxDelay (%v)", i, delay, config.RetryConfig.MaxDelay)
		}
	}

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
}

// TestRetryHandler_IsRetryableError tests IsRetryableError method
func TestRetryHandler_IsRetryableError(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "context canceled",
			err:      context.Canceled,
			expected: false,
		},
		{
			name:     "context deadline exceeded",
			err:      context.DeadlineExceeded,
			expected: false,
		},
		{
			name:     "generic error",
			err:      errors.New("some error"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := retryHandler.IsRetryableError(tt.err)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestNewValidationError tests creation of ValidationError
func TestNewValidationError(t *testing.T) {
	msg := "test error"
	validationErr := NewValidationError(msg)

	if validationErr == nil {
		t.Fatal("Expected non-nil error")
	}

	if validationErr.Message != msg {
		t.Errorf("Expected message '%s', got '%s'", msg, validationErr.Message)
	}
}

// TestValidationError_Error tests error message format
func TestValidationError_Error(t *testing.T) {
	field := "content"
	message := "cannot be empty"
	err := &ValidationError{
		Field:   field,
		Message: message,
	}

	expected := "validation error for content: cannot be empty"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}
