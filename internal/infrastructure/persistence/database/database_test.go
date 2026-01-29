package database

import (
	"context"
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/shared/logging"
	"github.com/google/uuid"
)

var testDB *DB
var testCtx context.Context

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

func setupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	// Create in-memory SQLite database
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	// Create schema
	if _, err := db.Exec(schemaSQL); err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db
}

func setupTestDBInstance(t *testing.T) *DB {
	t.Helper()

	db := setupTestDB(t)
	queries := New(db)

	return &DB{
		Queries: queries,
		db:      db,
		logger:  logging.NewNoopLogger(),
	}
}

func TestMain(m *testing.M) {
	testCtx = context.Background()

	// Run tests
	code := m.Run()

	os.Exit(code)
}

// Test DBConfig validation

func TestDBConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *DBConfig
		wantErr bool
		errMsg  string
	}{
		{
			name:    "Valid SQLite config",
			config:  &DBConfig{Type: "sqlite", Path: "./test.db"},
			wantErr: false,
		},
		{
			name:    "Valid PostgreSQL config",
			config:  &DBConfig{Type: "postgres", Path: "postgres://localhost/test"},
			wantErr: false,
		},
		{
			name:    "Missing type",
			config:  &DBConfig{Path: "./test.db"},
			wantErr: true,
			errMsg:  "database type is required",
		},
		{
			name:    "Missing path",
			config:  &DBConfig{Type: "sqlite"},
			wantErr: true,
			errMsg:  "database path is required",
		},
		{
			name:    "Unsupported database type",
			config:  &DBConfig{Type: "mysql", Path: "./test.db"},
			wantErr: true,
			errMsg:  "unsupported database type",
		},
		{
			name:    "Empty config",
			config:  &DBConfig{},
			wantErr: true,
			errMsg:  "database type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Error("Validate() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want to contain %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// Test Users

func TestCreateAndGetUser(t *testing.T) {
	db := setupTestDBInstance(t)

	userID := uuid.New().String()
	channel := "telegram"
	channelUserID := "123456"

	// Create user
	params := CreateUserParams{
		ID:            userID,
		Channel:       channel,
		ChannelUserID: channelUserID,
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	user, err := db.CreateUser(testCtx, params)
	if err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	if user.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, user.ID)
	}
	if user.Channel != channel {
		t.Errorf("expected channel %s, got %s", channel, user.Channel)
	}

	// Get user by ID
	gotUser, err := db.GetUserByID(testCtx, userID)
	if err != nil {
		t.Fatalf("failed to get user by ID: %v", err)
	}

	if gotUser.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, gotUser.ID)
	}

	// Get user by channel
	getParams := GetUserByChannelParams{
		Channel:       channel,
		ChannelUserID: channelUserID,
	}
	gotUser2, err := db.GetUserByChannel(testCtx, getParams)
	if err != nil {
		t.Fatalf("failed to get user by channel: %v", err)
	}

	if gotUser2.ID != userID {
		t.Errorf("expected user ID %s, got %s", userID, gotUser2.ID)
	}
}

