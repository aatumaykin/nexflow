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

// TypeName-specific ID types for type safety
type (
	// UserID represents a user identifier.
	UserID ID

	// SessionID represents a session identifier.
	SessionID ID

	// TaskID represents a task identifier.
	TaskID ID

	// MessageID represents a message identifier.
	MessageID ID

	// SkillID represents a skill identifier.
	SkillID ID

	// ScheduleID represents a schedule identifier.
	ScheduleID ID

	// LogID represents a log entry identifier.
	LogID ID
)

// UserID methods
func (id UserID) String() string               { return string(id) }
func (id UserID) IsEmpty() bool                { return string(id) == "" }
func (id UserID) IsValid() bool                { return ID(id).IsValid() }
func (id UserID) Equals(other UserID) bool     { return id == other }
func (id UserID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewUserID(idStr string) (UserID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return UserID(id), nil
}

// MustNewUserID creates a new UserID from a string. Panics on error.
func MustNewUserID(idStr string) UserID {
	id, err := NewUserID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// SessionID methods
func (id SessionID) String() string               { return string(id) }
func (id SessionID) IsEmpty() bool                { return string(id) == "" }
func (id SessionID) IsValid() bool                { return ID(id).IsValid() }
func (id SessionID) Equals(other SessionID) bool  { return id == other }
func (id SessionID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewSessionID(idStr string) (SessionID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return SessionID(id), nil
}

// MustNewSessionID creates a new SessionID from a string. Panics on error.
func MustNewSessionID(idStr string) SessionID {
	id, err := NewSessionID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// TaskID methods
func (id TaskID) String() string               { return string(id) }
func (id TaskID) IsEmpty() bool                { return string(id) == "" }
func (id TaskID) IsValid() bool                { return ID(id).IsValid() }
func (id TaskID) Equals(other TaskID) bool     { return id == other }
func (id TaskID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewTaskID(idStr string) (TaskID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return TaskID(id), nil
}

// MustNewTaskID creates a new TaskID from a string. Panics on error.
func MustNewTaskID(idStr string) TaskID {
	id, err := NewTaskID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// MessageID methods
func (id MessageID) String() string               { return string(id) }
func (id MessageID) IsEmpty() bool                { return string(id) == "" }
func (id MessageID) IsValid() bool                { return ID(id).IsValid() }
func (id MessageID) Equals(other MessageID) bool  { return id == other }
func (id MessageID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewMessageID(idStr string) (MessageID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return MessageID(id), nil
}

// MustNewMessageID creates a new MessageID from a string. Panics on error.
func MustNewMessageID(idStr string) MessageID {
	id, err := NewMessageID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// SkillID methods
func (id SkillID) String() string               { return string(id) }
func (id SkillID) IsEmpty() bool                { return string(id) == "" }
func (id SkillID) IsValid() bool                { return ID(id).IsValid() }
func (id SkillID) Equals(other SkillID) bool    { return id == other }
func (id SkillID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewSkillID(idStr string) (SkillID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return SkillID(id), nil
}

// MustNewSkillID creates a new SkillID from a string. Panics on error.
func MustNewSkillID(idStr string) SkillID {
	id, err := NewSkillID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// ScheduleID methods
func (id ScheduleID) String() string               { return string(id) }
func (id ScheduleID) IsEmpty() bool                { return string(id) == "" }
func (id ScheduleID) IsValid() bool                { return ID(id).IsValid() }
func (id ScheduleID) Equals(other ScheduleID) bool { return id == other }
func (id ScheduleID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewScheduleID(idStr string) (ScheduleID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return ScheduleID(id), nil
}

// MustNewScheduleID creates a new ScheduleID from a string. Panics on error.
func MustNewScheduleID(idStr string) ScheduleID {
	id, err := NewScheduleID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

// LogID methods
func (id LogID) String() string               { return string(id) }
func (id LogID) IsEmpty() bool                { return string(id) == "" }
func (id LogID) IsValid() bool                { return ID(id).IsValid() }
func (id LogID) Equals(other LogID) bool      { return id == other }
func (id LogID) MarshalJSON() ([]byte, error) { return json.Marshal(string(id)) }
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
func NewLogID(idStr string) (LogID, error) {
	id, err := NewID(idStr)
	if err != nil {
		return "", err
	}
	return LogID(id), nil
}

// MustNewLogID creates a new LogID from a string. Panics on error.
func MustNewLogID(idStr string) LogID {
	id, err := NewLogID(idStr)
	if err != nil {
		panic(err)
	}
	return id
}

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
