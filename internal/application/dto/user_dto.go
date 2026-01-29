package dto

// UserDTO represents a user data transfer object.
type UserDTO struct {
	ID        string `json:"id"`         // Unique identifier for the user
	Channel   string `json:"channel"`    // Channel type: "telegram", "discord", "web", etc.
	ChannelID string `json:"channel_id"` // Channel-specific user identifier
	CreatedAt string `json:"created_at"` // ISO 8601 format timestamp when the user was created
}

// CreateUserRequest represents a request to create a new user.
type CreateUserRequest struct {
	Channel   string `json:"channel" yaml:"channel"`       // Channel type
	ChannelID string `json:"channel_id" yaml:"channel_id"` // Channel-specific user identifier
}

// UpdateUserRequest represents a request to update an existing user.
type UpdateUserRequest struct {
	ChannelID string `json:"channel_id,omitempty" yaml:"channel_id,omitempty"` // New channel ID (optional)
}

// UserResponse represents a response containing a single user.
type UserResponse struct {
	Success bool     `json:"success"`         // Whether the operation was successful
	User    *UserDTO `json:"user,omitempty"`  // User data (if successful)
	Error   string   `json:"error,omitempty"` // Error message (if failed)
}

// UsersResponse represents a response containing multiple users.
type UsersResponse struct {
	Success bool       `json:"success"`         // Whether the operation was successful
	Users   []*UserDTO `json:"users,omitempty"` // List of users (if successful)
	Error   string     `json:"error,omitempty"` // Error message (if failed)
}
