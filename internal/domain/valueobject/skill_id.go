package valueobject

import (
	"encoding/json"
	"fmt"
)

// SkillID represents a skill identifier.
type SkillID ID

// String returns the string representation of the SkillID.
func (id SkillID) String() string {
	return string(id)
}

// IsEmpty returns true if the SkillID is empty.
func (id SkillID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the SkillID is valid (not empty and matches pattern).
func (id SkillID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the SkillID equals another SkillID.
func (id SkillID) Equals(other SkillID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id SkillID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *SkillID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !SkillID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = SkillID(str)
	return nil
}

// NewSkillID creates a new SkillID from a string.
// Returns an error if the string is not a valid ID.
func NewSkillID(idStr string) (SkillID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return SkillID(id), nil
}

// MustNewSkillID creates a new SkillID from a string.
// Panics if the string is not a valid ID.
func MustNewSkillID(idStr string) SkillID {
	id, err := NewSkillID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
