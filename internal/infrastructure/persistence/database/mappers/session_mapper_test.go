package mappers

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSession   *dbmodel.Session
		expected    *entity.Session
		expectedNil bool
	}{
		{
			name: "Valid session",
			dbSession: &dbmodel.Session{
				ID:        "session-id",
				UserID:    "user-id",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Session{
				ID:        valueobject.SessionID("session-id"),
				UserID:    valueobject.MustNewUserID("user-id"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid session with different timestamps",
			dbSession: &dbmodel.Session{
				ID:        "session-id-2",
				UserID:    "user-id-2",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Session{
				ID:        valueobject.SessionID("session-id-2"),
				UserID:    valueobject.MustNewUserID("user-id-2"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbSession:   nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SessionToDomain(tt.dbSession)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
			assert.WithinDuration(t, tt.expected.UpdatedAt, result.UpdatedAt, time.Second)
		})
	}
}

func TestSessionToDB(t *testing.T) {
	tests := []struct {
		name        string
		session     *entity.Session
		expected    *dbmodel.Session
		expectedNil bool
	}{
		{
			name: "Valid session",
			session: &entity.Session{
				ID:        valueobject.SessionID("session-id"),
				UserID:    valueobject.MustNewUserID("user-id"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &dbmodel.Session{
				ID:        "session-id",
				UserID:    "user-id",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid session with different timestamps",
			session: &entity.Session{
				ID:        valueobject.SessionID("session-id-2"),
				UserID:    valueobject.MustNewUserID("user-id-2"),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expected: &dbmodel.Session{
				ID:        "session-id-2",
				UserID:    "user-id-2",
				CreatedAt: time.Now().Format(time.RFC3339),
				UpdatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			session:     nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SessionToDB(tt.session)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.UserID, result.UserID)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
			assert.Equal(t, tt.expected.UpdatedAt, result.UpdatedAt)
		})
	}
}

func TestSessionsToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSessions  []dbmodel.Session
		expectedLen int
	}{
		{
			name: "Multiple sessions",
			dbSessions: []dbmodel.Session{
				{
					ID:        "session-1",
					UserID:    "user-1",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "session-2",
					UserID:    "user-2",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "session-3",
					UserID:    "user-3",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbSessions:  []dbmodel.Session{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbSessions:  nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SessionsToDomain(tt.dbSessions)

			if tt.dbSessions == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, session := range result {
				assert.Equal(t, tt.dbSessions[i].ID, string(session.ID))
				assert.Equal(t, tt.dbSessions[i].UserID, string(session.UserID))
			}
		})
	}
}
