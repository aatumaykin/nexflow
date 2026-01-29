package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockSkillRuntime(t *testing.T) {
	skillRuntime := NewMockSkillRuntime()

	require.NotNil(t, skillRuntime)
	require.NotNil(t, skillRuntime.ExecuteFunc)
	require.NotNil(t, skillRuntime.ValidateFunc)
	require.NotNil(t, skillRuntime.ListFunc)
	require.NotNil(t, skillRuntime.GetSkillFunc)
}

func TestMockSkillRuntime_Execute(t *testing.T) {
	skillRuntime := NewMockSkillRuntime()
	ctx := context.Background()

	input := map[string]interface{}{
		"test_param": "test_value",
	}

	exec, err := skillRuntime.Execute(ctx, "test-skill", input)

	require.NoError(t, err)
	assert.NotNil(t, exec)
	assert.True(t, exec.Success)
	assert.Contains(t, exec.Output, "mock skill execution for test-skill")
	assert.Contains(t, exec.Output, "test_value")
}

func TestMockSkillRuntime_Validate(t *testing.T) {
	skillRuntime := NewMockSkillRuntime()

	err := skillRuntime.Validate("test-skill")

	assert.NoError(t, err)
}

func TestMockSkillRuntime_List(t *testing.T) {
	skillRuntime := NewMockSkillRuntime()

	skills, err := skillRuntime.List()

	require.NoError(t, err)
	assert.NotNil(t, skills)
	assert.Len(t, skills, 3)
	assert.Contains(t, skills, "mock-skill-1")
	assert.Contains(t, skills, "mock-skill-2")
	assert.Contains(t, skills, "mock-skill-3")
}

func TestMockSkillRuntime_GetSkill(t *testing.T) {
	skillRuntime := NewMockSkillRuntime()

	skill, err := skillRuntime.GetSkill("test-skill")

	require.NoError(t, err)
	assert.NotNil(t, skill)
	assert.Equal(t, "test-skill", skill["name"])
	assert.Equal(t, "1.0.0", skill["version"])
	assert.Contains(t, skill["description"], "test-skill")
}
