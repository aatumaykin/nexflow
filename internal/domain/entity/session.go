package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Session represents a conversation session between a user and the AI.
// A session contains all messages exchanged during a conversation.
type Session struct {
	ID        string    `json:"id"`         // Unique identifier for the session
	UserID    string    `json:"user_id"`    // ID of the user who owns this session
	CreatedAt time.Time `json:"created_at"` // Timestamp when the session was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the session was last updated
}

// NewSession creates a new session for the specified user.
func NewSession(userID string) *Session {
	now := utils.Now()
	return &Session{
		ID:        utils.GenerateID(),
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateTimestamp updates the last modified timestamp to the current time.
func (s *Session) UpdateTimestamp() {
	s.UpdatedAt = utils.Now()
}

// IsOwnedBy returns true if the session belongs to the specified user.
func (s *Session) IsOwnedBy(userID string) bool {
	return s.UserID == userID
}
