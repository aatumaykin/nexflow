package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.SessionRepository = (*SessionRepository)(nil)

// SessionRepository implements repository.SessionRepository using SQLC-generated queries
type SessionRepository struct {
	db *sql.DB
}

// NewSessionRepository creates a new SessionRepository instance
func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Create saves a new session
func (r *SessionRepository) Create(ctx context.Context, session *entity.Session) error {
	dbSession := mappers.SessionToDB(session)
	if dbSession == nil {
		return fmt.Errorf("failed to convert session to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sessions (id, user_id, created_at, updated_at) VALUES (?, ?, ?, ?)`,
		dbSession.ID, dbSession.UserID, dbSession.CreatedAt, dbSession.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

// FindByID retrieves a session by ID
func (r *SessionRepository) FindByID(ctx context.Context, id string) (*entity.Session, error) {
	var dbSession dbmodel.Session

	err := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, created_at, updated_at FROM sessions WHERE id = ? LIMIT 1`,
		id,
	).Scan(&dbSession.ID, &dbSession.UserID, &dbSession.CreatedAt, &dbSession.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find session by id: %w", err)
	}

	return mappers.SessionToDomain(&dbSession), nil
}

// FindByUserID retrieves all sessions for a user
func (r *SessionRepository) FindByUserID(ctx context.Context, userID string) ([]*entity.Session, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, user_id, created_at, updated_at FROM sessions WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find sessions by user id: %w", err)
	}
	defer rows.Close()

	var dbSessions []dbmodel.Session
	for rows.Next() {
		var dbSession dbmodel.Session

		if err := rows.Scan(&dbSession.ID, &dbSession.UserID, &dbSession.CreatedAt, &dbSession.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan session: %w", err)
		}

		dbSessions = append(dbSessions, dbSession)
	}

	return mappers.SessionsToDomain(dbSessions), nil
}

// Update updates an existing session
func (r *SessionRepository) Update(ctx context.Context, session *entity.Session) error {
	// Update timestamp before saving
	session.UpdateTimestamp()

	dbSession := mappers.SessionToDB(session)
	if dbSession == nil {
		return fmt.Errorf("failed to convert session to db model")
	}

	now := time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET updated_at = ? WHERE id = ?`,
		now, dbSession.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	return nil
}

// Delete removes a session
func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found: %s", id)
	}

	return nil
}
