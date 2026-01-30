package mappers

import (
	"database/sql"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbTask      *dbmodel.Task
		expected    *entity.Task
		expectedNil bool
	}{
		{
			name: "Valid task with output",
			dbTask: &dbmodel.Task{
				ID:        "task-id",
				SessionID: "session-id",
				Skill:     "test-skill",
				Input:     "test input",
				Output:    sql.NullString{String: "test output", Valid: true},
				Status:    "completed",
				Error:     sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Task{
				ID:        valueobject.TaskID("task-id"),
				SessionID: valueobject.MustNewSessionID("session-id"),
				Skill:     "test-skill",
				Input:     "test input",
				Output:    "test output",
				Status:    valueobject.MustNewTaskStatus("completed"),
				Error:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid task without output",
			dbTask: &dbmodel.Task{
				ID:        "task-id-2",
				SessionID: "session-id-2",
				Skill:     "test-skill-2",
				Input:     "test input 2",
				Output:    sql.NullString{Valid: false},
				Status:    "pending",
				Error:     sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Task{
				ID:        valueobject.TaskID("task-id-2"),
				SessionID: valueobject.MustNewSessionID("session-id-2"),
				Skill:     "test-skill-2",
				Input:     "test input 2",
				Output:    "",
				Status:    valueobject.MustNewTaskStatus("pending"),
				Error:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid task with error",
			dbTask: &dbmodel.Task{
				ID:        "task-id-3",
				SessionID: "session-id-3",
				Skill:     "test-skill-3",
				Input:     "test input 3",
				Output:    sql.NullString{Valid: false},
				Status:    "failed",
				Error:     sql.NullString{String: "test error", Valid: true},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Task{
				ID:        valueobject.TaskID("task-id-3"),
				SessionID: valueobject.MustNewSessionID("session-id-3"),
				Skill:     "test-skill-3",
				Input:     "test input 3",
				Output:    "",
				Status:    valueobject.MustNewTaskStatus("failed"),
				Error:     "test error",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbTask:      nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TaskToDomain(tt.dbTask)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SessionID, result.SessionID)
			assert.Equal(t, tt.expected.Skill, result.Skill)
			assert.Equal(t, tt.expected.Input, result.Input)
			assert.Equal(t, tt.expected.Output, result.Output)
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Error, result.Error)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
			assert.WithinDuration(t, tt.expected.UpdatedAt, result.UpdatedAt, time.Second)
		})
	}
}

func TestTaskToDB(t *testing.T) {
	tests := []struct {
		name        string
		task        *entity.Task
		expected    *dbmodel.Task
		expectedNil bool
	}{
		{
			name: "Valid task with output",
			task: &entity.Task{
				ID:        valueobject.TaskID("task-id"),
				SessionID: valueobject.MustNewSessionID("session-id"),
				Skill:     "test-skill",
				Input:     "test input",
				Output:    "test output",
				Status:    valueobject.MustNewTaskStatus("completed"),
				Error:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &dbmodel.Task{
				ID:        "task-id",
				SessionID: "session-id",
				Skill:     "test-skill",
				Input:     "test input",
				Output:    sql.NullString{String: "test output", Valid: true},
				Status:    "completed",
				Error:     sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid task without output",
			task: &entity.Task{
				ID:        valueobject.TaskID("task-id-2"),
				SessionID: valueobject.MustNewSessionID("session-id-2"),
				Skill:     "test-skill-2",
				Input:     "test input 2",
				Output:    "",
				Status:    valueobject.MustNewTaskStatus("pending"),
				Error:     "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &dbmodel.Task{
				ID:        "task-id-2",
				SessionID: "session-id-2",
				Skill:     "test-skill-2",
				Input:     "test input 2",
				Output:    sql.NullString{Valid: false},
				Status:    "pending",
				Error:     sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid task with error",
			task: &entity.Task{
				ID:        valueobject.TaskID("task-id-3"),
				SessionID: valueobject.MustNewSessionID("session-id-3"),
				Skill:     "test-skill-3",
				Input:     "test input 3",
				Output:    "",
				Status:    valueobject.MustNewTaskStatus("failed"),
				Error:     "test error",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &dbmodel.Task{
				ID:        "task-id-3",
				SessionID: "session-id-3",
				Skill:     "test-skill-3",
				Input:     "test input 3",
				Output:    sql.NullString{Valid: false},
				Status:    "failed",
				Error:     sql.NullString{String: "test error", Valid: true},
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			task:        nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TaskToDB(tt.task)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SessionID, result.SessionID)
			assert.Equal(t, tt.expected.Skill, result.Skill)
			assert.Equal(t, tt.expected.Input, result.Input)
			assert.Equal(t, tt.expected.Output.Valid, tt.expected.Output.Valid)
			if tt.expected.Output.Valid {
				assert.Equal(t, tt.expected.Output.String, result.Output.String)
			}
			assert.Equal(t, tt.expected.Status, result.Status)
			assert.Equal(t, tt.expected.Error.Valid, tt.expected.Error.Valid)
			if tt.expected.Error.Valid {
				assert.Equal(t, tt.expected.Error.String, result.Error.String)
			}
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestTasksToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbTasks     []dbmodel.Task
		expectedLen int
	}{
		{
			name: "Multiple tasks",
			dbTasks: []dbmodel.Task{
				{
					ID:        "task-1",
					SessionID: "session-1",
					Skill:     "skill-1",
					Input:     "input-1",
					Output:    sql.NullString{String: "output-1", Valid: true},
					Status:    "completed",
					Error:     sql.NullString{Valid: false},
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "task-2",
					SessionID: "session-2",
					Skill:     "skill-2",
					Input:     "input-2",
					Output:    sql.NullString{Valid: false},
					Status:    "pending",
					Error:     sql.NullString{Valid: false},
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "task-3",
					SessionID: "session-3",
					Skill:     "skill-3",
					Input:     "input-3",
					Output:    sql.NullString{Valid: false},
					Status:    "failed",
					Error:     sql.NullString{String: "error-3", Valid: true},
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbTasks:     []dbmodel.Task{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbTasks:     nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TasksToDomain(tt.dbTasks)

			if tt.dbTasks == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, task := range result {
				assert.Equal(t, tt.dbTasks[i].ID, string(task.ID))
				assert.Equal(t, tt.dbTasks[i].SessionID, string(task.SessionID))
				assert.Equal(t, tt.dbTasks[i].Skill, task.Skill)
				assert.Equal(t, tt.dbTasks[i].Status, string(task.Status))
			}
		})
	}
}
