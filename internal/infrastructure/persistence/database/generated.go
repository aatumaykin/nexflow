package database

import (
	gendb "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/gen"
)

// Re-export generated types
type (
	Log      = gendb.Log
	Message  = gendb.Message
	Schedule = gendb.Schedule
	Session  = gendb.Session
	Skill    = gendb.Skill
	Task     = gendb.Task
	User     = gendb.User

	CreateLogParams          = gendb.CreateLogParams
	CreateMessageParams      = gendb.CreateMessageParams
	CreateScheduleParams     = gendb.CreateScheduleParams
	CreateSessionParams      = gendb.CreateSessionParams
	CreateSkillParams        = gendb.CreateSkillParams
	CreateTaskParams         = gendb.CreateTaskParams
	CreateUserParams         = gendb.CreateUserParams
	GetLogsByDateRangeParams = gendb.GetLogsByDateRangeParams
	GetLogsByLevelParams     = gendb.GetLogsByLevelParams
	GetLogsBySourceParams    = gendb.GetLogsBySourceParams
	GetUserByChannelParams   = gendb.GetUserByChannelParams
	UpdateScheduleParams     = gendb.UpdateScheduleParams
	UpdateSessionParams      = gendb.UpdateSessionParams
	UpdateSkillParams        = gendb.UpdateSkillParams
	UpdateTaskParams         = gendb.UpdateTaskParams

	DBTX    = gendb.DBTX
	Querier = gendb.Querier
	Queries = gendb.Queries
)

// New creates a new Queries instance
func New(db DBTX) *Queries {
	return gendb.New(db)
}
