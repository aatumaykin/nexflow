package ports

import (
	"context"
)

// SkillExecution represents the result of a skill execution.
type SkillExecution struct {
	Success bool   `json:"success"`         // Whether the execution was successful
	Output  string `json:"output"`          // Output result in JSON format
	Error   string `json:"error,omitempty"` // Error message if execution failed
}

// SkillRuntime defines the interface for skill execution.
// Skills are reusable components that perform specific tasks.
type SkillRuntime interface {
	// Execute runs a skill with the given input parameters.
	Execute(ctx context.Context, skillName string, input map[string]interface{}) (*SkillExecution, error)

	// Validate checks if a skill is valid (permissions, configuration, etc.).
	Validate(skillName string) error

	// List returns all available skill names.
	List() ([]string, error)

	// GetSkill returns skill details by name.
	GetSkill(skillName string) (map[string]interface{}, error)
}
