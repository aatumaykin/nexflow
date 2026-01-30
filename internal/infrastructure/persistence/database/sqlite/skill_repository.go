package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.SkillRepository = (*SkillRepository)(nil)

type SkillRepository struct {
	queries *database.Queries
}

func NewSkillRepository(queries *database.Queries) *SkillRepository {
	return &SkillRepository{queries: queries}
}

func (r *SkillRepository) Create(ctx context.Context, skill *entity.Skill) error {
	dbSkill := mappers.SkillToDB(skill)
	if dbSkill == nil {
		return fmt.Errorf("failed to convert skill to db model")
	}

	_, err := r.queries.CreateSkill(ctx, database.CreateSkillParams{
		ID:          dbSkill.ID,
		Name:        dbSkill.Name,
		Version:     dbSkill.Version,
		Location:    dbSkill.Location,
		Permissions: dbSkill.Permissions,
		Metadata:    dbSkill.Metadata,
		CreatedAt:   dbSkill.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}

func (r *SkillRepository) FindByID(ctx context.Context, id string) (*entity.Skill, error) {
	dbSkill, err := r.queries.GetSkillByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find skill by id: %w", err)
	}

	return mappers.SkillToDomain(&dbSkill), nil
}

func (r *SkillRepository) FindByName(ctx context.Context, name string) (*entity.Skill, error) {
	dbSkill, err := r.queries.GetSkillByName(ctx, name)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find skill by name: %w", err)
	}

	return mappers.SkillToDomain(&dbSkill), nil
}

func (r *SkillRepository) List(ctx context.Context) ([]*entity.Skill, error) {
	dbSkills, err := r.queries.ListSkills(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}

	return mappers.SkillsToDomain(dbSkills), nil
}

func (r *SkillRepository) Update(ctx context.Context, skill *entity.Skill) error {
	dbSkill := mappers.SkillToDB(skill)
	if dbSkill == nil {
		return fmt.Errorf("failed to convert skill to db model")
	}

	_, err := r.queries.UpdateSkill(ctx, database.UpdateSkillParams{
		Version:     dbSkill.Version,
		Location:    dbSkill.Location,
		Permissions: dbSkill.Permissions,
		Metadata:    dbSkill.Metadata,
		ID:          dbSkill.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	return nil
}

func (r *SkillRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetSkillByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("skill not found: %s", id)
		}
		return fmt.Errorf("failed to check skill existence: %w", err)
	}

	err = r.queries.DeleteSkill(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	return nil
}
