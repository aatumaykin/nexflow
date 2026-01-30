package mappers

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSkillToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSkill     *dbmodel.Skill
		expected    *entity.Skill
		expectedNil bool
	}{
		{
			name: "Valid skill with metadata",
			dbSkill: &dbmodel.Skill{
				ID:          "test-id",
				Name:        "test-skill",
				Version:     "1.0.0",
				Location:    "/skills/test",
				Permissions: `["read", "write"]`,
				Metadata:    `{"description": "Test skill"}`,
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expected: &entity.Skill{
				ID:          valueobject.SkillID("test-id"),
				Name:        "test-skill",
				Version:     valueobject.MustNewVersion("1.0.0"),
				Location:    "/skills/test",
				Permissions: `["read", "write"]`,
				Metadata:    `{"description": "Test skill"}`,
				CreatedAt:   time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid skill without metadata",
			dbSkill: &dbmodel.Skill{
				ID:          "test-id-2",
				Name:        "test-skill-2",
				Version:     "2.0.0",
				Location:    "/skills/test-2",
				Permissions: `["read", "write", "execute"]`,
				Metadata:    "",
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expected: &entity.Skill{
				ID:          valueobject.SkillID("test-id-2"),
				Name:        "test-skill-2",
				Version:     valueobject.MustNewVersion("2.0.0"),
				Location:    "/skills/test-2",
				Permissions: `["read", "write", "execute"]`,
				Metadata:    "",
				CreatedAt:   time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbSkill:     nil,
			expected:    nil,
			expectedNil: true,
		},
		{
			name: "Skill with empty metadata",
			dbSkill: &dbmodel.Skill{
				ID:          "test-id-3",
				Name:        "test-skill-3",
				Version:     "3.0.0",
				Location:    "/skills/test-3",
				Permissions: `["read"]`,
				Metadata:    "",
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expected: &entity.Skill{
				ID:          valueobject.SkillID("test-id-3"),
				Name:        "test-skill-3",
				Version:     valueobject.MustNewVersion("3.0.0"),
				Location:    "/skills/test-3",
				Permissions: `["read"]`,
				Metadata:    "",
				CreatedAt:   time.Now(),
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SkillToDomain(tt.dbSkill)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Version, result.Version)
			assert.Equal(t, tt.expected.Location, result.Location)
			assert.Equal(t, tt.expected.Permissions, result.Permissions)
			assert.Equal(t, tt.expected.Metadata, result.Metadata)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
		})
	}
}

func TestSkillToDB(t *testing.T) {
	tests := []struct {
		name        string
		skill       *entity.Skill
		expected    *dbmodel.Skill
		expectedNil bool
	}{
		{
			name: "Valid skill with metadata",
			skill: &entity.Skill{
				ID:          valueobject.SkillID("test-id"),
				Name:        "test-skill",
				Version:     valueobject.MustNewVersion("1.0.0"),
				Location:    "/skills/test",
				Permissions: `["read", "write"]`,
				Metadata:    `{"description": "Test skill"}`,
				CreatedAt:   time.Now(),
			},
			expected: &dbmodel.Skill{
				ID:          "test-id",
				Name:        "test-skill",
				Version:     "1.0.0",
				Location:    "/skills/test",
				Permissions: `["read", "write"]`,
				Metadata:    `{"description": "Test skill"}`,
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid skill without metadata",
			skill: &entity.Skill{
				ID:          valueobject.SkillID("test-id-2"),
				Name:        "test-skill-2",
				Version:     valueobject.MustNewVersion("2.0.0"),
				Location:    "/skills/test-2",
				Permissions: `["read", "write", "execute"]`,
				Metadata:    "",
				CreatedAt:   time.Now(),
			},
			expected: &dbmodel.Skill{
				ID:          "test-id-2",
				Name:        "test-skill-2",
				Version:     "2.0.0",
				Location:    "/skills/test-2",
				Permissions: `["read", "write", "execute"]`,
				Metadata:    "",
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			skill:       nil,
			expected:    nil,
			expectedNil: true,
		},
		{
			name: "Skill with empty metadata",
			skill: &entity.Skill{
				ID:          valueobject.SkillID("test-id-3"),
				Name:        "test-skill-3",
				Version:     valueobject.MustNewVersion("3.0.0"),
				Location:    "/skills/test-3",
				Permissions: `["read"]`,
				Metadata:    "",
				CreatedAt:   time.Now(),
			},
			expected: &dbmodel.Skill{
				ID:          "test-id-3",
				Name:        "test-skill-3",
				Version:     "3.0.0",
				Location:    "/skills/test-3",
				Permissions: `["read"]`,
				Metadata:    "",
				CreatedAt:   time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SkillToDB(tt.skill)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Version, result.Version)
			assert.Equal(t, tt.expected.Location, result.Location)
			assert.Equal(t, tt.expected.Permissions, result.Permissions)
			assert.Equal(t, tt.expected.Metadata, result.Metadata)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestSkillsToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbSkills    []dbmodel.Skill
		expectedLen int
	}{
		{
			name: "Multiple skills",
			dbSkills: []dbmodel.Skill{
				{
					ID:          "skill-1",
					Name:        "skill-1",
					Version:     "1.0.0",
					Location:    "/skills/1",
					Permissions: `["read"]`,
					Metadata:    "",
					CreatedAt:   time.Now().Format(time.RFC3339),
				},
				{
					ID:          "skill-2",
					Name:        "skill-2",
					Version:     "2.0.0",
					Location:    "/skills/2",
					Permissions: `["read", "write"]`,
					Metadata:    "",
					CreatedAt:   time.Now().Format(time.RFC3339),
				},
				{
					ID:          "skill-3",
					Name:        "skill-3",
					Version:     "3.0.0",
					Location:    "/skills/3",
					Permissions: `["read", "write", "execute"]`,
					Metadata:    "",
					CreatedAt:   time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbSkills:    []dbmodel.Skill{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbSkills:    nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SkillsToDomain(tt.dbSkills)

			if tt.dbSkills == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, skill := range result {
				assert.Equal(t, tt.dbSkills[i].ID, string(skill.ID))
				assert.Equal(t, tt.dbSkills[i].Name, skill.Name)
				assert.Equal(t, tt.dbSkills[i].Permissions, skill.Permissions)
			}
		})
	}
}
