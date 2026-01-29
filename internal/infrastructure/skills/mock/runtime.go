package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
)

// MockSkillRuntime is a mock implementation of SkillRuntime for testing
type MockSkillRuntime struct {
	ExecuteFunc  func(context.Context, string, map[string]interface{}) (*ports.SkillExecution, error)
	ValidateFunc func(string) error
	ListFunc     func() ([]string, error)
	GetSkillFunc func(string) (map[string]interface{}, error)
}

// NewMockSkillRuntime creates a new mock skill runtime
func NewMockSkillRuntime() *MockSkillRuntime {
	return &MockSkillRuntime{
		ExecuteFunc: func(ctx context.Context, skillName string, input map[string]interface{}) (*ports.SkillExecution, error) {
			// Simulate skill execution
			output := fmt.Sprintf(`{"result": "mock skill execution for %s", "input": %v}`, skillName, input)
			return &ports.SkillExecution{
				Success: true,
				Output:  output,
			}, nil
		},
		ValidateFunc: func(skillName string) error {
			// Mock validation - all skills are valid
			return nil
		},
		ListFunc: func() ([]string, error) {
			// Mock skill list
			return []string{"mock-skill-1", "mock-skill-2", "mock-skill-3"}, nil
		},
		GetSkillFunc: func(skillName string) (map[string]interface{}, error) {
			// Mock skill details
			return map[string]interface{}{
				"name":        skillName,
				"version":     "1.0.0",
				"description": fmt.Sprintf("Mock skill: %s", skillName),
			}, nil
		},
	}
}

// Execute implements SkillRuntime interface
func (m *MockSkillRuntime) Execute(ctx context.Context, skillName string, input map[string]interface{}) (*ports.SkillExecution, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, skillName, input)
	}
	return nil, fmt.Errorf("ExecuteFunc not set")
}

// Validate implements SkillRuntime interface
func (m *MockSkillRuntime) Validate(skillName string) error {
	if m.ValidateFunc != nil {
		return m.ValidateFunc(skillName)
	}
	return fmt.Errorf("ValidateFunc not set")
}

// List implements SkillRuntime interface
func (m *MockSkillRuntime) List() ([]string, error) {
	if m.ListFunc != nil {
		return m.ListFunc()
	}
	return nil, fmt.Errorf("ListFunc not set")
}

// GetSkill implements SkillRuntime interface
func (m *MockSkillRuntime) GetSkill(skillName string) (map[string]interface{}, error) {
	if m.GetSkillFunc != nil {
		return m.GetSkillFunc(skillName)
	}
	return nil, fmt.Errorf("GetSkillFunc not set")
}
