package valueobject

import (
	"encoding/json"
	"fmt"
)

// LogID represents a log entry identifier.
type LogID ID

// String returns the string representation of the LogID.
func (id LogID) String() string {
	return string(id)
}

// IsEmpty returns true if the LogID is empty.
func (id LogID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the LogID is valid (not empty and matches pattern).
func (id LogID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the LogID equals another LogID.
func (id LogID) Equals(other LogID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id LogID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *LogID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !LogID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = LogID(str)
	return nil
}

// NewLogID creates a new LogID from a string.
// Returns an error if the string is not a valid ID.
func NewLogID(idStr string) (LogID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return LogID(id), nil
}

// MustNewLogID creates a new LogID from a string.
// Panics if the string is not a valid ID.
func MustNewLogID(idStr string) LogID {
	id, err := NewLogID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
