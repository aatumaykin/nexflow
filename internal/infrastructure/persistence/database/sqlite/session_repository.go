package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.SessionRepository = (*SessionRepository)(nil)

type SessionRepository struct {
	queries *database.Queries
}

func NewSessionRepository(queries *database.Queries) *SessionRepository {
	return &SessionRepository{queries: queries}
}

func (r *SessionRepository) Create(ctx context.Context, session *entity.Session) error {
	dbSession := mappers.SessionToDB(session)
	if dbSession == nil {
		return fmt.Errorf("failed to convert session to db model")
	}

	_, err := r.queries.CreateSession(ctx, database.CreateSessionParams{
		ID:        dbSession.ID,
		UserID:    dbSession.UserID,
		CreatedAt: dbSession.CreatedAt,
		UpdatedAt: dbSession.UpdatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	dbSession, err := r.queries.GetSessionByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find session by id: %w", err)
	}

	return mappers.SessionToDomain(&dbSession), nil
}

func (r *SessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	dbSessions, err := r.queries.GetSessionsByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions by user id: %w", err)
	}

	return mappers.SessionsToDomain(dbSessions), nil
}

func (r *SessionRepository) Update(ctx context.Context, session *entity.Session) error {
	session.UpdateTimestamp()

	dbSession := mappers.SessionToDB(session)
	if dbSession == nil {
		return fmt.Errorf("failed to convert session to db model")
	}

	_, err := r.queries.UpdateSession(ctx, database.UpdateSessionParams{
		UpdatedAt: time.Now().Format(time.RFC3339),
		ID:        dbSession.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetSessionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("session not found: %s", id)
		}
		return fmt.Errorf("failed to check session existence: %w", err)
	}

	err = r.queries.DeleteSession(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}
