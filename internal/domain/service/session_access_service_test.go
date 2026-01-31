package service

import (
	"context"
	"testing"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockSessionRepository is a mock implementation of SessionRepository for testing
type mockSessionRepository struct {
	findByIDFunc func(ctx context.Context, id string) (*entity.Session, error)
}

// FindByID implements SessionRepository interface
func (m *mockSessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	if m.findByIDFunc != nil {
		return m.findByIDFunc(ctx, id)
	}
	return nil, nil
}

// Other SessionRepository methods (not used in this test)
func (m *mockSessionRepository) Create(ctx context.Context, session *entity.Session) error {
	return nil
}

func (m *mockSessionRepository) Update(ctx context.Context, session *entity.Session) error {
	return nil
}

func (m *mockSessionRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (m *mockSessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	return nil, nil
}

// TestSessionAccessService_CanAccessSession tests session access control
func TestSessionAccessService_CanAccessSession(t *testing.T) {
	tests := []struct {
		name       string
		userID     string
		sessionID  string
		session    *entity.Session
		findError  error
		expectedOK bool
	}{
		{
			name:       "user can access own session",
			userID:     "user-1",
			sessionID:  "session-1",
			session:    entity.NewSession("user-1"),
			findError:  nil,
			expectedOK: true,
		},
		{
			name:       "user cannot access another user's session",
			userID:     "user-1",
			sessionID:  "session-2",
			session:    entity.NewSession("user-2"),
			findError:  nil,
			expectedOK: false,
		},
		{
			name:       "access denied when session not found",
			userID:     "user-1",
			sessionID:  "non-existent-session",
			session:    nil,
			findError:  assert.AnError,
			expectedOK: false,
		},
		{
			name:       "access denied when session is nil",
			userID:     "user-1",
			sessionID:  "session-1",
			session:    nil,
			findError:  nil,
			expectedOK: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockRepo := &mockSessionRepository{
				findByIDFunc: func(ctx context.Context, id string) (*entity.Session, error) {
					return tt.session, tt.findError
				},
			}
			service := NewSessionAccessService(mockRepo)

			// Act
			result := service.CanAccessSession(context.Background(),
				valueobject.MustNewUserID(tt.userID),
				valueobject.MustNewSessionID(tt.sessionID))

			// Assert
			assert.Equal(t, tt.expectedOK, result)
		})
	}
}

// TestNewSessionAccessService tests service creation
func TestNewSessionAccessService(t *testing.T) {
	// Arrange
	mockRepo := &mockSessionRepository{}

	// Act
	service := NewSessionAccessService(mockRepo)

	// Assert
	require.NotNil(t, service)
	assert.Equal(t, mockRepo, service.sessionRepo)
}
