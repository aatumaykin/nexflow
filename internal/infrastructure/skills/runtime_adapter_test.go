package skills

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockExecutor implements Executor interface for testing
type mockExecutor struct {
	name       string
	skills     map[string]mockSkill
	results    map[string]*ExecutionResult
	err        error
	executions []ExecutionRequest
}

type mockSkill struct {
	name      string
	metadata  map[string]interface{}
	available bool
}

func (m *mockExecutor) Name() string {
	return m.name
}

func (m *mockExecutor) Execute(ctx context.Context, req *ExecutionRequest) (*ExecutionResult, error) {
	m.executions = append(m.executions, *req)

	if m.err != nil {
		return nil, m.err
	}

	if result, ok := m.results[req.SkillName]; ok {
		return result, nil
	}

	// Default success response
	return &ExecutionResult{
		Success: true,
		Output: map[string]interface{}{
			"skill":  req.SkillName,
			"input":  req.Input,
			"result": "mock execution successful",
		},
		Error:    nil,
		Metadata: map[string]interface{}{"executor": m.name},
	}, nil
}

func (m *mockExecutor) List(ctx context.Context) ([]string, error) {
	if m.err != nil {
		return nil, m.err
	}

	skillNames := make([]string, 0, len(m.skills))
	for name := range m.skills {
		skillNames = append(skillNames, name)
	}
	return skillNames, nil
}

func (m *mockExecutor) GetMetadata(ctx context.Context, skillName string) (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}

	skill, exists := m.skills[skillName]
	if !exists {
		return nil, errors.New("skill not found")
	}

	return skill.metadata, nil
}

func (m *mockExecutor) IsAvailable(ctx context.Context, skillName string) bool {
	skill, exists := m.skills[skillName]
	if !exists {
		return false
	}
	return skill.available
}

func (m *mockExecutor) addSkill(name string, metadata map[string]interface{}, available bool) {
	if m.skills == nil {
		m.skills = make(map[string]mockSkill)
	}
	m.skills[name] = mockSkill{
		name:      name,
		metadata:  metadata,
		available: available,
	}
}

func (m *mockExecutor) setResult(skillName string, result *ExecutionResult) {
	if m.results == nil {
		m.results = make(map[string]*ExecutionResult)
	}
	m.results[skillName] = result
}

func TestNewRuntimeAdapter(t *testing.T) {
	executor := &mockExecutor{name: "test"}
	adapter := NewRuntimeAdapter(executor)

	require.NotNil(t, adapter)
	assert.IsType(t, &RuntimeAdapter{}, adapter)
}

func TestRuntimeAdapter_Execute(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("test-skill", map[string]interface{}{
		"name":        "test-skill",
		"description": "A test skill",
		"version":     "1.0.0",
	}, true)

	adapter := NewRuntimeAdapter(executor)

	ctx := context.Background()
	input := map[string]interface{}{
		"param1": "value1",
		"param2": 42,
	}

	exec, err := adapter.Execute(ctx, "test-skill", input)

	require.NoError(t, err)
	assert.NotNil(t, exec)
	assert.True(t, exec.Success)
	assert.Contains(t, exec.Output, "test-skill")
	assert.Empty(t, exec.Error)
}

func TestRuntimeAdapter_Execute_WithError(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.setResult("test-skill", &ExecutionResult{
		Success: false,
		Output:  map[string]interface{}{},
		Error:   errors.New("skill execution failed"),
	})

	adapter := NewRuntimeAdapter(executor)

	ctx := context.Background()
	input := map[string]interface{}{}

	exec, err := adapter.Execute(ctx, "test-skill", input)

	require.NoError(t, err)
	assert.NotNil(t, exec)
	assert.False(t, exec.Success)
	assert.Equal(t, "skill execution failed", exec.Error)
}