func TestListUsers(t *testing.T) {
	db := setupTestDBInstance(t)

	// Create multiple users
	for i := 0; i < 3; i++ {
		params := CreateUserParams{
			ID:            uuid.New().String(),
			Channel:       "telegram",
			ChannelUserID: uuid.New().String(),
			CreatedAt:     time.Now().Format(time.RFC3339),
		}
		if _, err := db.CreateUser(testCtx, params); err != nil {
			t.Fatalf("failed to create user: %v", err)
		}
	}

	users, err := db.ListUsers(testCtx)
	if err != nil {
		t.Fatalf("failed to list users: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("expected 3 users, got %d", len(users))
	}
}

func TestDeleteUser(t *testing.T) {
	db := setupTestDBInstance(t)

	userID := uuid.New().String()

	params := CreateUserParams{
		ID:            userID,
		Channel:       "telegram",
		ChannelUserID: "123456",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}

	if _, err := db.CreateUser(testCtx, params); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Delete user
	if err := db.DeleteUser(testCtx, userID); err != nil {
		t.Fatalf("failed to delete user: %v", err)
	}

	// Verify user is deleted
	_, err := db.GetUserByID(testCtx, userID)
	if err == nil {
		t.Error("expected error when getting deleted user, got nil")
	}
}

// Test Sessions

func TestCreateAndGetSession(t *testing.T) {
	db := setupTestDBInstance(t)

	// Create user first
	userID := uuid.New().String()
	userParams := CreateUserParams{
		ID:            userID,
		Channel:       "telegram",
		ChannelUserID: "123456",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	if _, err := db.CreateUser(testCtx, userParams); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	// Create session
	sessionID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	sessionParams := CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	session, err := db.CreateSession(testCtx, sessionParams)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	if session.ID != sessionID {
		t.Errorf("expected session ID %s, got %s", sessionID, session.ID)
	}

	// Get session by ID
	gotSession, err := db.GetSessionByID(testCtx, sessionID)
	if err != nil {
		t.Fatalf("failed to get session by ID: %v", err)
	}

	if gotSession.ID != sessionID {
		t.Errorf("expected session ID %s, got %s", sessionID, gotSession.ID)
	}

	// Get sessions by user ID
	sessions, err := db.GetSessionsByUserID(testCtx, userID)
	if err != nil {
		t.Fatalf("failed to get sessions by user ID: %v", err)
	}

	if len(sessions) != 1 {
		t.Errorf("expected 1 session, got %d", len(sessions))
	}
}

// Test Messages

func TestCreateAndGetMessages(t *testing.T) {
	db := setupTestDBInstance(t)

	// Create user and session
	userID := uuid.New().String()
	userParams := CreateUserParams{
		ID:            userID,
		Channel:       "telegram",
		ChannelUserID: "123456",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	if _, err := db.CreateUser(testCtx, userParams); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	sessionID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	sessionParams := CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := db.CreateSession(testCtx, sessionParams); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Create message
	messageID := uuid.New().String()
	messageParams := CreateMessageParams{
		ID:        messageID,
		SessionID: sessionID,
		Role:      "user",
		Content:   "Hello, world!",
		CreatedAt: now,
	}

	message, err := db.CreateMessage(testCtx, messageParams)
	if err != nil {
		t.Fatalf("failed to create message: %v", err)
	}

	if message.ID != messageID {
		t.Errorf("expected message ID %s, got %s", messageID, message.ID)
	}

	// Get messages by session ID
	messages, err := db.GetMessagesBySessionID(testCtx, sessionID)
	if err != nil {
		t.Fatalf("failed to get messages by session ID: %v", err)
	}

	if len(messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(messages))
	}
}

// Test Tasks

func TestCreateAndUpdateTask(t *testing.T) {
	db := setupTestDBInstance(t)

	// Create user and session
	userID := uuid.New().String()
	userParams := CreateUserParams{
		ID:            userID,
		Channel:       "telegram",
		ChannelUserID: "123456",
		CreatedAt:     time.Now().Format(time.RFC3339),
	}
	if _, err := db.CreateUser(testCtx, userParams); err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	sessionID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)
	sessionParams := CreateSessionParams{
		ID:        sessionID,
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if _, err := db.CreateSession(testCtx, sessionParams); err != nil {
		t.Fatalf("failed to create session: %v", err)
	}

	// Create task
	taskID := uuid.New().String()
	taskParams := CreateTaskParams{
		ID:        taskID,
		SessionID: sessionID,
		Skill:     "test-skill",
		Input:     "test input",
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}

	task, err := db.CreateTask(testCtx, taskParams)
	if err != nil {
		t.Fatalf("failed to create task: %v", err)
	}

	if task.ID != taskID {
		t.Errorf("expected task ID %s, got %s", taskID, task.ID)
	}

	// Update task
	output := "test output"
	updateParams := UpdateTaskParams{
		Output:    sql.NullString{String: output, Valid: true},
		Status:    "completed",
		Error:     sql.NullString{Valid: false},
		UpdatedAt: time.Now().Format(time.RFC3339),
		ID:        taskID,
	}

	updatedTask, err := db.UpdateTask(testCtx, updateParams)
	if err != nil {
		t.Fatalf("failed to update task: %v", err)
	}

	if updatedTask.Status != "completed" {
		t.Errorf("expected status 'completed', got %s", updatedTask.Status)
	}

	if !updatedTask.Output.Valid || updatedTask.Output.String != output {
		t.Errorf("expected output %s, got %v", output, updatedTask.Output)
	}
}

// Test Skills

func TestCreateAndGetSkill(t *testing.T) {
	db := setupTestDBInstance(t)

	skillID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	params := CreateSkillParams{
		ID:          skillID,
		Name:        "test-skill",
		Version:     "1.0.0",
		Location:    "/skills/test",
		Permissions: `["read", "write"]`,
		Metadata:    `{"description": "Test skill"}`,
		CreatedAt:   now,
	}

	skill, err := db.CreateSkill(testCtx, params)
	if err != nil {
		t.Fatalf("failed to create skill: %v", err)
	}

	if skill.ID != skillID {
		t.Errorf("expected skill ID %s, got %s", skillID, skill.ID)
	}

	// Get skill by ID
	gotSkill, err := db.GetSkillByID(testCtx, skillID)
	if err != nil {
		t.Fatalf("failed to get skill by ID: %v", err)
	}

	if gotSkill.Name != "test-skill" {
		t.Errorf("expected skill name 'test-skill', got %s", gotSkill.Name)
	}

	// Get skill by name
	gotSkill2, err := db.GetSkillByName(testCtx, "test-skill")
	if err != nil {
		t.Fatalf("failed to get skill by name: %v", err)
	}

	if gotSkill2.ID != skillID {
		t.Errorf("expected skill ID %s, got %s", skillID, gotSkill2.ID)
	}

	// List skills
	skills, err := db.ListSkills(testCtx)
	if err != nil {
		t.Fatalf("failed to list skills: %v", err)
	}

	if len(skills) != 1 {
		t.Errorf("expected 1 skill, got %d", len(skills))
	}
}

// Test Logs

func TestCreateAndGetLogs(t *testing.T) {
	db := setupTestDBInstance(t)

	logID := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	params := CreateLogParams{
		ID:        logID,
		Level:     "info",
		Source:    "test-source",
		Message:   "Test log message",
		Metadata:  sql.NullString{String: `{"key": "value"}`, Valid: true},
		CreatedAt: now,
	}

	log, err := db.CreateLog(testCtx, params)
	if err != nil {
		t.Fatalf("failed to create log: %v", err)
	}

	if log.ID != logID {
		t.Errorf("expected log ID %s, got %s", logID, log.ID)
	}

	// Get log by ID
	gotLog, err := db.GetLogByID(testCtx, logID)
	if err != nil {
		t.Fatalf("failed to get log by ID: %v", err)
	}

	if gotLog.Level != "info" {
		t.Errorf("expected log level 'info', got %s", gotLog.Level)
	}

	// Get logs by level
	getParams := GetLogsByLevelParams{
		Level: "info",
		Limit: 10,
	}
	logs, err := db.GetLogsByLevel(testCtx, getParams)
	if err != nil {
		t.Fatalf("failed to get logs by level: %v", err)
	}

	if len(logs) != 1 {
		t.Errorf("expected 1 log, got %d", len(logs))
	}
}
