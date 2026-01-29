package entity

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSkill(t *testing.T) {
	// Arrange
	name := "my-skill"
	version := "1.0.0"
	location := "/skills/my-skill"
	permissions := []string{"read", "write"}
	metadata := map[string]interface{}{"timeout": 30}

	// Act
	skill := NewSkill(name, version, location, permissions, metadata)

	// Assert
	require.NotEmpty(t, skill.ID)
	assert.Equal(t, name, skill.Name)
	assert.Equal(t, valueobject.Version(version), skill.Version)
	assert.Equal(t, location, skill.Location)
	assert.Equal(t, `["read","write"]`, skill.Permissions) // marshaled
	assert.Equal(t, `{"timeout":30}`, skill.Metadata)      // marshaled
	assert.WithinDuration(t, time.Now(), skill.CreatedAt, time.Second)
}

func TestSkill_RequiresPermission(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{"read"}, nil)

	// Act - currently always returns true
	result := skill.RequiresPermission("read")

	// Assert
	assert.True(t, result)
}

func TestSkill_RequiresSandbox(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{}, nil)

	// Act - empty permissions should not require sandbox
	result := skill.RequiresSandbox()

	// Assert
	assert.False(t, result)
}

func TestSkill_RequiresSandbox_WithDangerousPermissions(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{"shell", "read"}, nil)

	// Act - shell permission requires sandbox
	result := skill.RequiresSandbox()

	// Assert
	assert.True(t, result)
}

func TestSkill_GetTimeout(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{}, nil)

	// Act - currently returns default 30
	timeout := skill.GetTimeout()

	// Assert
	assert.Equal(t, 30, timeout)
}

func TestSkill_HasPermission(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{"read", "write"}, nil)

	// Act - currently always returns true
	result := skill.HasPermission("read")

	// Assert
	assert.True(t, result)
}

func TestSkill_MetadataMap(t *testing.T) {
	// Arrange
	skill := NewSkill("skill", "1.0.0", "/path", []string{}, nil)

	// Act & Assert
	// MetadataMap is not persisted, but can be used in memory
	assert.Nil(t, skill.MetadataMap)
}

func TestSkill_DifferentVersions(t *testing.T) {
	// Arrange
	skill1 := NewSkill("skill", "1.0.0", "/path1", []string{}, nil)
	skill2 := NewSkill("skill", "2.0.0", "/path2", []string{}, nil)

	// Act & Assert
	assert.Equal(t, "skill", skill1.Name)
	assert.Equal(t, "skill", skill2.Name)
	assert.NotEqual(t, skill1.Version, skill2.Version)
	assert.NotEqual(t, skill1.ID, skill2.ID)
}
