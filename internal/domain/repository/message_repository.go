package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	// Create saves a new message
	Create(ctx context.Context, message *entity.Message) error

	// FindByID retrieves a message by ID
	FindByID(ctx context.Context, id string) (*entity.Message, error)

	// FindBySessionID retrieves all messages for a session
	FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Message, error)

	// Delete removes a message
	Delete(ctx context.Context, id string) error

	// DeleteBySessionID removes all messages for a session
	DeleteBySessionID(ctx context.Context, sessionID string) error
}
