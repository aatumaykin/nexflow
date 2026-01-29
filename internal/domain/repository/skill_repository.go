package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// SkillRepository defines the interface for skill data operations
type SkillRepository interface {
	// Create saves a new skill
	Create(ctx context.Context, skill *entity.Skill) error

	// FindByID retrieves a skill by ID
	FindByID(ctx context.Context, id string) (*entity.Skill, error)

	// FindByName retrieves a skill by name
	FindByName(ctx context.Context, name string) (*entity.Skill, error)

	// List retrieves all skills
	List(ctx context.Context) ([]*entity.Skill, error)

	// Update updates an existing skill
	Update(ctx context.Context, skill *entity.Skill) error

	// Delete removes a skill
	Delete(ctx context.Context, id string) error
}
