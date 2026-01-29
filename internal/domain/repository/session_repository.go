package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// SessionRepository defines the interface for session data operations
type SessionRepository interface {
	// Create saves a new session
	Create(ctx context.Context, session *entity.Session) error

	// FindByID retrieves a session by ID
	FindByID(ctx context.Context, id string) (*entity.Session, error)

	// FindByUserID retrieves all sessions for a user
	FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error)

	// Update updates an existing session
	Update(ctx context.Context, session *entity.Session) error

	// Delete removes a session
	Delete(ctx context.Context, id string) error
}
