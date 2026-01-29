package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Session represents a conversation session between a user and the AI.
// A session contains all messages exchanged during a conversation.
type Session struct {
	ID        valueobject.SessionID `json:"id"`         // Unique identifier for the session
	UserID    valueobject.UserID    `json:"user_id"`    // ID of the user who owns this session
	CreatedAt time.Time             `json:"created_at"` // Timestamp when the session was created
	UpdatedAt time.Time             `json:"updated_at"` // Timestamp when the session was last updated
}

// NewSession creates a new session for the specified user.
func NewSession(userID string) *Session {
	now := utils.Now()
	return &Session{
		ID:        valueobject.SessionID(utils.GenerateID()),
		UserID:    valueobject.MustNewUserID(userID),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UpdateTimestamp updates the last modified timestamp to the current time.
func (s *Session) UpdateTimestamp() {
	s.UpdatedAt = utils.Now()
}

// IsOwnedBy returns true if the session belongs to the specified user.
func (s *Session) IsOwnedBy(userID valueobject.UserID) bool {
	return s.UserID.Equals(userID)
}
