package router

import (
	"testing"

	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/logging"
	"github.com/stretchr/testify/assert"
)

func TestConfig_DefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 10000, config.MaxMessageLength)
	assert.True(t, config.ValidationEnabled)
	assert.Equal(t, 3, config.RetryConfig.MaxAttempts)
	assert.Equal(t, int64(100), config.RetryConfig.InitialDelay.Milliseconds())
	assert.Equal(t, int64(5000), config.RetryConfig.MaxDelay.Milliseconds())
	assert.Equal(t, 2.0, config.RetryConfig.BackoffMultiplier)
}

func TestConfig_Validate_DefaultConfig(t *testing.T) {
	config := DefaultConfig()

	err := config.Validate()

	assert.NoError(t, err)
}

func TestConfig_Validate_InvalidMaxMessageLength(t *testing.T) {
	config := &Config{
		MaxMessageLength: -1,
		RetryConfig: RetryConfig{
			MaxAttempts:       1,
			InitialDelay:      100,
			MaxDelay:          200,
			BackoffMultiplier: 2.0,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxMessageLength")
}

func TestConfig_Validate_InvalidMaxAttempts(t *testing.T) {
	config := &Config{
		MaxMessageLength: 1000,
		RetryConfig: RetryConfig{
			MaxAttempts: -1,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxAttempts")
}

func TestConfig_Validate_InvalidInitialDelay(t *testing.T) {
	config := &Config{
		MaxMessageLength: 1000,
		RetryConfig: RetryConfig{
			MaxAttempts:       1,
			InitialDelay:      0,
			MaxDelay:          200,
			BackoffMultiplier: 2.0,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InitialDelay")
}

func TestConfig_Validate_InvalidMaxDelay(t *testing.T) {
	config := &Config{
		MaxMessageLength: 1000,
		RetryConfig: RetryConfig{
			MaxAttempts:       1,
			InitialDelay:      100,
			MaxDelay:          0,
			BackoffMultiplier: 2.0,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxDelay")
}

func TestConfig_Validate_MaxDelayLessThanInitialDelay(t *testing.T) {
	config := &Config{
		MaxMessageLength: 1000,
		RetryConfig: RetryConfig{
			MaxAttempts:       1,
			InitialDelay:      200,
			MaxDelay:          100,
			BackoffMultiplier: 2.0,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "MaxDelay")
}

func TestConfig_Validate_InvalidBackoffMultiplier(t *testing.T) {
	config := &Config{
		MaxMessageLength: 1000,
		RetryConfig: RetryConfig{
			MaxAttempts:       1,
			InitialDelay:      100,
			MaxDelay:          200,
			BackoffMultiplier: 1.0,
		},
	}

	err := config.Validate()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "BackoffMultiplier")
}

func TestMessageValidator_Validate_NilMessage(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	err := validator.Validate(nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestMessageValidator_Validate_EmptyUserID(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "",
		Content: "Hello",
	}

	err := validator.Validate(msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user ID")
}

func TestMessageValidator_Validate_EmptyContent(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "user123",
		Content: "",
	}

	err := validator.Validate(msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message content")
}

func TestMessageValidator_Validate_TooLongContent(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "user123",
		Content: string(make([]byte, 10001)), // 10001 characters, max is 10000
	}

	err := validator.Validate(msg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "exceeds maximum length")
}

func TestMessageValidator_Validate_ValidMessage(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := DefaultConfig()
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "user123",
		Content: "Hello, world!",
	}

	err := validator.Validate(msg)

	assert.NoError(t, err)
}

func TestMessageValidator_Validate_ValidationDisabled(t *testing.T) {
	logger := logging.NewNoopLogger()
	config := &Config{
		ValidationEnabled: false,
	}
	validator := NewMessageValidator(config, logger)

	msg := &channels.Message{
		UserID:  "",
		Content: "",
	}

	err := validator.Validate(msg)

	assert.NoError(t, err)
}
