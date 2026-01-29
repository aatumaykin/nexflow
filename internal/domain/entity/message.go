package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Message represents a message in a conversation session.
// Messages can be from user, assistant (AI), or system.
type Message struct {
	ID        valueobject.MessageID   `json:"id"`         // Unique identifier for the message
	SessionID valueobject.SessionID   `json:"session_id"` // ID of the session this message belongs to
	Role      valueobject.MessageRole `json:"role"`       // Message role: "user", "assistant", "system"
	Content   string                  `json:"content"`    // Message content
	CreatedAt time.Time               `json:"created_at"` // Timestamp when the message was created
}

// NewUserMessage creates a new user message in the specified session.
func NewUserMessage(sessionID, content string) *Message {
	return &Message{
		ID:        valueobject.MessageID(utils.GenerateID()),
		SessionID: valueobject.MustNewSessionID(sessionID),
		Role:      valueobject.RoleUser,
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// NewAssistantMessage creates a new assistant (AI) message in the specified session.
func NewAssistantMessage(sessionID, content string) *Message {
	return &Message{
		ID:        valueobject.MessageID(utils.GenerateID()),
		SessionID: valueobject.MustNewSessionID(sessionID),
		Role:      valueobject.RoleAssistant,
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// NewSystemMessage creates a new system message in the specified session.
func NewSystemMessage(sessionID, content string) *Message {
	return &Message{
		ID:        valueobject.MessageID(utils.GenerateID()),
		SessionID: valueobject.MustNewSessionID(sessionID),
		Role:      valueobject.RoleSystem,
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// IsFromUser returns true if the message is from a user.
func (m *Message) IsFromUser() bool {
	return m.Role == valueobject.RoleUser
}

// IsFromAssistant returns true if the message is from the AI assistant.
func (m *Message) IsFromAssistant() bool {
	return m.Role == valueobject.RoleAssistant
}

// IsSystem returns true if the message is a system message.
func (m *Message) IsSystem() bool {
	return m.Role == valueobject.RoleSystem
}

// IsPartOfSession returns true if the message belongs to the specified session.
func (m *Message) IsPartOfSession(sessionID valueobject.SessionID) bool {
	return m.SessionID.Equals(sessionID)
}
