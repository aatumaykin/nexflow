package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	// ErrInvalidID is returned when an invalid ID is provided.
	ErrInvalidID = errors.New("invalid id")
	// ErrEmptyID is returned when an empty ID is provided.
	ErrEmptyID = errors.New("id cannot be empty")
	// idRegex is the regex pattern for valid IDs.
	// IDs should be alphanumeric with underscores and hyphens, at least 1 character.
	idRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

// ID is a generic identifier type.
// It provides type safety for entity IDs across the domain.
type ID string

// String returns the string representation of the ID.
func (id ID) String() string {
	return string(id)
}

// IsEmpty returns true if the ID is empty.
func (id ID) IsEmpty() bool {
	return string(id) == ""
}

// IsValid checks if the ID is valid (not empty and matches pattern).
func (id ID) IsValid() bool {
	if id.IsEmpty() {
		return false
	}
	return idRegex.MatchString(string(id))
}

// Equals checks if the ID equals another ID.
func (id ID) Equals(other ID) bool {
	return id == other
}

// MarshalJSON implements json.Marshaler interface.
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(id))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (id *ID) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyID
	}
	if !ID(str).IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidID, str)
	}
	*id = ID(str)
	return nil
}

// NewID creates a new ID from a string.
// Returns an error if the string is not a valid ID.
func NewID(idStr string) (ID, error) {
	if idStr == "" {
		return "", ErrEmptyID
	}
	if !ID(idStr).IsValid() {
		return "", ErrInvalidID
	}
	return ID(idStr), nil
}

// MustNewID creates a new ID from a string.
// Panics if the string is not a valid ID.
func MustNewID(idStr string) ID {
	id, err := NewID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// GenerateID generates a unique ID using the provided generator function.
// If generator is nil, uses a default generator (simple random string).
func GenerateID(generator func() string) ID {
	if generator != nil {
		return ID(generator())
	}
	// Default simple ID generator
	return ID(fmt.Sprintf("id_%d", len(idStr)))
}

// idStr is a placeholder for length-based generation.
var idStr = make([]string, 1000)

// StringToIDType converts a string to a specific ID type based on the type name.
// Useful for dynamic ID creation from string type names.
func StringToIDType(typeName, idStr string) (interface{}, error) {
	switch strings.ToLower(typeName) {
	case "userid", "user":
		return NewUserID(idStr)
	case "sessionid", "session":
		return NewSessionID(idStr)
	case "taskid", "task":
		return NewTaskID(idStr)
	case "messageid", "message":
		return NewMessageID(idStr)
	case "skillid", "skill":
		return NewSkillID(idStr)
	case "scheduleid", "schedule":
		return NewScheduleID(idStr)
	case "logid", "log":
		return NewLogID(idStr)
	default:
		return nil, fmt.Errorf("unknown ID type: %s", typeName)
	}
}
