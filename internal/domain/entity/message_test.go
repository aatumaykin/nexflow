package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUserMessage(t *testing.T) {
	// Arrange & Act
	message := NewUserMessage("session-1", "Hello, world!")

	// Assert
	require.NotEmpty(t, message.ID)
	assert.Equal(t, "session-1", message.SessionID)
	assert.Equal(t, string(RoleUser), message.Role)
	assert.Equal(t, "Hello, world!", message.Content)
	assert.WithinDuration(t, time.Now(), message.CreatedAt, time.Second)
}

func TestNewAssistantMessage(t *testing.T) {
	// Arrange & Act
	message := NewAssistantMessage("session-1", "Hi there!")

	// Assert
	require.NotEmpty(t, message.ID)
	assert.Equal(t, "session-1", message.SessionID)
	assert.Equal(t, string(RoleAssistant), message.Role)
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
		ID:        "msg-1",
		SessionID: "session-1",
		Role:      string(RoleSystem),
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
	assert.True(t, message.IsPartOfSession("session-1"))
	assert.False(t, message.IsPartOfSession("session-2"))
}

func TestMessage_DifferentRoles(t *testing.T) {
	// Arrange
	userMsg := NewUserMessage("session-1", "user")
	assistantMsg := NewAssistantMessage("session-1", "assistant")
	systemMsg := &Message{
		ID:        "msg-3",
		SessionID: "session-1",
		Role:      string(RoleSystem),
		Content:   "system",
		CreatedAt: time.Now(),
	}

	// Act & Assert
	assert.Equal(t, string(RoleUser), userMsg.Role)
	assert.Equal(t, string(RoleAssistant), assistantMsg.Role)
	assert.Equal(t, string(RoleSystem), systemMsg.Role)
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

	assert.True(t, msgs[0].IsPartOfSession(sessionID))
	assert.True(t, msgs[1].IsPartOfSession(sessionID))
	assert.True(t, msgs[2].IsPartOfSession(sessionID))
}

func TestMessage_UniqueIDs(t *testing.T) {
	// Arrange
	msg1 := NewUserMessage("session-1", "test")
	msg2 := NewUserMessage("session-1", "test")

	// Act & Assert
	assert.NotEqual(t, msg1.ID, msg2.ID)
}
