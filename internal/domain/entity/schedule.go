package entity

import (
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// Schedule represents a cron-based scheduled task.
// Schedules allow automatic skill execution at specific times defined by cron expressions.
type Schedule struct {
	ID             valueobject.ScheduleID     `json:"id"`              // Unique identifier for the schedule
	Skill          string                     `json:"skill"`           // Name of the skill to execute
	CronExpression valueobject.CronExpression `json:"cron_expression"` // Cron syntax (e.g., "0 * * * *")
	Input          string                     `json:"input"`           // Input parameters in JSON format
	Enabled        bool                       `json:"enabled"`         // Whether the schedule is active
	CreatedAt      time.Time                  `json:"created_at"`      // Timestamp when the schedule was created
}

// NewSchedule creates a new enabled schedule for the specified skill with a cron expression and input.
func NewSchedule(skill, cronExpression, input string) *Schedule {
	return &Schedule{
		ID:             valueobject.ScheduleID(utils.GenerateID()),
		Skill:          skill,
		CronExpression: valueobject.MustNewCronExpression(cronExpression),
		Input:          input,
		Enabled:        true,
		CreatedAt:      utils.Now(),
	}
}

// Enable sets the schedule as enabled.
func (s *Schedule) Enable() {
	s.Enabled = true
}

// Disable sets the schedule as disabled.
func (s *Schedule) Disable() {
	s.Enabled = false
}

// IsEnabled returns true if the schedule is enabled.
func (s *Schedule) IsEnabled() bool {
	return s.Enabled
}

// BelongsToSkill returns true if the schedule belongs to the specified skill.
func (s *Schedule) BelongsToSkill(skill string) bool {
	return s.Skill == skill
}

// GetInput parses and returns the input parameters as a map.
// Returns nil if parsing fails or input is empty.
func (s *Schedule) GetInput() map[string]interface{} {
	return utils.UnmarshalJSONToMap(s.Input)
}
