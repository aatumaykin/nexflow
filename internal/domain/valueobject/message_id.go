package valueobject

import (
	"encoding/json"
	"fmt"
)

// MessageID represents a message identifier.
type MessageID ID

// String returns the string representation of the MessageID.
func (id MessageID) String() string {
	return string(id)
}

// IsEmpty returns true if the MessageID is empty.
func (id MessageID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the MessageID is valid (not empty and matches pattern).
func (id MessageID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the MessageID equals another MessageID.
func (id MessageID) Equals(other MessageID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id MessageID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *MessageID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !MessageID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = MessageID(str)
	return nil
}

// NewMessageID creates a new MessageID from a string.
// Returns an error if the string is not a valid ID.
func NewMessageID(idStr string) (MessageID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return MessageID(id), nil
}

// MustNewMessageID creates a new MessageID from a string.
// Panics if the string is not a valid ID.
func MustNewMessageID(idStr string) MessageID {
	id, err := NewMessageID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
