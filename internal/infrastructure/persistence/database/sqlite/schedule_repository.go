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

var _ repository.ScheduleRepository = (*ScheduleRepository)(nil)

// ScheduleRepository implements repository.ScheduleRepository using SQLC-generated queries
type ScheduleRepository struct {
	db *sql.DB
}

// NewScheduleRepository creates a new ScheduleRepository instance
func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// Create saves a new schedule
func (r *ScheduleRepository) Create(ctx context.Context, schedule *entity.Schedule) error {
	dbSchedule := mappers.ScheduleToDB(schedule)
	if dbSchedule == nil {
		return fmt.Errorf("failed to convert schedule to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO schedules (id, skill, cron_expression, input, enabled, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		dbSchedule.ID, dbSchedule.Skill, dbSchedule.CronExpression, dbSchedule.Input, dbSchedule.Enabled, dbSchedule.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}

// FindByID retrieves a schedule by ID
func (r *ScheduleRepository) FindByID(ctx context.Context, id string) (*entity.Schedule, error) {
	var dbSchedule dbmodel.Schedule

	err := r.db.QueryRowContext(ctx,
		`SELECT id, skill, cron_expression, input, enabled, created_at FROM schedules WHERE id = ? LIMIT 1`,
		id,
	).Scan(&dbSchedule.ID, &dbSchedule.Skill, &dbSchedule.CronExpression, &dbSchedule.Input, &dbSchedule.Enabled, &dbSchedule.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("schedule not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule by id: %w", err)
	}

	return mappers.ScheduleToDomain(&dbSchedule), nil
}

// FindBySkill retrieves all schedules for a skill
func (r *ScheduleRepository) FindBySkill(ctx context.Context, skill string) ([]*entity.Schedule, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, skill, cron_expression, input, enabled, created_at FROM schedules WHERE skill = ? ORDER BY created_at DESC`,
		skill,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find schedules by skill: %w", err)
	}
	defer rows.Close()

	var dbSchedules []dbmodel.Schedule
	for rows.Next() {
		var dbSchedule dbmodel.Schedule

		if err := rows.Scan(&dbSchedule.ID, &dbSchedule.Skill, &dbSchedule.CronExpression, &dbSchedule.Input, &dbSchedule.Enabled, &dbSchedule.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		dbSchedules = append(dbSchedules, dbSchedule)
	}

	return mappers.SchedulesToDomain(dbSchedules), nil
}

// List retrieves all schedules
func (r *ScheduleRepository) List(ctx context.Context) ([]*entity.Schedule, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, skill, cron_expression, input, enabled, created_at FROM schedules ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}
	defer rows.Close()

	var dbSchedules []dbmodel.Schedule
	for rows.Next() {
		var dbSchedule dbmodel.Schedule

		if err := rows.Scan(&dbSchedule.ID, &dbSchedule.Skill, &dbSchedule.CronExpression, &dbSchedule.Input, &dbSchedule.Enabled, &dbSchedule.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		dbSchedules = append(dbSchedules, dbSchedule)
	}

	return mappers.SchedulesToDomain(dbSchedules), nil
}

// Update updates an existing schedule
func (r *ScheduleRepository) Update(ctx context.Context, schedule *entity.Schedule) error {
	dbSchedule := mappers.ScheduleToDB(schedule)
	if dbSchedule == nil {
		return fmt.Errorf("failed to convert schedule to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`UPDATE schedules SET cron_expression = ?, input = ?, enabled = ? WHERE id = ?`,
		dbSchedule.CronExpression, dbSchedule.Input, dbSchedule.Enabled, dbSchedule.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

// Delete removes a schedule
func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM schedules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("schedule not found: %s", id)
	}

	return nil
}

// FindEnabled retrieves all enabled schedules
func (r *ScheduleRepository) FindEnabled(ctx context.Context) ([]*entity.Schedule, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, skill, cron_expression, input, enabled, created_at FROM schedules WHERE enabled = 1 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find enabled schedules: %w", err)
	}
	defer rows.Close()

	var dbSchedules []dbmodel.Schedule
	for rows.Next() {
		var dbSchedule dbmodel.Schedule

		if err := rows.Scan(&dbSchedule.ID, &dbSchedule.Skill, &dbSchedule.CronExpression, &dbSchedule.Input, &dbSchedule.Enabled, &dbSchedule.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan schedule: %w", err)
		}

		dbSchedules = append(dbSchedules, dbSchedule)
	}

	return mappers.SchedulesToDomain(dbSchedules), nil
}
