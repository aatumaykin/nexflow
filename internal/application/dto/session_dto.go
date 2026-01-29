package dto

// SessionDTO represents a session data transfer object.
type SessionDTO struct {
	ID        string `json:"id"`         // Unique identifier for the session
	UserID    string `json:"user_id"`    // ID of the user who owns the session
	CreatedAt string `json:"created_at"` // ISO 8601 format timestamp when the session was created
	UpdatedAt string `json:"updated_at"` // ISO 8601 format timestamp when the session was last updated
}

// CreateSessionRequest represents a request to create a new session.
type CreateSessionRequest struct {
	UserID string `json:"user_id" yaml:"user_id"` // ID of the user who will own the session
}

// UpdateSessionRequest represents a request to update an existing session.
type UpdateSessionRequest struct {
	UserID string `json:"user_id,omitempty" yaml:"user_id,omitempty"` // New user ID (optional)
}

// SessionResponse represents a response containing a single session.
type SessionResponse struct {
	Success bool        `json:"success"`           // Whether the operation was successful
	Session *SessionDTO `json:"session,omitempty"` // Session data (if successful)
	Error   string      `json:"error,omitempty"`   // Error message (if failed)
}

// SessionsResponse represents a response containing multiple sessions.
type SessionsResponse struct {
	Success  bool          `json:"success"`            // Whether the operation was successful
	Sessions []*SessionDTO `json:"sessions,omitempty"` // List of sessions (if successful)
	Error    string        `json:"error,omitempty"`    // Error message (if failed)
}