func TestRuntimeAdapter_Execute_Error(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
		err:  assert.AnError,
	}
	executor.addSkill("test-skill", map[string]interface{}{}, true)

	adapter := NewRuntimeAdapter(executor)

	ctx := context.Background()
	input := map[string]interface{}{}

	exec, err := adapter.Execute(ctx, "test-skill", input)

	require.Error(t, err)
	assert.Nil(t, exec)
	assert.Contains(t, err.Error(), "SkillRuntimeAdapter.Execute")
}

func TestRuntimeAdapter_Validate(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("test-skill", map[string]interface{}{
		"name":        "test-skill",
		"description": "A test skill",
		"version":     "1.0.0",
	}, true)

	adapter := NewRuntimeAdapter(executor)

	err := adapter.Validate("test-skill")

	assert.NoError(t, err)
}

func TestRuntimeAdapter_Validate_NotAvailable(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("test-skill", map[string]interface{}{
		"name": "test-skill",
	}, false) // Not available

	adapter := NewRuntimeAdapter(executor)

	err := adapter.Validate("test-skill")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "is not available")
}

func TestRuntimeAdapter_Validate_NotFound(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("other-skill", map[string]interface{}{
		"name": "other-skill",
	}, true)

	adapter := NewRuntimeAdapter(executor)

	err := adapter.Validate("test-skill")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")
}

func TestRuntimeAdapter_List(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("skill-1", map[string]interface{}{"name": "skill-1"}, true)
	executor.addSkill("skill-2", map[string]interface{}{"name": "skill-2"}, true)
	executor.addSkill("skill-3", map[string]interface{}{"name": "skill-3"}, true)

	adapter := NewRuntimeAdapter(executor)

	skills, err := adapter.List()

	require.NoError(t, err)
	assert.NotNil(t, skills)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "skill-1")
	assert.Contains(t, skills, "skill-2")
	assert.Contains(t, skills, "skill-3")
}

func TestRuntimeAdapter_List_Error(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
		err:  assert.AnError,
	}
	adapter := NewRuntimeAdapter(executor)

	skills, err := adapter.List()

	assert.Error(t, err)
	assert.Nil(t, skills)
	assert.Contains(t, err.Error(), "SkillRuntimeAdapter.List")
}

func TestRuntimeAdapter_GetSkill(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	expectedMetadata := map[string]interface{}{
		"name":        "test-skill",
		"description": "A test skill",
		"version":     "1.0.0",
		"author":      "test-author",
	}
	executor.addSkill("test-skill", expectedMetadata, true)

	adapter := NewRuntimeAdapter(executor)

	skill, err := adapter.GetSkill("test-skill")

	require.NoError(t, err)
	assert.NotNil(t, skill)
	assert.Equal(t, "test-skill", skill["name"])
	assert.Equal(t, "A test skill", skill["description"])
	assert.Equal(t, "1.0.0", skill["version"])
	assert.Equal(t, "test-author", skill["author"])
}

func TestRuntimeAdapter_GetSkill_NotFound(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	executor.addSkill("other-skill", map[string]interface{}{"name": "other-skill"}, true)

	adapter := NewRuntimeAdapter(executor)

	skill, err := adapter.GetSkill("test-skill")

	assert.Error(t, err)
	assert.Nil(t, skill)
	assert.Contains(t, err.Error(), "SkillRuntimeAdapter.GetSkill")
}

func TestRuntimeAdapter_Execute_OutputFormat(t *testing.T) {
	executor := &mockExecutor{
		name: "test",
	}
	output := map[string]interface{}{
		"result": "success",
		"data":   map[string]interface{}{"key": "value"},
		"number": 42,
	}
	executor.setResult("test-skill", &ExecutionResult{
		Success: true,
		Output:  output,
		Error:   nil,
	})

	adapter := NewRuntimeAdapter(executor)

	ctx := context.Background()
	input := map[string]interface{}{}

	exec, err := adapter.Execute(ctx, "test-skill", input)

	require.NoError(t, err)
	assert.NotNil(t, exec)
	assert.True(t, exec.Success)
	assert.Contains(t, exec.Output, "success")
	// Output should be valid JSON
	assert.Contains(t, exec.Output, "{")
	assert.Contains(t, exec.Output, "}")
}
