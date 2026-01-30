package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.MessageRepository = (*MessageRepository)(nil)

type MessageRepository struct {
	queries *database.Queries
}

func NewMessageRepository(queries *database.Queries) *MessageRepository {
	return &MessageRepository{queries: queries}
}

func (r *MessageRepository) Create(ctx context.Context, message *entity.Message) error {
	dbMessage := mappers.MessageToDB(message)
	if dbMessage == nil {
		return fmt.Errorf("failed to convert message to db model")
	}

	_, err := r.queries.CreateMessage(ctx, database.CreateMessageParams{
		ID:        dbMessage.ID,
		SessionID: dbMessage.SessionID,
		Role:      dbMessage.Role,
		Content:   dbMessage.Content,
		CreatedAt: dbMessage.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

func (r *MessageRepository) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	dbMessage, err := r.queries.GetMessageByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("message not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find message by id: %w", err)
	}

	return mappers.MessageToDomain(&dbMessage), nil
}

func (r *MessageRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Message, error) {
	dbMessages, err := r.queries.GetMessagesBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find messages by session id: %w", err)
	}

	return mappers.MessagesToDomain(dbMessages), nil
}

func (r *MessageRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetMessageByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("message not found: %s", id)
		}
		return fmt.Errorf("failed to check message existence: %w", err)
	}

	err = r.queries.DeleteMessage(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	return nil
}

func (r *MessageRepository) DeleteBySessionID(ctx context.Context, sessionID string) error {
	messages, err := r.queries.GetMessagesBySessionID(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to find messages by session id: %w", err)
	}

	for _, msg := range messages {
		if err := r.queries.DeleteMessage(ctx, msg.ID); err != nil {
			return fmt.Errorf("failed to delete message: %w", err)
		}
	}

	return nil
}
