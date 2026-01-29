package entity

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	// Arrange & Act
	user := NewUser("telegram", "user123")

	// Assert
	require.NotEmpty(t, user.ID)
	assert.Equal(t, valueobject.Channel("telegram"), user.Channel)
	assert.Equal(t, "user123", user.ChannelID)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
}

func TestUser_CanAccessSession(t *testing.T) {
	// Arrange
	user := NewUser("telegram", "user123")

	// Act
	result := user.CanAccessSession("session-1")

	// Assert - Currently always returns true per implementation
	assert.True(t, result)
}

func TestUser_GetChannelUserID(t *testing.T) {
	// Arrange
	user := NewUser("discord", "user456")

	// Act
	channelID := user.GetChannelUserID()

	// Assert
	assert.Equal(t, "user456", channelID)
}

func TestUser_IsSameChannel(t *testing.T) {
	// Arrange
	user1 := NewUser("telegram", "user1")
	user2 := NewUser("telegram", "user2")
	user3 := NewUser("discord", "user3")

	// Act & Assert
	assert.True(t, user1.IsSameChannel(user2))
	assert.False(t, user1.IsSameChannel(user3))
}

func TestUser_DifferentChannels(t *testing.T) {
	// Arrange
	user := NewUser("web", "webuser1")

	// Act & Assert
	assert.Equal(t, valueobject.Channel("web"), user.Channel)
	assert.Equal(t, "webuser1", user.ChannelID)
	assert.NotEmpty(t, user.ID)
}
