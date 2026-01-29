package entity

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	// Arrange & Act
	session := NewSession("user-1")

	// Assert
	require.NotEmpty(t, session.ID)
	assert.Equal(t, "user-1", string(session.UserID))
	assert.WithinDuration(t, time.Now(), session.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), session.UpdatedAt, time.Second)
}

func TestSession_UpdateTimestamp(t *testing.T) {
	// Arrange
	session := NewSession("user-1")
	initialUpdatedAt := session.UpdatedAt
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	// Act
	session.UpdateTimestamp()

	// Assert
	assert.True(t, session.UpdatedAt.After(initialUpdatedAt))
}

func TestSession_IsOwnedBy(t *testing.T) {
	// Arrange
	session := NewSession("user-1")

	// Act & Assert
	assert.True(t, session.IsOwnedBy(valueobject.UserID("user-1")))
	assert.False(t, session.IsOwnedBy(valueobject.UserID("user-2")))
}

func TestSession_MultipleUsers(t *testing.T) {
	// Arrange
	session1 := NewSession("user-1")
	session2 := NewSession("user-2")
	session3 := NewSession("user-1")

	// Act & Assert
	assert.Equal(t, "user-1", string(session1.UserID))
	assert.Equal(t, "user-2", string(session2.UserID))
	assert.Equal(t, "user-1", string(session3.UserID))

	assert.True(t, session1.IsOwnedBy(valueobject.UserID("user-1")))
	assert.True(t, session3.IsOwnedBy(valueobject.UserID("user-1")))
	assert.False(t, session1.IsOwnedBy(valueobject.UserID("user-2")))
}

func TestSession_UniqueIDs(t *testing.T) {
	// Arrange
	session1 := NewSession("user-1")
	session2 := NewSession("user-1")

	// Act & Assert
	assert.NotEqual(t, session1.ID, session2.ID)
	assert.Equal(t, session1.UserID, session2.UserID)
}
