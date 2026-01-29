package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	// ErrInvalidTaskStatus is returned when an invalid task status is provided.
	ErrInvalidTaskStatus = errors.New("invalid task status")
)

// TaskStatus represents the status of a task.
// It's a value object that ensures type safety for task statuses.
type TaskStatus string

const (
	// TaskStatusPending represents a task waiting to be executed.
	TaskStatusPending TaskStatus = "pending"
	// TaskStatusRunning represents a task currently running.
	TaskStatusRunning TaskStatus = "running"
	// TaskStatusCompleted represents a task completed successfully.
	TaskStatusCompleted TaskStatus = "completed"
	// TaskStatusFailed represents a task failed with an error.
	TaskStatusFailed TaskStatus = "failed"
)

// String returns the string representation of the task status.
func (s TaskStatus) String() string {
	return string(s)
}

// IsValid checks if the task status is valid.
func (s TaskStatus) IsValid() bool {
	switch s {
	case TaskStatusPending, TaskStatusRunning, TaskStatusCompleted, TaskStatusFailed:
		return true
	default:
		return false
	}
}

// IsPending returns true if the task status is pending.
func (s TaskStatus) IsPending() bool {
	return s == TaskStatusPending
}

// IsRunning returns true if the task status is running.
func (s TaskStatus) IsRunning() bool {
	return s == TaskStatusRunning
}

// IsCompleted returns true if the task status is completed.
func (s TaskStatus) IsCompleted() bool {
	return s == TaskStatusCompleted
}

// IsFailed returns true if the task status is failed.
func (s TaskStatus) IsFailed() bool {
	return s == TaskStatusFailed
}

// IsTerminal returns true if the task status is a terminal state (completed or failed).
func (s TaskStatus) IsTerminal() bool {
	return s == TaskStatusCompleted || s == TaskStatusFailed
}

// MarshalJSON implements json.Marshaler interface.
func (s TaskStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *TaskStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = TaskStatus(str)
	if !s.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidTaskStatus, str)
	}
	return nil
}

// NewTaskStatus creates a new TaskStatus from a string.
// Returns an error if the string is not a valid task status.
func NewTaskStatus(status string) (TaskStatus, error) {
	s := TaskStatus(status)
	if !s.IsValid() {
		return "", ErrInvalidTaskStatus
	}
	return s, nil
}

// MustNewTaskStatus creates a new TaskStatus from a string.
// Panics if the string is not a valid task status.
func MustNewTaskStatus(status string) TaskStatus {
	s, err := NewTaskStatus(status)
	if err != nil {
		panic(err)
	}
	return s
}
