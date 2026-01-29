package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Log represents an application log entry.
// Logs are stored in the database for observability and debugging.
type Log struct {
	ID        valueobject.LogID    `json:"id"`         // Unique identifier for the log entry
	Level     valueobject.LogLevel `json:"level"`      // Log level: "debug", "info", "warn", "error"
	Source    string               `json:"source"`     // Source component/module that generated the log
	Message   string               `json:"message"`    // Log message content
	Metadata  string               `json:"metadata"`   // Additional metadata in JSON format
	CreatedAt time.Time            `json:"created_at"` // Timestamp when the log was created
}

// NewLog creates a new log entry with the specified level, source, message, and metadata.
// The log entry is assigned a unique ID and the current timestamp.
func NewLog(level valueobject.LogLevel, source, message string, metadata map[string]interface{}) *Log {
	return &Log{
		ID:        valueobject.LogID(utils.GenerateID()),
		Level:     level,
		Source:    source,
		Message:   message,
		Metadata:  utils.MarshalJSON(metadata),
		CreatedAt: utils.Now(),
	}
}

// IsDebug returns true if the log is at debug level.
func (l *Log) IsDebug() bool {
	return l.Level == valueobject.LogLevelDebug
}

// IsInfo returns true if the log is at info level.
func (l *Log) IsInfo() bool {
	return l.Level == valueobject.LogLevelInfo
}

// IsWarn returns true if the log is at warn level.
func (l *Log) IsWarn() bool {
	return l.Level == valueobject.LogLevelWarn
}

// IsError returns true if the log is at error level.
func (l *Log) IsError() bool {
	return l.Level == valueobject.LogLevelError
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
