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

var _ repository.ScheduleRepository = (*ScheduleRepository)(nil)

type ScheduleRepository struct {
	queries *database.Queries
}

func NewScheduleRepository(queries *database.Queries) *ScheduleRepository {
	return &ScheduleRepository{queries: queries}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *entity.Schedule) error {
	dbSchedule := mappers.ScheduleToDB(schedule)
	if dbSchedule == nil {
		return fmt.Errorf("failed to convert schedule to db model")
	}

	_, err := r.queries.CreateSchedule(ctx, database.CreateScheduleParams{
		ID:             dbSchedule.ID,
		Skill:          dbSchedule.Skill,
		CronExpression: dbSchedule.CronExpression,
		Input:          dbSchedule.Input,
		Enabled:        dbSchedule.Enabled,
		CreatedAt:      dbSchedule.CreatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) FindByID(ctx context.Context, id string) (*entity.Schedule, error) {
	dbSchedule, err := r.queries.GetScheduleByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("schedule not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find schedule by id: %w", err)
	}

	return mappers.ScheduleToDomain(&dbSchedule), nil
}

func (r *ScheduleRepository) FindBySkill(ctx context.Context, skill string) ([]*entity.Schedule, error) {
	dbSchedules, err := r.queries.GetSchedulesBySkill(ctx, skill)
	if err != nil {
		return nil, fmt.Errorf("failed to find schedules by skill: %w", err)
	}

	return mappers.SchedulesToDomain(dbSchedules), nil
}

func (r *ScheduleRepository) List(ctx context.Context) ([]*entity.Schedule, error) {
	dbSchedules, err := r.queries.ListSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list schedules: %w", err)
	}

	return mappers.SchedulesToDomain(dbSchedules), nil
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule *entity.Schedule) error {
	dbSchedule := mappers.ScheduleToDB(schedule)
	if dbSchedule == nil {
		return fmt.Errorf("failed to convert schedule to db model")
	}

	_, err := r.queries.UpdateSchedule(ctx, database.UpdateScheduleParams{
		CronExpression: dbSchedule.CronExpression,
		Input:          dbSchedule.Input,
		Enabled:        dbSchedule.Enabled,
		ID:             dbSchedule.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to update schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetScheduleByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("schedule not found: %s", id)
		}
		return fmt.Errorf("failed to check schedule existence: %w", err)
	}

	err = r.queries.DeleteSchedule(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete schedule: %w", err)
	}

	return nil
}

func (r *ScheduleRepository) FindEnabled(ctx context.Context) ([]*entity.Schedule, error) {
	dbSchedules, err := r.queries.ListSchedules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find enabled schedules: %w", err)
	}

	var enabled []*entity.Schedule
	for _, s := range dbSchedules {
		if s.Enabled == 1 {
			enabled = append(enabled, mappers.ScheduleToDomain(&s))
		}
	}

	return enabled, nil
}
