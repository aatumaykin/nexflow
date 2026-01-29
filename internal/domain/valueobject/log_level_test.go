package valueobject

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogLevel_String(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected string
	}{
		{"debug", LogLevelDebug, "debug"},
		{"info", LogLevelInfo, "info"},
		{"warn", LogLevelWarn, "warn"},
		{"error", LogLevelError, "error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.String())
		})
	}
}

func TestLogLevel_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"debug", LogLevelDebug, true},
		{"info", LogLevelInfo, true},
		{"warn", LogLevelWarn, true},
		{"error", LogLevelError, true},
		{"invalid", LogLevel("invalid"), false},
		{"empty", LogLevel(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.IsValid())
		})
	}
}

func TestLogLevel_IsDebug(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"debug", LogLevelDebug, true},
		{"info", LogLevelInfo, false},
		{"warn", LogLevelWarn, false},
		{"error", LogLevelError, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.IsDebug())
		})
	}
}

func TestLogLevel_IsInfo(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"debug", LogLevelDebug, false},
		{"info", LogLevelInfo, true},
		{"warn", LogLevelWarn, false},
		{"error", LogLevelError, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.IsInfo())
		})
	}
}

func TestLogLevel_IsWarn(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"debug", LogLevelDebug, false},
		{"info", LogLevelInfo, false},
		{"warn", LogLevelWarn, true},
		{"error", LogLevelError, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.IsWarn())
		})
	}
}

func TestLogLevel_IsError(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected bool
	}{
		{"debug", LogLevelDebug, false},
		{"info", LogLevelInfo, false},
		{"warn", LogLevelWarn, false},
		{"error", LogLevelError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.IsError())
		})
	}
}

func TestLogLevel_Priority(t *testing.T) {
	tests := []struct {
		name     string
		level    LogLevel
		expected int
	}{
		{"debug", LogLevelDebug, 0},
		{"info", LogLevelInfo, 1},
		{"warn", LogLevelWarn, 2},
		{"error", LogLevelError, 3},
		{"invalid", LogLevel("invalid"), -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.level.Priority())
		})
	}
}

func TestLogLevel_ShouldLog(t *testing.T) {
	tests := []struct {
		name      string
		level     LogLevel
		threshold LogLevel
		shouldLog bool
	}{
		{"debug at debug threshold", LogLevelDebug, LogLevelDebug, true},
		{"debug at info threshold", LogLevelDebug, LogLevelInfo, false},
		{"info at info threshold", LogLevelInfo, LogLevelInfo, true},
		{"info at warn threshold", LogLevelInfo, LogLevelWarn, false},
		{"warn at warn threshold", LogLevelWarn, LogLevelWarn, true},
		{"warn at error threshold", LogLevelWarn, LogLevelError, false},
		{"error at error threshold", LogLevelError, LogLevelError, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.shouldLog, tt.level.ShouldLog(tt.threshold))
		})
	}
}

func TestLogLevel_MarshalJSON(t *testing.T) {
	level := LogLevelInfo
	data, err := json.Marshal(level)

	require.NoError(t, err)
	assert.JSONEq(t, `"info"`, string(data))
}

func TestLogLevel_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expected  LogLevel
		expectErr bool
	}{
		{"debug", `"debug"`, LogLevelDebug, false},
		{"info", `"info"`, LogLevelInfo, false},
		{"warn", `"warn"`, LogLevelWarn, false},
		{"error", `"error"`, LogLevelError, false},
		{"invalid", `"invalid"`, LogLevel(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var level LogLevel
			err := json.Unmarshal([]byte(tt.data), &level)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, level)
			}
		})
	}
}

func TestNewLogLevel(t *testing.T) {
	tests := []struct {
		name      string
		levelStr  string
		expected  LogLevel
		expectErr bool
	}{
		{"debug", "debug", LogLevelDebug, false},
		{"info", "info", LogLevelInfo, false},
		{"warn", "warn", LogLevelWarn, false},
		{"error", "error", LogLevelError, false},
		{"invalid", "invalid", LogLevel(""), true},
		{"empty", "", LogLevel(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level, err := NewLogLevel(tt.levelStr)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, level)
			}
		})
	}
}

func TestMustNewLogLevel(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewLogLevel("info")
	})

	assert.Panics(t, func() {
		MustNewLogLevel("invalid")
	})
}
