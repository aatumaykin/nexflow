package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidLogLevel is returned when an invalid log level is provided.
	ErrInvalidLogLevel = errors.New("invalid log level")
)

// LogLevel represents the severity level of a log message.
// It's a value object that ensures type safety for log levels.
type LogLevel string

const (
	// LogLevelDebug represents the debug level for detailed information.
	LogLevelDebug LogLevel = "debug"
	// LogLevelInfo represents the info level for general information.
	LogLevelInfo LogLevel = "info"
	// LogLevelWarn represents the warning level for potential issues.
	LogLevelWarn LogLevel = "warn"
	// LogLevelError represents the error level for errors and failures.
	LogLevelError LogLevel = "error"
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	return string(l)
}

// IsValid checks if the log level is valid.
func (l LogLevel) IsValid() bool {
	switch l {
	case LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError:
		return true
	default:
		return false
	}
}

// IsDebug returns true if the log level is debug.
func (l LogLevel) IsDebug() bool {
	return l == LogLevelDebug
}

// IsInfo returns true if the log level is info.
func (l LogLevel) IsInfo() bool {
	return l == LogLevelInfo
}

// IsWarn returns true if the log level is warn.
func (l LogLevel) IsWarn() bool {
	return l == LogLevelWarn
}

// IsError returns true if the log level is error.
func (l LogLevel) IsError() bool {
	return l == LogLevelError
}

// Priority returns the numeric priority of the log level.
// Higher values indicate higher severity (error=3, warn=2, info=1, debug=0).
func (l LogLevel) Priority() int {
	switch l {
	case LogLevelDebug:
		return 0
	case LogLevelInfo:
		return 1
	case LogLevelWarn:
		return 2
	case LogLevelError:
		return 3
	default:
		return -1
	}
}

// ShouldLog returns true if the receiver log level should be logged
// when comparing against the threshold level.
func (l LogLevel) ShouldLog(threshold LogLevel) bool {
	return l.Priority() >= threshold.Priority()
}

// MarshalJSON implements json.Marshaler interface.
func (l LogLevel) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(l))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (l *LogLevel) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*l = LogLevel(str)
	if !l.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidLogLevel, str)
	}
	return nil
}

// NewLogLevel creates a new LogLevel from a string.
// Returns an error if the string is not a valid log level.
func NewLogLevel(level string) (LogLevel, error) {
	l := LogLevel(level)
	if !l.IsValid() {
		return "", ErrInvalidLogLevel
	}
	return l, nil
}

// MustNewLogLevel creates a new LogLevel from a string.
// Panics if the string is not a valid log level.
func MustNewLogLevel(level string) LogLevel {
	l, err := NewLogLevel(level)
	if err != nil {
		panic(err)
	}
	return l
}
