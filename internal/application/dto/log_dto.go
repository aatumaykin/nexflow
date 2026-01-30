package dto

// LogDTO represents a log data transfer object
type LogDTO struct {
	ID        string `json:"id"`
	Level     string `json:"level"`      // "debug", "info", "warn", "error"
	Source    string `json:"source"`     // Source component/module
	Message   string `json:"message"`    // Log message content
	Metadata  string `json:"metadata"`   // Additional metadata in JSON format
	CreatedAt string `json:"created_at"` // ISO 8601 format
}

// CreateLogRequest represents a request to create a log
type CreateLogRequest struct {
	Level    string                 `json:"level" yaml:"level"`   // "debug", "info", "warn", "error"
	Source   string                 `json:"source" yaml:"source"` // Source component/module
	Message  string                 `json:"message" yaml:"message"`
	Metadata map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// LogResponse represents a log response
type LogResponse struct {
	Success bool    `json:"success"`
	Log     *LogDTO `json:"log,omitempty"`
	Error   string  `json:"error,omitempty"`
}

// LogsResponse represents a list of logs response
type LogsResponse struct {
	Success bool      `json:"success"`
	Logs    []*LogDTO `json:"logs,omitempty"`
	Error   string    `json:"error,omitempty"`
}
