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

var _ repository.UserRepository = (*UserRepository)(nil)

// UserRepository implements repository.UserRepository using SQLC-generated queries
type UserRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create saves a new user
func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
	dbUser := mappers.UserToDB(user)
	if dbUser == nil {
		return fmt.Errorf("failed to convert user to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, channel, channel_user_id, created_at) VALUES (?, ?, ?, ?)`,
		dbUser.ID, dbUser.Channel, dbUser.ChannelUserID, dbUser.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// FindByID retrieves a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
	var sqlcUser dbmodel.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, channel, channel_user_id, created_at FROM users WHERE id = ? LIMIT 1`,
		id,
	).Scan(&sqlcUser.ID, &sqlcUser.Channel, &sqlcUser.ChannelUserID, &sqlcUser.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return mappers.UserToDomain(&sqlcUser), nil
}

// FindByChannel retrieves a user by channel and channel ID
func (r *UserRepository) FindByChannel(ctx context.Context, channel, channelID string) (*entity.User, error) {
	var sqlcUser dbmodel.User

	err := r.db.QueryRowContext(ctx,
		`SELECT id, channel, channel_user_id, created_at FROM users WHERE channel = ? AND channel_user_id = ? LIMIT 1`,
		channel, channelID,
	).Scan(&sqlcUser.ID, &sqlcUser.Channel, &sqlcUser.ChannelUserID, &sqlcUser.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found: channel=%s, channelID=%s", channel, channelID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by channel: %w", err)
	}

	return mappers.UserToDomain(&sqlcUser), nil
}

// List retrieves all users
func (r *UserRepository) List(ctx context.Context) ([]*entity.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, channel, channel_user_id, created_at FROM users ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var dbUsers []dbmodel.User
	for rows.Next() {
		var dbUser dbmodel.User

		if err := rows.Scan(&dbUser.ID, &dbUser.Channel, &dbUser.ChannelUserID, &dbUser.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}

		dbUsers = append(dbUsers, dbUser)
	}

	return mappers.UsersToDomain(dbUsers), nil
}

// Delete removes a user
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found: %s", id)
	}

	return nil
}
