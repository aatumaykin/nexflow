package mappers

import (
	"github.com/atumaikin/nexflow/internal/domain/entity"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// ScheduleToDomain converts SQLC Schedule model to domain Schedule entity.
func ScheduleToDomain(dbSchedule *dbmodel.Schedule) *entity.Schedule {
	if dbSchedule == nil {
		return nil
	}

	return &entity.Schedule{
		ID:             dbSchedule.ID,
		Skill:          dbSchedule.Skill,
		CronExpression: dbSchedule.CronExpression,
		Input:          dbSchedule.Input,
		Enabled:        dbSchedule.Enabled == 1,
		CreatedAt:      utils.ParseTimeRFC3339(dbSchedule.CreatedAt),
	}
}

// ScheduleToDB converts domain Schedule entity to SQLC Schedule model.
func ScheduleToDB(schedule *entity.Schedule) *dbmodel.Schedule {
	if schedule == nil {
		return nil
	}

	enabled := int64(0)
	if schedule.Enabled {
		enabled = 1
	}

	return &dbmodel.Schedule{
		ID:             schedule.ID,
		Skill:          schedule.Skill,
		CronExpression: schedule.CronExpression,
		Input:          schedule.Input,
		Enabled:        enabled,
		CreatedAt:      utils.FormatTimeRFC3339(schedule.CreatedAt),
	}
}

// SchedulesToDomain converts slice of SQLC Schedule models to domain Schedule entities.
func SchedulesToDomain(dbSchedules []dbmodel.Schedule) []*entity.Schedule {
	schedules := make([]*entity.Schedule, 0, len(dbSchedules))
	for i := range dbSchedules {
		schedules = append(schedules, ScheduleToDomain(&dbSchedules[i]))
	}
	return schedules
}
