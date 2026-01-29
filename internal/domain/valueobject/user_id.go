package valueobject

import (
	"encoding/json"
	"fmt"
)

// UserID represents a user identifier.
type UserID ID

// String returns the string representation of the UserID.
func (id UserID) String() string {
	return string(id)
}

// IsEmpty returns true if the UserID is empty.
func (id UserID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the UserID is valid (not empty and matches pattern).
func (id UserID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the UserID equals another UserID.
func (id UserID) Equals(other UserID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id UserID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *UserID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !UserID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = UserID(str)
	return nil
}

// NewUserID creates a new UserID from a string.
// Returns an error if the string is not a valid ID.
func NewUserID(idStr string) (UserID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return UserID(id), nil
}

// MustNewUserID creates a new UserID from a string.
// Panics if the string is not a valid ID.
func MustNewUserID(idStr string) UserID {
	id, err := NewUserID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
