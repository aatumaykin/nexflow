package skills

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
)

// RuntimeAdapter adapts infrastructure.Executor to ports.SkillRuntime
type RuntimeAdapter struct {
	executor Executor
}

// NewRuntimeAdapter creates a new adapter that implements ports.SkillRuntime
func NewRuntimeAdapter(executor Executor) ports.SkillRuntime {
	return &RuntimeAdapter{
		executor: executor,
	}
}

// Execute implements ports.SkillRuntime.Execute
func (a *RuntimeAdapter) Execute(ctx context.Context, skillName string, input map[string]interface{}) (*ports.SkillExecution, error) {
	// Create execution request
	req := &ExecutionRequest{
		SkillName: skillName,
		Input:     input,
		Context:   make(map[string]interface{}),
	}

	// Execute the skill
	result, err := a.executor.Execute(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("SkillRuntimeAdapter.Execute: %w", err)
	}

	// Convert output to JSON string for compatibility with ports.SkillExecution
	outputJSON, err := json.Marshal(result.Output)
	if err != nil {
		return nil, fmt.Errorf("SkillRuntimeAdapter.Execute: failed to marshal output: %w", err)
	}

	// Create skill execution response
	execution := &ports.SkillExecution{
		Success: result.Success,
		Output:  string(outputJSON),
	}

	// Add error if present
	if result.Error != nil {
		execution.Error = result.Error.Error()
	}

	return execution, nil
}

// Validate implements ports.SkillRuntime.Validate
func (a *RuntimeAdapter) Validate(skillName string) error {
	ctx := context.Background()

	// First check if skill exists by getting metadata
	_, err := a.executor.GetMetadata(ctx, skillName)
	if err != nil {
		return fmt.Errorf("skill '%s' validation failed: %w", skillName, err)
	}

	// Then check if skill is available
	if !a.executor.IsAvailable(ctx, skillName) {
		return fmt.Errorf("skill '%s' is not available", skillName)
	}

	return nil
}

// List implements ports.SkillRuntime.List
func (a *RuntimeAdapter) List() ([]string, error) {
	ctx := context.Background()

	// Get list of skills from executor
	skills, err := a.executor.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("SkillRuntimeAdapter.List: %w", err)
	}

	return skills, nil
}

// GetSkill implements ports.SkillRuntime.GetSkill
func (a *RuntimeAdapter) GetSkill(skillName string) (map[string]interface{}, error) {
	ctx := context.Background()

	// Get metadata from executor
	metadata, err := a.executor.GetMetadata(ctx, skillName)
	if err != nil {
		return nil, fmt.Errorf("SkillRuntimeAdapter.GetSkill: %w", err)
	}

	return metadata, nil
}
