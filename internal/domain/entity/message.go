package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Message represents a message in a conversation session.
// Messages can be from user, assistant (AI), or system.
type Message struct {
	ID        string    `json:"id"`         // Unique identifier for the message
	SessionID string    `json:"session_id"` // ID of the session this message belongs to
	Role      string    `json:"role"`       // Message role: "user", "assistant", "system"
	Content   string    `json:"content"`    // Message content
	CreatedAt time.Time `json:"created_at"` // Timestamp when the message was created
}

// MessageRole represents the role of the message sender.
type MessageRole string

const (
	RoleUser      MessageRole = "user"      // Message from a human user
	RoleAssistant MessageRole = "assistant" // Message from the AI assistant
	RoleSystem    MessageRole = "system"    // System-level message
)

// NewUserMessage creates a new user message in the specified session.
func NewUserMessage(sessionID, content string) *Message {
	return &Message{
		ID:        utils.GenerateID(),
		SessionID: sessionID,
		Role:      string(RoleUser),
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// NewAssistantMessage creates a new assistant (AI) message in the specified session.
func NewAssistantMessage(sessionID, content string) *Message {
	return &Message{
		ID:        utils.GenerateID(),
		SessionID: sessionID,
		Role:      string(RoleAssistant),
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// NewSystemMessage creates a new system message in the specified session.
func NewSystemMessage(sessionID, content string) *Message {
	return &Message{
		ID:        utils.GenerateID(),
		SessionID: sessionID,
		Role:      string(RoleSystem),
		Content:   content,
		CreatedAt: utils.Now(),
	}
}

// IsFromUser returns true if the message is from a user.
func (m *Message) IsFromUser() bool {
	return m.Role == string(RoleUser)
}

// IsFromAssistant returns true if the message is from the AI assistant.
func (m *Message) IsFromAssistant() bool {
	return m.Role == string(RoleAssistant)
}

// IsSystem returns true if the message is a system message.
func (m *Message) IsSystem() bool {
	return m.Role == string(RoleSystem)
}

// IsPartOfSession returns true if the message belongs to the specified session.
func (m *Message) IsPartOfSession(sessionID string) bool {
	return m.SessionID == sessionID
}
