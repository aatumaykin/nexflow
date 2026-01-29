package valueobject

import (
	"encoding/json"
	"fmt"
)

// TaskID represents a task identifier.
type TaskID ID

// String returns the string representation of the TaskID.
func (id TaskID) String() string {
	return string(id)
}

// IsEmpty returns true if the TaskID is empty.
func (id TaskID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the TaskID is valid (not empty and matches pattern).
func (id TaskID) IsValid() bool {
	return ID(id).IsValid()
}

// Equals checks if the TaskID equals another TaskID.
func (id TaskID) Equals(other TaskID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id TaskID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *TaskID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !TaskID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = TaskID(str)
	return nil
}

// NewTaskID creates a new TaskID from a string.
// Returns an error if the string is not a valid ID.
func NewTaskID(idStr string) (TaskID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return TaskID(id), nil
}

// MustNewTaskID creates a new TaskID from a string.
// Panics if the string is not a valid ID.
func MustNewTaskID(idStr string) TaskID {
	id, err := NewTaskID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}
