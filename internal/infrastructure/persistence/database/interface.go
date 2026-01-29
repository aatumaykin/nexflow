package database

import "context"

// Database is a legacy composite interface for backward compatibility.
// It combines all repository methods into a single interface.
//
// DEPRECATED: Prefer using specific repository interfaces (UserRepository, SessionRepository, etc.)
// directly instead of this monolithic interface following Interface Segregation Principle.
//
// New code should use individual repository interfaces:
// - UserRepository for user operations
// - SessionRepository for session operations
// - MessageRepository for message operations
// - TaskRepository for task operations
// - SkillRepository for skill operations
// - ScheduleRepository for schedule operations
// - LogRepository for log operations
// - Migration for database migrations
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
