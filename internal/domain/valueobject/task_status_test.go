package valueobject

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskStatus_String(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected string
	}{
		{"pending", TaskStatusPending, "pending"},
		{"running", TaskStatusRunning, "running"},
		{"completed", TaskStatusCompleted, "completed"},
		{"failed", TaskStatusFailed, "failed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.String())
		})
	}
}

func TestTaskStatus_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, true},
		{"running", TaskStatusRunning, true},
		{"completed", TaskStatusCompleted, true},
		{"failed", TaskStatusFailed, true},
		{"invalid", TaskStatus("invalid"), false},
		{"empty", TaskStatus(""), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsValid())
		})
	}
}

func TestTaskStatus_IsPending(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, true},
		{"running", TaskStatusRunning, false},
		{"completed", TaskStatusCompleted, false},
		{"failed", TaskStatusFailed, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsPending())
		})
	}
}

func TestTaskStatus_IsRunning(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, false},
		{"running", TaskStatusRunning, true},
		{"completed", TaskStatusCompleted, false},
		{"failed", TaskStatusFailed, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsRunning())
		})
	}
}

func TestTaskStatus_IsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, false},
		{"running", TaskStatusRunning, false},
		{"completed", TaskStatusCompleted, true},
		{"failed", TaskStatusFailed, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsCompleted())
		})
	}
}

func TestTaskStatus_IsFailed(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, false},
		{"running", TaskStatusRunning, false},
		{"completed", TaskStatusCompleted, false},
		{"failed", TaskStatusFailed, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsFailed())
		})
	}
}

func TestTaskStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		name     string
		status   TaskStatus
		expected bool
	}{
		{"pending", TaskStatusPending, false},
		{"running", TaskStatusRunning, false},
		{"completed", TaskStatusCompleted, true},
		{"failed", TaskStatusFailed, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.status.IsTerminal())
		})
	}
}

func TestTaskStatus_MarshalJSON(t *testing.T) {
	status := TaskStatusPending
	data, err := json.Marshal(status)

	require.NoError(t, err)
	assert.JSONEq(t, `"pending"`, string(data))
}

func TestTaskStatus_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expected  TaskStatus
		expectErr bool
	}{
		{"pending", `"pending"`, TaskStatusPending, false},
		{"running", `"running"`, TaskStatusRunning, false},
		{"completed", `"completed"`, TaskStatusCompleted, false},
		{"failed", `"failed"`, TaskStatusFailed, false},
		{"invalid", `"invalid"`, TaskStatus(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var status TaskStatus
			err := json.Unmarshal([]byte(tt.data), &status)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, status)
			}
		})
	}
}

func TestNewTaskStatus(t *testing.T) {
	tests := []struct {
		name      string
		statusStr string
		expected  TaskStatus
		expectErr bool
	}{
		{"pending", "pending", TaskStatusPending, false},
		{"running", "running", TaskStatusRunning, false},
		{"completed", "completed", TaskStatusCompleted, false},
		{"failed", "failed", TaskStatusFailed, false},
		{"invalid", "invalid", TaskStatus(""), true},
		{"empty", "", TaskStatus(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewTaskStatus(tt.statusStr)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, status)
			}
		})
	}
}

func TestMustNewTaskStatus(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewTaskStatus("pending")
	})

	assert.Panics(t, func() {
		MustNewTaskStatus("invalid")
	})
}
