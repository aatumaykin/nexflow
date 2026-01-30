package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateSchedule creates a new schedule
func (uc *ScheduleUseCase) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	inputJSON, err := dto.MapToString(req.Input)
	if err != nil {
		return handleScheduleError(err, "failed to marshal input")
	}

	schedule := entity.NewSchedule(req.Skill, req.CronExpression, inputJSON)

	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return handleScheduleError(err, "failed to create schedule")
	}

	uc.logger.Info("schedule created", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return dto.SuccessScheduleResponse(dto.ScheduleDTOFromEntity(schedule)), nil
}
