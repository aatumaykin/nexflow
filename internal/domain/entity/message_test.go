package entity

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserMessage(t *testing.T) {
	// Arrange & Act
	message := NewUserMessage("session-1", "Hello, world!")

	// Assert
	require.NotEmpty(t, message.ID)
	assert.Equal(t, "session-1", string(message.SessionID))
	assert.Equal(t, valueobject.RoleUser, message.Role)
	assert.Equal(t, "Hello, world!", message.Content)
	assert.WithinDuration(t, time.Now(), message.CreatedAt, time.Second)
}

func TestNewAssistantMessage(t *testing.T) {
	// Arrange & Act
	message := NewAssistantMessage("session-1", "Hi there!")

	// Assert
	require.NotEmpty(t, message.ID)
	assert.Equal(t, "session-1", string(message.SessionID))
	assert.Equal(t, valueobject.RoleAssistant, message.Role)
	assert.Equal(t, "Hi there!", message.Content)
	assert.WithinDuration(t, time.Now(), message.CreatedAt, time.Second)
}

func TestMessage_IsFromUser(t *testing.T) {
	// Arrange
	message := NewUserMessage("session-1", "test")

	// Act & Assert
	assert.True(t, message.IsFromUser())
	assert.False(t, message.IsFromAssistant())
	assert.False(t, message.IsSystem())
}

func TestMessage_IsFromAssistant(t *testing.T) {
	// Arrange
	message := NewAssistantMessage("session-1", "test")

	// Act & Assert
	assert.False(t, message.IsFromUser())
	assert.True(t, message.IsFromAssistant())
	assert.False(t, message.IsSystem())
}

func TestMessage_IsSystem(t *testing.T) {
	// Arrange
	message := &Message{
		ID:        valueobject.MessageID("msg-1"),
		SessionID: valueobject.SessionID("session-1"),
		Role:      valueobject.RoleSystem,
		Content:   "System message",
		CreatedAt: time.Now(),
	}

	// Act & Assert
	assert.False(t, message.IsFromUser())
	assert.False(t, message.IsFromAssistant())
	assert.True(t, message.IsSystem())
}

func TestMessage_IsPartOfSession(t *testing.T) {
	// Arrange
	message := NewUserMessage("session-1", "test")

	// Act & Assert
	assert.True(t, message.IsPartOfSession(valueobject.SessionID("session-1")))
	assert.False(t, message.IsPartOfSession(valueobject.SessionID("session-2")))
}

func TestMessage_DifferentRoles(t *testing.T) {
	// Arrange
	userMsg := NewUserMessage("session-1", "user")
	assistantMsg := NewAssistantMessage("session-1", "assistant")
	systemMsg := &Message{
		ID:        valueobject.MessageID("msg-3"),
		SessionID: valueobject.SessionID("session-1"),
		Role:      valueobject.RoleSystem,
		Content:   "system",
		CreatedAt: time.Now(),
	}

	// Act & Assert
	assert.Equal(t, valueobject.RoleUser, userMsg.Role)
	assert.Equal(t, valueobject.RoleAssistant, assistantMsg.Role)
	assert.Equal(t, valueobject.RoleSystem, systemMsg.Role)
}

func TestMessage_MultipleMessages(t *testing.T) {
	// Arrange
	sessionID := "session-1"
	msgs := []*Message{
		NewUserMessage(sessionID, "message 1"),
		NewAssistantMessage(sessionID, "message 2"),
		NewUserMessage(sessionID, "message 3"),
	}

	// Act & Assert
	assert.Equal(t, 3, len(msgs))
	assert.True(t, msgs[0].IsFromUser())
	assert.True(t, msgs[1].IsFromAssistant())
	assert.True(t, msgs[2].IsFromUser())

	assert.True(t, msgs[0].IsPartOfSession(valueobject.SessionID(sessionID)))
	assert.True(t, msgs[1].IsPartOfSession(valueobject.SessionID(sessionID)))
	assert.True(t, msgs[2].IsPartOfSession(valueobject.SessionID(sessionID)))
}

func TestMessage_UniqueIDs(t *testing.T) {
	// Arrange
	msg1 := NewUserMessage("session-1", "test")
	msg2 := NewUserMessage("session-1", "test")

	// Act & Assert
	assert.NotEqual(t, msg1.ID, msg2.ID)
}
