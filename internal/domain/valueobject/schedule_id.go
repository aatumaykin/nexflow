package valueobject

import (
	"encoding/json"
	"fmt"
)

// ScheduleID represents a schedule identifier.
type ScheduleID ID

// String returns the string representation of the ScheduleID.
func (id ScheduleID) String() string {
	return string(id)
}

// IsEmpty returns true if the ScheduleID is empty.
func (id ScheduleID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the ScheduleID is valid (not empty and matches pattern).
func (id ScheduleID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the ScheduleID equals another ScheduleID.
func (id ScheduleID) Equals(other ScheduleID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id ScheduleID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *ScheduleID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !ScheduleID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = ScheduleID(str)
	return nil
}

// NewScheduleID creates a new ScheduleID from a string.
// Returns an error if the string is not a valid ID.
func NewScheduleID(idStr string) (ScheduleID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return ScheduleID(id), nil
}

// MustNewScheduleID creates a new ScheduleID from a string.
// Panics if the string is not a valid ID.
func MustNewScheduleID(idStr string) ScheduleID {
	id, err := NewScheduleID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
