package mappers

import (
	"github.com/atumaikin/nexflow/internal/domain/entity"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// SkillToDomain converts SQLC Skill model to domain Skill entity.
func SkillToDomain(dbSkill *dbmodel.Skill) *entity.Skill {
	if dbSkill == nil {
		return nil
	}

	return &entity.Skill{
		ID:          dbSkill.ID,
		Name:        dbSkill.Name,
		Version:     dbSkill.Version,
		Location:    dbSkill.Location,
		Permissions: dbSkill.Permissions,
		Metadata:    dbSkill.Metadata,
		CreatedAt:   utils.ParseTimeRFC3339(dbSkill.CreatedAt),
	}
}

// SkillToDB converts domain Skill entity to SQLC Skill model.
func SkillToDB(skill *entity.Skill) *dbmodel.Skill {
	if skill == nil {
		return nil
	}

	return &dbmodel.Skill{
		ID:          skill.ID,
		Name:        skill.Name,
		Version:     skill.Version,
		Location:    skill.Location,
		Permissions: skill.Permissions,
		Metadata:    skill.Metadata,
		CreatedAt:   utils.FormatTimeRFC3339(skill.CreatedAt),
	}
}

// SkillsToDomain converts slice of SQLC Skill models to domain Skill entities.
func SkillsToDomain(dbSkills []dbmodel.Skill) []*entity.Skill {
	skills := make([]*entity.Skill, 0, len(dbSkills))
	for i := range dbSkills {
		skills = append(skills, SkillToDomain(&dbSkills[i]))
	}
	return skills
}
