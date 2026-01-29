package valueobject

import (
	"encoding/json"
	"fmt"
)

// SessionID represents a session identifier.
type SessionID ID

// String returns the string representation of the SessionID.
func (id SessionID) String() string {
	return string(id)
}

// IsEmpty returns true if the SessionID is empty.
func (id SessionID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the SessionID is valid (not empty and matches pattern).
func (id SessionID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the SessionID equals another SessionID.
func (id SessionID) Equals(other SessionID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id SessionID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *SessionID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !SessionID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = SessionID(str)
	return nil
}

// NewSessionID creates a new SessionID from a string.
// Returns an error if the string is not a valid ID.
func NewSessionID(idStr string) (SessionID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return SessionID(id), nil
}

// MustNewSessionID creates a new SessionID from a string.
// Panics if the string is not a valid ID.
func MustNewSessionID(idStr string) SessionID {
	id, err := NewSessionID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
