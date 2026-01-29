package database

import "context"

// UserRepository defines operations for User entity
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, arg CreateUserParams) (User, error)
	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id string) (User, error)
	// GetByChannel retrieves a user by channel and channel user ID
	GetByChannel(ctx context.Context, arg GetUserByChannelParams) (User, error)
	// List retrieves all users
	List(ctx context.Context) ([]User, error)
	// Delete removes a user
	Delete(ctx context.Context, id string) error
}

// SessionRepository defines operations for Session entity
type SessionRepository interface {
	// Create creates a new session
	Create(ctx context.Context, arg CreateSessionParams) (Session, error)
	// GetByID retrieves a session by ID
	GetByID(ctx context.Context, id string) (Session, error)
	// GetByUserID retrieves all sessions for a user
	GetByUserID(ctx context.Context, userID string) ([]Session, error)
	// Update updates a session
	Update(ctx context.Context, arg UpdateSessionParams) (Session, error)
	// Delete removes a session
	Delete(ctx context.Context, id string) error
}

// MessageRepository defines operations for Message entity
type MessageRepository interface {
	// Create creates a new message
	Create(ctx context.Context, arg CreateMessageParams) (Message, error)
	// GetByID retrieves a message by ID
	GetByID(ctx context.Context, id string) (Message, error)
	// GetBySessionID retrieves all messages for a session
	GetBySessionID(ctx context.Context, sessionID string) ([]Message, error)
	// Delete removes a message
	Delete(ctx context.Context, id string) error
}

// TaskRepository defines operations for Task entity
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, arg CreateTaskParams) (Task, error)
	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, id string) (Task, error)
	// GetBySessionID retrieves all tasks for a session
	GetBySessionID(ctx context.Context, sessionID string) ([]Task, error)
	// Update updates a task
	Update(ctx context.Context, arg UpdateTaskParams) (Task, error)
	// Delete removes a task
	Delete(ctx context.Context, id string) error
}

// SkillRepository defines operations for Skill entity
type SkillRepository interface {
	// Create creates a new skill
	Create(ctx context.Context, arg CreateSkillParams) (Skill, error)
	// GetByID retrieves a skill by ID
	GetByID(ctx context.Context, id string) (Skill, error)
	// GetByName retrieves a skill by name
	GetByName(ctx context.Context, name string) (Skill, error)
	// List retrieves all skills
	List(ctx context.Context) ([]Skill, error)
	// Update updates a skill
	Update(ctx context.Context, arg UpdateSkillParams) (Skill, error)
	// Delete removes a skill
	Delete(ctx context.Context, id string) error
}

// ScheduleRepository defines operations for Schedule entity
type ScheduleRepository interface {
	// Create creates a new schedule
	Create(ctx context.Context, arg CreateScheduleParams) (Schedule, error)
	// GetByID retrieves a schedule by ID
	GetByID(ctx context.Context, id string) (Schedule, error)
	// GetBySkill retrieves all schedules for a skill
	GetBySkill(ctx context.Context, skill string) ([]Schedule, error)
	// List retrieves all schedules
	List(ctx context.Context) ([]Schedule, error)
	// Update updates a schedule
	Update(ctx context.Context, arg UpdateScheduleParams) (Schedule, error)
	// Delete removes a schedule
	Delete(ctx context.Context, id string) error
}

// LogRepository defines operations for Log entity
type LogRepository interface {
	// Create creates a new log entry
	Create(ctx context.Context, arg CreateLogParams) (Log, error)
	// GetByID retrieves a log entry by ID
	GetByID(ctx context.Context, id string) (Log, error)
	// GetByLevel retrieves all logs with a specific level
	GetByLevel(ctx context.Context, arg GetLogsByLevelParams) ([]Log, error)
	// GetBySource retrieves all logs from a specific source
	GetBySource(ctx context.Context, arg GetLogsBySourceParams) ([]Log, error)
	// GetByDateRange retrieves all logs within a date range
	GetByDateRange(ctx context.Context, arg GetLogsByDateRangeParams) ([]Log, error)
	// Delete removes a log entry
	Delete(ctx context.Context, id string) error
	// DeleteOlderThan removes logs older than a specific date
	DeleteOlderThan(ctx context.Context, date string) error
}

// Migration defines operations for database migrations
type Migration interface {
	// Migrate runs all pending migrations
	Migrate(ctx context.Context) error
}

// Closer defines the operation to close database connection
type Closer interface {
	// Close closes the database connection
	Close() error
}
