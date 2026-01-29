package entity

import (
	"errors"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTask(t *testing.T) {
	// Arrange & Act
	task := NewTask("session-1", "my-skill", `{"param": "value"}`)

	// Assert
	require.NotEmpty(t, task.ID)
	assert.Equal(t, "session-1", string(task.SessionID))
	assert.Equal(t, "my-skill", task.Skill)
	assert.Equal(t, `{"param": "value"}`, task.Input)
	assert.Equal(t, valueobject.TaskStatusPending, task.Status)
	assert.Empty(t, task.Output)
	assert.Empty(t, task.Error)
	assert.WithinDuration(t, time.Now(), task.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), task.UpdatedAt, time.Second)
}

func TestTask_SetRunning(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	time.Sleep(10 * time.Millisecond) // Ensure time difference

	// Act
	task.SetRunning()

	// Assert
	assert.Equal(t, valueobject.TaskStatusRunning, task.Status)
	assert.True(t, task.UpdatedAt.After(task.CreatedAt))
}

func TestTask_SetCompleted(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	output := `{"result": "success"}`

	// Act
	task.SetCompleted(output)

	// Assert
	assert.Equal(t, valueobject.TaskStatusCompleted, task.Status)
	assert.Equal(t, output, task.Output)
	assert.Empty(t, task.Error)
	assert.True(t, task.UpdatedAt.After(task.CreatedAt))
}

func TestTask_SetFailed(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	time.Sleep(10 * time.Millisecond) // Ensure time difference
	err := errors.New("execution failed")

	// Act
	task.SetFailed(err.Error())

	// Assert
	assert.Equal(t, valueobject.TaskStatusFailed, task.Status)
	assert.Empty(t, task.Output)
	assert.Equal(t, "execution failed", task.Error)
	assert.True(t, task.UpdatedAt.After(task.CreatedAt))
}

func TestTask_SetFailed_NilError(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")

	// Act
	task.SetFailed("")

	// Assert
	assert.Equal(t, valueobject.TaskStatusFailed, task.Status)
	assert.Empty(t, task.Error)
}

func TestTask_IsPending(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")

	// Act & Assert
	assert.True(t, task.IsPending())
	assert.False(t, task.IsRunning())
	assert.False(t, task.IsCompleted())
	assert.False(t, task.IsFailed())
}

func TestTask_IsRunning(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	task.SetRunning()

	// Act & Assert
	assert.False(t, task.IsPending())
	assert.True(t, task.IsRunning())
	assert.False(t, task.IsCompleted())
	assert.False(t, task.IsFailed())
}

func TestTask_IsCompleted(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	task.SetCompleted(`{"result": "ok"}`)

	// Act & Assert
	assert.False(t, task.IsPending())
	assert.False(t, task.IsRunning())
	assert.True(t, task.IsCompleted())
	assert.False(t, task.IsFailed())
}

func TestTask_IsFailed(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	task.SetFailed("error")

	// Act & Assert
	assert.False(t, task.IsPending())
	assert.False(t, task.IsRunning())
	assert.False(t, task.IsCompleted())
	assert.True(t, task.IsFailed())
}

func TestTask_BelongsToSession(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")

	// Act & Assert
	assert.True(t, task.BelongsToSession(valueobject.SessionID("session-1")))
	assert.False(t, task.BelongsToSession(valueobject.SessionID("session-2")))
}

func TestTask_StatusTransitions(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")

	// Act & Assert - Pending to Running
	assert.True(t, task.IsPending())
	task.SetRunning()
	assert.True(t, task.IsRunning())

	// Running to Completed
	task.SetCompleted(`{}`)
	assert.True(t, task.IsCompleted())
}

func TestTask_UpdatedTimestampChanges(t *testing.T) {
	// Arrange
	task := NewTask("session-1", "skill", "{}")
	initialUpdatedAt := task.UpdatedAt
	time.Sleep(10 * time.Millisecond)

	// Act - trigger update
	task.SetRunning()

	// Assert
	assert.True(t, task.UpdatedAt.After(initialUpdatedAt))
}
