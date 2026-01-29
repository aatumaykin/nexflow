package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.SkillRepository = (*SkillRepository)(nil)

// SkillRepository implements repository.SkillRepository using SQLC-generated queries
type SkillRepository struct {
	db *sql.DB
}

// NewSkillRepository creates a new SkillRepository instance
func NewSkillRepository(db *sql.DB) *SkillRepository {
	return &SkillRepository{db: db}
}

// Create saves a new skill
func (r *SkillRepository) Create(ctx context.Context, skill *entity.Skill) error {
	dbSkill := mappers.SkillToDB(skill)
	if dbSkill == nil {
		return fmt.Errorf("failed to convert skill to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO skills (id, name, version, location, permissions, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		dbSkill.ID, dbSkill.Name, dbSkill.Version, dbSkill.Location, dbSkill.Permissions, dbSkill.Metadata, dbSkill.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create skill: %w", err)
	}

	return nil
}

// FindByID retrieves a skill by ID
func (r *SkillRepository) FindByID(ctx context.Context, id string) (*entity.Skill, error) {
	var dbSkill dbmodel.Skill

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, version, location, permissions, metadata, created_at FROM skills WHERE id = ? LIMIT 1`,
		id,
	).Scan(&dbSkill.ID, &dbSkill.Name, &dbSkill.Version, &dbSkill.Location, &dbSkill.Permissions, &dbSkill.Metadata, &dbSkill.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find skill by id: %w", err)
	}

	return mappers.SkillToDomain(&dbSkill), nil
}

// FindByName retrieves a skill by name
func (r *SkillRepository) FindByName(ctx context.Context, name string) (*entity.Skill, error) {
	var dbSkill dbmodel.Skill

	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, version, location, permissions, metadata, created_at FROM skills WHERE name = ? LIMIT 1`,
		name,
	).Scan(&dbSkill.ID, &dbSkill.Name, &dbSkill.Version, &dbSkill.Location, &dbSkill.Permissions, &dbSkill.Metadata, &dbSkill.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("skill not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find skill by name: %w", err)
	}

	return mappers.SkillToDomain(&dbSkill), nil
}

// List retrieves all skills
func (r *SkillRepository) List(ctx context.Context) ([]*entity.Skill, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, version, location, permissions, metadata, created_at FROM skills ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list skills: %w", err)
	}
	defer rows.Close()

	var dbSkills []dbmodel.Skill
	for rows.Next() {
		var dbSkill dbmodel.Skill

		if err := rows.Scan(&dbSkill.ID, &dbSkill.Name, &dbSkill.Version, &dbSkill.Location, &dbSkill.Permissions, &dbSkill.Metadata, &dbSkill.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan skill: %w", err)
		}

		dbSkills = append(dbSkills, dbSkill)
	}

	return mappers.SkillsToDomain(dbSkills), nil
}

// Update updates an existing skill
func (r *SkillRepository) Update(ctx context.Context, skill *entity.Skill) error {
	dbSkill := mappers.SkillToDB(skill)
	if dbSkill == nil {
		return fmt.Errorf("failed to convert skill to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE skills SET version = ?, location = ?, permissions = ?, metadata = ? WHERE id = ?`,
		dbSkill.Version, dbSkill.Location, dbSkill.Permissions, dbSkill.Metadata, dbSkill.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update skill: %w", err)
	}

	return nil
}

// Delete removes a skill
func (r *SkillRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM skills WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete skill: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("skill not found: %s", id)
	}

	return nil
}
