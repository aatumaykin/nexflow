package sqlite

import (
	"context"
	"testing"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionRepository_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	// Create user first
	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))

	err := sessionRepo.Create(ctx, session)
	require.NoError(t, err)
	assert.NotEmpty(t, session.ID)
}

func TestSessionRepository_FindByID(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	foundSession, err := sessionRepo.FindByID(ctx, string(session.ID))
	require.NoError(t, err)
	assert.Equal(t, session.ID, foundSession.ID)
	assert.Equal(t, session.UserID, foundSession.UserID)
}

func TestSessionRepository_FindByUserID(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session1 := entity.NewSession(string(user.ID))
	session2 := entity.NewSession(string(user.ID))
	session3 := entity.NewSession(string(user.ID))

	require.NoError(t, sessionRepo.Create(ctx, session1))
	require.NoError(t, sessionRepo.Create(ctx, session2))
	require.NoError(t, sessionRepo.Create(ctx, session3))

	sessions, err := sessionRepo.FindByUserID(ctx, string(user.ID))
	require.NoError(t, err)
	assert.Len(t, sessions, 3)
}

func TestSessionRepository_Update(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	// Update timestamp
	session.UpdateTimestamp()
	err := sessionRepo.Update(ctx, session)
	require.NoError(t, err)

	// Verify update - just check that session can be retrieved after update
	foundSession, err := sessionRepo.FindByID(ctx, string(session.ID))
	require.NoError(t, err)
	assert.Equal(t, session.ID, foundSession.ID)
	assert.Equal(t, session.UserID, foundSession.UserID)
}

func TestSessionRepository_Delete(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	err := sessionRepo.Delete(ctx, string(session.ID))
	require.NoError(t, err)

	// Verify deletion
	foundSession, err := sessionRepo.FindByID(ctx, string(session.ID))
	assert.Error(t, err)
	assert.Nil(t, foundSession)
}

func TestMessageRepository_Create(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	// Create user and session
	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	messageRepo := NewMessageRepository(db)
	message := entity.NewUserMessage(string(session.ID), "Hello, world!")

	err := messageRepo.Create(ctx, message)
	require.NoError(t, err)
	assert.NotEmpty(t, message.ID)
}

func TestMessageRepository_FindBySessionID(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	// Create user and session
	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	messageRepo := NewMessageRepository(db)
	msg1 := entity.NewUserMessage(string(session.ID), "Hello")
	msg2 := entity.NewAssistantMessage(string(session.ID), "Hi there!")

	require.NoError(t, messageRepo.Create(ctx, msg1))
	require.NoError(t, messageRepo.Create(ctx, msg2))

	messages, err := messageRepo.FindBySessionID(ctx, string(session.ID))
	require.NoError(t, err)
	assert.Len(t, messages, 2)
}

func TestMessageRepository_Roles(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	// Create user and session
	userRepo := NewUserRepository(db)
	user := entity.NewUser("telegram", "user123")
	require.NoError(t, userRepo.Create(ctx, user))

	sessionRepo := NewSessionRepository(db)
	session := entity.NewSession(string(user.ID))
	require.NoError(t, sessionRepo.Create(ctx, session))

	messageRepo := NewMessageRepository(db)
	userMsg := entity.NewUserMessage(string(session.ID), "Hello")
	assistantMsg := entity.NewAssistantMessage(string(session.ID), "Hi!")

	require.NoError(t, messageRepo.Create(ctx, userMsg))
	require.NoError(t, messageRepo.Create(ctx, assistantMsg))

	messages, err := messageRepo.FindBySessionID(ctx, string(session.ID))
	require.NoError(t, err)

	assert.Len(t, messages, 2)
	assert.True(t, messages[0].IsFromUser())
	assert.True(t, messages[1].IsFromAssistant())
}
