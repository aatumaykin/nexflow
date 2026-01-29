package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidMessageRole is returned when an invalid message role is provided.
	ErrInvalidMessageRole = errors.New("invalid message role")
)

// MessageRole represents the role of the message sender.
// It's a value object that ensures type safety for message roles.
type MessageRole string

const (
	// RoleUser represents a message from a human user.
	RoleUser MessageRole = "user"
	// RoleAssistant represents a message from the AI assistant.
	RoleAssistant MessageRole = "assistant"
	// RoleSystem represents a system-level message.
	RoleSystem MessageRole = "system"
)

// String returns the string representation of the message role.
func (r MessageRole) String() string {
	return string(r)
}

// IsValid checks if the message role is valid.
func (r MessageRole) IsValid() bool {
	switch r {
	case RoleUser, RoleAssistant, RoleSystem:
		return true
	default:
		return false
	}
}

// IsUser returns true if the message role is user.
func (r MessageRole) IsUser() bool {
	return r == RoleUser
}

// IsAssistant returns true if the message role is assistant.
func (r MessageRole) IsAssistant() bool {
	return r == RoleAssistant
}

// IsSystem returns true if the message role is system.
func (r MessageRole) IsSystem() bool {
	return r == RoleSystem
}

// MarshalJSON implements json.Marshaler interface.
func (r MessageRole) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(r))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (r *MessageRole) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*r = MessageRole(str)
	if !r.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidMessageRole, str)
	}
	return nil
}

// NewMessageRole creates a new MessageRole from a string.
// Returns an error if the string is not a valid message role.
func NewMessageRole(role string) (MessageRole, error) {
	r := MessageRole(role)
	if !r.IsValid() {
		return "", ErrInvalidMessageRole
	}
	return r, nil
}

// MustNewMessageRole creates a new MessageRole from a string.
// Panics if the string is not a valid message role.
func MustNewMessageRole(role string) MessageRole {
	r, err := NewMessageRole(role)
	if err != nil {
		panic(err)
	}
	return r
}
