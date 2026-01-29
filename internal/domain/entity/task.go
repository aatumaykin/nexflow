package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Task represents a skill execution task.
// Tasks track skill execution, status, and results.
type Task struct {
	ID        valueobject.TaskID     `json:"id"`         // Unique identifier for the task
	SessionID valueobject.SessionID  `json:"session_id"` // ID of the session this task belongs to
	Skill     string                 `json:"skill"`      // Name of the skill to execute
	Input     string                 `json:"input"`      // Input parameters in JSON format
	Output    string                 `json:"output"`     // Output result in JSON format
	Status    valueobject.TaskStatus `json:"status"`     // Task status: "pending", "running", "completed", "failed"
	Error     string                 `json:"error"`      // Error message if the task failed
	CreatedAt time.Time              `json:"created_at"` // Timestamp when the task was created
	UpdatedAt time.Time              `json:"updated_at"` // Timestamp when the task was last updated
}

// NewTask creates a new pending task for the specified session and skill with input parameters.
func NewTask(sessionID, skill, input string) *Task {
	now := utils.Now()
	return &Task{
		ID:        valueobject.TaskID(utils.GenerateID()),
		SessionID: valueobject.MustNewSessionID(sessionID),
		Skill:     skill,
		Input:     input,
		Status:    valueobject.TaskStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// SetRunning sets the task status to running and updates the timestamp.
func (t *Task) SetRunning() {
	t.Status = valueobject.TaskStatusRunning
	t.UpdatedAt = utils.Now()
}

// SetCompleted sets the task status to completed with the output and updates the timestamp.
func (t *Task) SetCompleted(output string) {
	t.Status = valueobject.TaskStatusCompleted
	t.Output = output
	t.UpdatedAt = utils.Now()
}

// SetFailed sets the task status to failed with an error message and updates the timestamp.
func (t *Task) SetFailed(err string) {
	t.Status = valueobject.TaskStatusFailed
	if err != "" {
		t.Error = err
	}
	t.UpdatedAt = utils.Now()
}

// IsPending returns true if the task is pending.
func (t *Task) IsPending() bool {
	return t.Status == valueobject.TaskStatusPending
}

// IsRunning returns true if the task is currently running.
func (t *Task) IsRunning() bool {
	return t.Status == valueobject.TaskStatusRunning
}

// IsCompleted returns true if the task completed successfully.
func (t *Task) IsCompleted() bool {
	return t.Status == valueobject.TaskStatusCompleted
}

// IsFailed returns true if the task failed.
func (t *Task) IsFailed() bool {
	return t.Status == valueobject.TaskStatusFailed
}

// BelongsToSession returns true if the task belongs to the specified session.
func (t *Task) BelongsToSession(sessionID valueobject.SessionID) bool {
	return t.SessionID.Equals(sessionID)
}

// GetInput parses and returns the input parameters as a map.
// Returns nil if parsing fails or input is empty.
func (t *Task) GetInput() map[string]interface{} {
	return utils.UnmarshalJSONToMap(t.Input)
}

// GetOutput parses and returns the output result as a map.
// Returns nil if parsing fails or output is empty.
func (t *Task) GetOutput() map[string]interface{} {
	return utils.UnmarshalJSONToMap(t.Output)
}
