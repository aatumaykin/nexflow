package skills

import (
	"context"
)

// ExecutionRequest represents a request to execute a skill
type ExecutionRequest struct {
	SkillName string                 // Name of the skill to execute
	Input     map[string]interface{} // Input parameters for the skill
	Context   map[string]interface{} // Additional context (user ID, session ID, etc.)
}

// ExecutionResult represents the result of skill execution
type ExecutionResult struct {
	Success  bool                   // Whether the execution was successful
	Output   map[string]interface{} // Output from the skill
	Error    error                  // Error if execution failed
	Metadata map[string]interface{} // Additional metadata (execution time, etc.)
}

// Executor defines the interface for skill execution
type Executor interface {
	// Name returns the name of the executor
	Name() string

	// Execute runs a skill with the given input
	Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResult, error)

	// List returns a list of available skills
	List(ctx context.Context) ([]string, error)

	// GetMetadata returns metadata for a specific skill
	GetMetadata(ctx context.Context, skillName string) (map[string]interface{}, error)

	// IsAvailable checks if a skill is available for execution
	IsAvailable(ctx context.Context, skillName string) bool
}
