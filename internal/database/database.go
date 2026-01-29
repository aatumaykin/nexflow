package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/config"
)

// Database interface defines all database operations
type Database interface {
	// Users
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetUserByID(ctx context.Context, id string) (User, error)
	GetUserByChannel(ctx context.Context, arg GetUserByChannelParams) (User, error)
	ListUsers(ctx context.Context) ([]User, error)
	DeleteUser(ctx context.Context, id string) error

	// Sessions
	CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error)
	GetSessionByID(ctx context.Context, id string) (Session, error)
	GetSessionsByUserID(ctx context.Context, userID string) ([]Session, error)
	UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error)
	DeleteSession(ctx context.Context, id string) error

	// Messages
	CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error)
	GetMessageByID(ctx context.Context, id string) (Message, error)
	GetMessagesBySessionID(ctx context.Context, sessionID string) ([]Message, error)
	DeleteMessage(ctx context.Context, id string) error

	// Tasks
	CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error)
	GetTaskByID(ctx context.Context, id string) (Task, error)
	GetTasksBySessionID(ctx context.Context, sessionID string) ([]Task, error)
	UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error)
	DeleteTask(ctx context.Context, id string) error

	// Skills
	CreateSkill(ctx context.Context, arg CreateSkillParams) (Skill, error)
	GetSkillByID(ctx context.Context, id string) (Skill, error)
	GetSkillByName(ctx context.Context, name string) (Skill, error)
	ListSkills(ctx context.Context) ([]Skill, error)
	UpdateSkill(ctx context.Context, arg UpdateSkillParams) (Skill, error)
	DeleteSkill(ctx context.Context, id string) error

	// Schedules
	CreateSchedule(ctx context.Context, arg CreateScheduleParams) (Schedule, error)
	GetScheduleByID(ctx context.Context, id string) (Schedule, error)
	GetSchedulesBySkill(ctx context.Context, skill string) ([]Schedule, error)
	ListSchedules(ctx context.Context) ([]Schedule, error)
	UpdateSchedule(ctx context.Context, arg UpdateScheduleParams) (Schedule, error)
	DeleteSchedule(ctx context.Context, id string) error

	// Logs
	CreateLog(ctx context.Context, arg CreateLogParams) (Log, error)
	GetLogByID(ctx context.Context, id string) (Log, error)
	GetLogsByLevel(ctx context.Context, arg GetLogsByLevelParams) ([]Log, error)
	GetLogsBySource(ctx context.Context, arg GetLogsBySourceParams) ([]Log, error)
	GetLogsByDateRange(ctx context.Context, arg GetLogsByDateRangeParams) ([]Log, error)
	DeleteLog(ctx context.Context, id string) error
	DeleteLogsOlderThan(ctx context.Context, date string) error

	// Migration
	Migrate(ctx context.Context) error
	Close() error
}

// DBConfig holds database configuration
type DBConfig struct {
	Type string // "sqlite" or "postgres"
	Path string // connection string or file path
}

// DB is the main database implementation
type DB struct {
	*Queries
	db     *sql.DB
	config *DBConfig
}

// NewDatabase creates a new database instance
func NewDatabase(cfg *config.DatabaseConfig) (Database, error) {
	dbConfig := &DBConfig{
		Type: cfg.Type,
		Path: cfg.Path,
	}

	var db *sql.DB
	var err error

	switch cfg.Type {
	case "sqlite":
		db, err = openSQLite(cfg.Path)
	case "postgres":
		db, err = openPostgres(cfg.Path)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", cfg.Type)
	}

	if err != nil {
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	queries := New(db)

	return &DB{
		Queries: queries,
		db:      db,
		config:  dbConfig,
	}, nil
}

// openSQLite opens a SQLite database connection
func openSQLite(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

// openPostgres opens a PostgreSQL database connection
func openPostgres(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres database: %w", err)
	}

	return db, nil
}

// Close closes the database connection
func (d *DB) Close() error {
	return d.db.Close()
}

// GetDB returns the underlying *sql.DB instance
func (d *DB) GetDB() *sql.DB {
	return d.db
}
