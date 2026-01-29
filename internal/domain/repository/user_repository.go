package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create saves a new user
	Create(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by ID
	FindByID(ctx context.Context, id string) (*entity.User, error)

	// FindByChannel retrieves a user by channel and channel ID
	FindByChannel(ctx context.Context, channel, channelID string) (*entity.User, error)

	// List retrieves all users
	List(ctx context.Context) ([]*entity.User, error)

	// Delete removes a user
	Delete(ctx context.Context, id string) error
}
