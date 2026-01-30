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

func TestLogToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbLog       *dbmodel.Log
		expected    *entity.Log
		expectedNil bool
	}{
		{
			name: "Valid log with metadata",
			dbLog: &dbmodel.Log{
				ID:        "test-id",
				Level:     "info",
				Source:    "test-source",
				Message:   "Test message",
				Metadata:  sql.NullString{String: "test-metadata", Valid: true},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Log{
				ID:        valueobject.LogID("test-id"),
				Level:     valueobject.MustNewLogLevel("info"),
				Source:    "test-source",
				Message:   "Test message",
				Metadata:  "test-metadata",
				CreatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid log without metadata",
			dbLog: &dbmodel.Log{
				ID:        "test-id-2",
				Level:     "error",
				Source:    "test-source-2",
				Message:   "Test message 2",
				Metadata:  sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Log{
				ID:        valueobject.LogID("test-id-2"),
				Level:     valueobject.MustNewLogLevel("error"),
				Source:    "test-source-2",
				Message:   "Test message 2",
				Metadata:  "",
				CreatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbLog:       nil,
			expected:    nil,
			expectedNil: true,
		},
		{
			name: "Log with SQL NULL metadata",
			dbLog: &dbmodel.Log{
				ID:        "test-id-3",
				Level:     "warn",
				Source:    "test-source-3",
				Message:   "Test message 3",
				Metadata:  sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Log{
				ID:        valueobject.LogID("test-id-3"),
				Level:     valueobject.MustNewLogLevel("warn"),
				Source:    "test-source-3",
				Message:   "Test message 3",
				Metadata:  "",
				CreatedAt: time.Now(),
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LogToDomain(tt.dbLog)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Level, result.Level)
			assert.Equal(t, tt.expected.Source, result.Source)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Metadata, result.Metadata)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
		})
	}
}

func TestLogToDB(t *testing.T) {
	tests := []struct {
		name        string
		log         *entity.Log
		expected    *dbmodel.Log
		expectedNil bool
	}{
		{
			name: "Valid log with metadata",
			log: &entity.Log{
				ID:        valueobject.LogID("test-id"),
				Level:     valueobject.MustNewLogLevel("info"),
				Source:    "test-source",
				Message:   "Test message",
				Metadata:  "test-metadata",
				CreatedAt: time.Now(),
			},
			expected: &dbmodel.Log{
				ID:        "test-id",
				Level:     "info",
				Source:    "test-source",
				Message:   "Test message",
				Metadata:  sql.NullString{String: "test-metadata", Valid: true},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid log without metadata",
			log: &entity.Log{
				ID:        valueobject.LogID("test-id-2"),
				Level:     valueobject.MustNewLogLevel("error"),
				Source:    "test-source-2",
				Message:   "Test message 2",
				Metadata:  "",
				CreatedAt: time.Now(),
			},
			expected: &dbmodel.Log{
				ID:        "test-id-2",
				Level:     "error",
				Source:    "test-source-2",
				Message:   "Test message 2",
				Metadata:  sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			log:         nil,
			expected:    nil,
			expectedNil: true,
		},
		{
			name: "Log with empty metadata",
			log: &entity.Log{
				ID:        valueobject.LogID("test-id-3"),
				Level:     valueobject.MustNewLogLevel("warn"),
				Source:    "test-source-3",
				Message:   "Test message 3",
				Metadata:  "",
				CreatedAt: time.Now(),
			},
			expected: &dbmodel.Log{
				ID:        "test-id-3",
				Level:     "warn",
				Source:    "test-source-3",
				Message:   "Test message 3",
				Metadata:  sql.NullString{Valid: false},
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LogToDB(tt.log)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Level, result.Level)
			assert.Equal(t, tt.expected.Source, result.Source)
			assert.Equal(t, tt.expected.Message, result.Message)
			assert.Equal(t, tt.expected.Metadata.Valid, tt.expected.Metadata.Valid)
			if tt.expected.Metadata.Valid {
				assert.Equal(t, tt.expected.Metadata.String, result.Metadata.String)
			}
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestLogsToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbLogs      []dbmodel.Log
		expectedLen int
	}{
		{
			name: "Multiple logs",
			dbLogs: []dbmodel.Log{
				{
					ID:        "log-1",
					Level:     "info",
					Source:    "source-1",
					Message:   "Message 1",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "log-2",
					Level:     "error",
					Source:    "source-2",
					Message:   "Message 2",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "log-3",
					Level:     "warn",
					Source:    "source-3",
					Message:   "Message 3",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbLogs:      []dbmodel.Log{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbLogs:      nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LogsToDomain(tt.dbLogs)

			if tt.dbLogs == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, log := range result {
				assert.Equal(t, tt.dbLogs[i].ID, string(log.ID))
				assert.Equal(t, tt.dbLogs[i].Level, string(log.Level))
				assert.Equal(t, tt.dbLogs[i].Message, log.Message)
			}
		})
	}
}
