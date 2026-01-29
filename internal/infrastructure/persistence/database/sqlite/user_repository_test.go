package sqlite

import (
	"context"
	"database/sql"
	"testing"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// schemaSQL contains the database schema for testing
var schemaSQL = `
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    channel TEXT NOT NULL,
    channel_user_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    UNIQUE(channel, channel_user_id)
);

CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE TABLE tasks (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    skill TEXT NOT NULL,
    input TEXT NOT NULL,
    output TEXT,
    status TEXT NOT NULL,
    error TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
);

CREATE TABLE skills (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    version TEXT NOT NULL,
    location TEXT NOT NULL,
    permissions TEXT NOT NULL,
    metadata TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE schedules (
    id TEXT PRIMARY KEY,
    skill TEXT NOT NULL,
    cron_expression TEXT NOT NULL,
    input TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (skill) REFERENCES skills(name) ON DELETE CASCADE
);

CREATE TABLE logs (
    id TEXT PRIMARY KEY,
    level TEXT NOT NULL,
    source TEXT NOT NULL,
    message TEXT NOT NULL,
    metadata TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);
`

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	require.NoError(t, err)

	// Enable foreign keys
	_, err = db.Exec("PRAGMA foreign_keys = ON")
	require.NoError(t, err)

	// Create schema
	_, err = db.Exec(schemaSQL)
	require.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")

	err := repo.Create(ctx, user)
	require.NoError(t, err)
	assert.NotEmpty(t, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	foundUser, err := repo.FindByID(ctx, user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Channel, foundUser.Channel)
	assert.Equal(t, user.ChannelID, foundUser.ChannelID)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	foundUser, err := repo.FindByID(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_FindByChannel(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	foundUser, err := repo.FindByChannel(ctx, "telegram", "user123")
	require.NoError(t, err)
	assert.Equal(t, user.ID, foundUser.ID)
	assert.Equal(t, user.Channel, foundUser.Channel)
	assert.Equal(t, user.ChannelID, foundUser.ChannelID)
}

func TestUserRepository_FindByChannel_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	foundUser, err := repo.FindByChannel(ctx, "telegram", "non-existent-user")
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_List(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Create multiple users
	user1 := entity.NewUser("telegram", "user1")
	user2 := entity.NewUser("discord", "user2")
	user3 := entity.NewUser("web", "user3")

	require.NoError(t, repo.Create(ctx, user1))
	require.NoError(t, repo.Create(ctx, user2))
	require.NoError(t, repo.Create(ctx, user3))

	users, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 3)
}

func TestUserRepository_Delete(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")

	err := repo.Create(ctx, user)
	require.NoError(t, err)

	err = repo.Delete(ctx, user.ID)
	require.NoError(t, err)

	// Verify user is deleted
	foundUser, err := repo.FindByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, foundUser)
}

func TestUserRepository_Delete_NotFound(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	err := repo.Delete(ctx, "non-existent-id")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")
}

func TestUserRepository_UniqueConstraint(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	user1 := entity.NewUser("telegram", "user123")
	user2 := entity.NewUser("telegram", "user123") // Same channel and channel ID

	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	err = repo.Create(ctx, user2)
	assert.Error(t, err)
}

func TestUserRepository_DifferentChannels(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	repo := NewUserRepository(db)

	// Same channel ID but different channels
	user1 := entity.NewUser("telegram", "user123")
	user2 := entity.NewUser("discord", "user123")

	err := repo.Create(ctx, user1)
	require.NoError(t, err)

	err = repo.Create(ctx, user2)
	require.NoError(t, err)

	users, err := repo.List(ctx)
	require.NoError(t, err)
	assert.Len(t, users, 2)
}
