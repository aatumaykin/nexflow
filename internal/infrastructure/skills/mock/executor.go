package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/infrastructure/skills"
)

// MockExecutor is a mock implementation of the skills Executor interface
type MockExecutor struct {
	name       string
	skills     map[string]mockSkill
	executions []mockExecution
}

type mockSkill struct {
	name      string
	metadata  map[string]interface{}
	available bool
}

type mockExecution struct {
	request *skills.ExecutionRequest
	result  *skills.ExecutionResult
}

// NewMockExecutor creates a new mock skills executor
func NewMockExecutor() *MockExecutor {
	return &MockExecutor{
		name:       "mock",
		skills:     make(map[string]mockSkill),
		executions: make([]mockExecution, 0),
	}
}

// Name returns the name of the executor
func (e *MockExecutor) Name() string {
	return e.name
}

// Execute runs a skill with the given input
func (e *MockExecutor) Execute(ctx context.Context, req *skills.ExecutionRequest) (*skills.ExecutionResult, error) {
	skill, exists := e.skills[req.SkillName]
	if !exists {
		return nil, fmt.Errorf("skill not found: %s", req.SkillName)
	}

	if !skill.available {
		return nil, fmt.Errorf("skill is not available: %s", req.SkillName)
	}

	// Create a default successful response
	result := &skills.ExecutionResult{
		Success: true,
		Output: map[string]interface{}{
			"skill":  req.SkillName,
			"input":  req.Input,
			"result": "Mock execution successful",
		},
		Error: nil,
		Metadata: map[string]interface{}{
			"executor": e.name,
		},
	}

	// Store the execution for inspection
	e.executions = append(e.executions, mockExecution{
		request: req,
		result:  result,
	})

	return result, nil
}

// List returns a list of available skills
func (e *MockExecutor) List(ctx context.Context) ([]string, error) {
	var skillNames []string
	for name := range e.skills {
		skillNames = append(skillNames, name)
	}

	return skillNames, nil
}

// GetMetadata returns metadata for a specific skill
func (e *MockExecutor) GetMetadata(ctx context.Context, skillName string) (map[string]interface{}, error) {
	skill, exists := e.skills[skillName]
	if !exists {
		return nil, fmt.Errorf("skill not found: %s", skillName)
	}

	return skill.metadata, nil
}

// IsAvailable checks if a skill is available for execution
func (e *MockExecutor) IsAvailable(ctx context.Context, skillName string) bool {
	skill, exists := e.skills[skillName]
	if !exists {
		return false
	}

	return skill.available
}

// AddSkill adds a skill to the mock executor
func (e *MockExecutor) AddSkill(name string, metadata map[string]interface{}, available bool) {
	e.skills[name] = mockSkill{
		name:      name,
		metadata:  metadata,
		available: available,
	}
}

// RemoveSkill removes a skill from the mock executor
func (e *MockExecutor) RemoveSkill(name string) {
	delete(e.skills, name)
}

// SetSkillAvailability sets the availability status of a skill
func (e *MockExecutor) SetSkillAvailability(name string, available bool) {
	if skill, exists := e.skills[name]; exists {
		skill.available = available
		e.skills[name] = skill
	}
}

// GetExecutions returns all recorded executions
func (e *MockExecutor) GetExecutions() []mockExecution {
	executions := make([]mockExecution, len(e.executions))
	copy(executions, e.executions)
	return executions
}

// ClearExecutions clears all recorded executions
func (e *MockExecutor) ClearExecutions() {
	e.executions = make([]mockExecution, 0)
}

// SetCustomResult sets a custom result for the next execution of a specific skill
func (e *MockExecutor) SetCustomResult(skillName string, result *skills.ExecutionResult) {
	// In a more sophisticated implementation, we could store custom results per skill
	// For now, this is a placeholder for that functionality
}
