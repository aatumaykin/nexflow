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

var _ repository.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	queries *database.Queries
}

func NewUserRepository(queries *database.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	dbUser := mappers.UserToDB(user)
	if dbUser == nil {
		return fmt.Errorf("failed to convert user to db model")
	}

	_, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		ID:            dbUser.ID,
		Channel:       dbUser.Channel,
		ChannelUserID: dbUser.ChannelUserID,
		CreatedAt:     dbUser.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return mappers.UserToDomain(&sqlcUser), nil
}

func (r *UserRepository) FindByChannel(ctx context.Context, channel, channelID string) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByChannel(ctx, database.GetUserByChannelParams{
		Channel:       channel,
		ChannelUserID: channelID,
	})

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: channel=%s, channelID=%s", channel, channelID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by channel: %w", err)
	}

	return mappers.UserToDomain(&sqlcUser), nil
}

func (r *UserRepository) List(ctx context.Context) ([]*entity.User, error) {
	dbUsers, err := r.queries.ListUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return mappers.UsersToDomain(dbUsers), nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found: %s", id)
		}
		return fmt.Errorf("failed to check user existence: %w", err)
	}

	err = r.queries.DeleteUser(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
