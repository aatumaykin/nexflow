package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// User represents a user in the system.
// A user can interact through different channels (Telegram, Discord, Web, etc.).
type User struct {
	ID        valueobject.UserID  `json:"id"`         // Unique identifier for the user
	Channel   valueobject.Channel `json:"channel"`    // Channel type: "telegram", "discord", "web", etc.
	ChannelID string              `json:"channel_id"` // Channel-specific user identifier
	CreatedAt time.Time           `json:"created_at"` // Timestamp when the user was created
}

// NewUser creates a new user with the specified channel and channel ID.
func NewUser(channel, channelID string) *User {
	return &User{
		ID:        valueobject.UserID(utils.GenerateID()),
		Channel:   valueobject.MustNewChannel(channel),
		ChannelID: channelID,
		CreatedAt: utils.Now(),
	}
}

// CanAccessSession checks if the user can access the specified session.
//
// DEPRECATED: Use service.SessionAccessService for proper access control.
// This method is kept for backward compatibility and will be removed in future versions.
// SessionAccessService implements proper ownership checks against the database.
//
// For session access control, use:
//
//	service.SessionAccessService.CanAccessSession(ctx, userID, sessionID)
func (u *User) CanAccessSession(sessionID valueobject.SessionID) bool {
	// DEPRECATED: Always returns true for backward compatibility.
	// Use service.SessionAccessService for actual access control.
	return true
}

// GetChannelUserID returns the channel-specific user identifier.
func (u *User) GetChannelUserID() string {
	return u.ChannelID
}

// IsSameChannel returns true if the user is from the same channel as the other user.
func (u *User) IsSameChannel(other *User) bool {
	return u.Channel.Equals(other.Channel)
}
