package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Log represents an application log entry.
// Logs are stored in the database for observability and debugging.
type Log struct {
	ID        string    `json:"id"`         // Unique identifier for the log entry
	Level     string    `json:"level"`      // Log level: "debug", "info", "warn", "error"
	Source    string    `json:"source"`     // Source component/module that generated the log
	Message   string    `json:"message"`    // Log message content
	Metadata  string    `json:"metadata"`   // Additional metadata in JSON format
	CreatedAt time.Time `json:"created_at"` // Timestamp when the log was created
}

// LogLevel represents the severity level of a log message.
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug" // Debug level for detailed information
	LogLevelInfo  LogLevel = "info"  // Info level for general information
	LogLevelWarn  LogLevel = "warn"  // Warning level for potential issues
	LogLevelError LogLevel = "error" // Error level for errors and failures
)

// NewLog creates a new log entry with the specified level, source, message, and metadata.
// The log entry is assigned a unique ID and the current timestamp.
func NewLog(level LogLevel, source, message string, metadata map[string]interface{}) *Log {
	return &Log{
		ID:        utils.GenerateID(),
		Level:     string(level),
		Source:    source,
		Message:   message,
		Metadata:  utils.MarshalJSON(metadata),
		CreatedAt: utils.Now(),
	}
}

// IsDebug returns true if the log is at debug level.
func (l *Log) IsDebug() bool {
	return l.Level == string(LogLevelDebug)
}

// IsInfo returns true if the log is at info level.
func (l *Log) IsInfo() bool {
	return l.Level == string(LogLevelInfo)
}

// IsWarn returns true if the log is at warn level.
func (l *Log) IsWarn() bool {
	return l.Level == string(LogLevelWarn)
}

// IsError returns true if the log is at error level.
func (l *Log) IsError() bool {
	return l.Level == string(LogLevelError)
}

// IsFromSource returns true if the log originated from the specified source.
func (l *Log) IsFromSource(source string) bool {
	return l.Source == source
}

// GetMetadata parses and returns the metadata as a map.
// Returns nil if parsing fails or metadata is empty.
func (l *Log) GetMetadata() map[string]interface{} {
	return utils.UnmarshalJSONToMap(l.Metadata)
}
