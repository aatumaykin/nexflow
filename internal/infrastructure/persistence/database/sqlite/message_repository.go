package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.MessageRepository = (*MessageRepository)(nil)

// MessageRepository implements repository.MessageRepository using SQLC-generated queries
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new MessageRepository instance
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create saves a new message
func (r *MessageRepository) Create(ctx context.Context, message *entity.Message) error {
	dbMessage := mappers.MessageToDB(message)
	if dbMessage == nil {
		return fmt.Errorf("failed to convert message to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO messages (id, session_id, role, content, created_at) VALUES (?, ?, ?, ?, ?)`,
		dbMessage.ID, dbMessage.SessionID, dbMessage.Role, dbMessage.Content, dbMessage.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// FindByID retrieves a message by ID
func (r *MessageRepository) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	var dbMessage dbmodel.Message

	err := r.db.QueryRowContext(ctx,
		`SELECT id, session_id, role, content, created_at FROM messages WHERE id = ? LIMIT 1`,
		id,
	).Scan(&dbMessage.ID, &dbMessage.SessionID, &dbMessage.Role, &dbMessage.Content, &dbMessage.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find message by id: %w", err)
	}

	return mappers.MessageToDomain(&dbMessage), nil
}

// FindBySessionID retrieves all messages for a session
func (r *MessageRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Message, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, role, content, created_at FROM messages WHERE session_id = ? ORDER BY created_at ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages by session id: %w", err)
	}
	defer rows.Close()

	var dbMessages []dbmodel.Message
	for rows.Next() {
		var dbMessage dbmodel.Message

		if err := rows.Scan(&dbMessage.ID, &dbMessage.SessionID, &dbMessage.Role, &dbMessage.Content, &dbMessage.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}

		dbMessages = append(dbMessages, dbMessage)
	}

	return mappers.MessagesToDomain(dbMessages), nil
}

// Delete removes a message
func (r *MessageRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM messages WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("message not found: %s", id)
	}

	return nil
}

// DeleteBySessionID removes all messages for a session
func (r *MessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM messages WHERE session_id = ?`, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete messages by session id: %w", err)
	}

	return nil
}
