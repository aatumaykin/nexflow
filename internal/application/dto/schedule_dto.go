package dto

// ScheduleDTO represents a schedule data transfer object
type ScheduleDTO struct {
	ID             string `json:"id"`
	Skill          string `json:"skill"`           // Name of the skill to execute
	CronExpression string `json:"cron_expression"` // Cron syntax (e.g., "0 * * * *")
	Input          string `json:"input"`           // Input parameters (JSON)
	Enabled        bool   `json:"enabled"`         // Whether schedule is active
	CreatedAt      string `json:"created_at"`      // ISO 8601 format
}

// CreateScheduleRequest represents a request to create a schedule
type CreateScheduleRequest struct {
	Skill          string                 `json:"skill" yaml:"skill"`
	CronExpression string                 `json:"cron_expression" yaml:"cron_expression"`
	Input          map[string]interface{} `json:"input" yaml:"input"`
}

// UpdateScheduleRequest represents a request to update a schedule
type UpdateScheduleRequest struct {
	CronExpression string                 `json:"cron_expression,omitempty" yaml:"cron_expression,omitempty"`
	Input          map[string]interface{} `json:"input,omitempty" yaml:"input,omitempty"`
	Enabled        *bool                  `json:"enabled,omitempty" yaml:"enabled,omitempty"`
}

// ScheduleResponse represents a schedule response
type ScheduleResponse struct {
	Success  bool         `json:"success"`
	Schedule *ScheduleDTO `json:"schedule,omitempty"`
	Error    string       `json:"error,omitempty"`
}

// SchedulesResponse represents a list of schedules response
type SchedulesResponse struct {
	Success   bool           `json:"success"`
	Schedules []*ScheduleDTO `json:"schedules,omitempty"`
	Error     string         `json:"error,omitempty"`
}

// ToggleScheduleRequest represents a request to toggle schedule enabled status
type ToggleScheduleRequest struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
}
