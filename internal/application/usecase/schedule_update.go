package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
)

// UpdateSchedule updates an existing schedule
func (uc *ScheduleUseCase) UpdateSchedule(ctx context.Context, id string, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return handleScheduleError(err, "schedule not found")
	}

	if err := uc.updateScheduleFields(schedule, req); err != nil {
		return handleScheduleError(err, "failed to update schedule fields")
	}

	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return handleScheduleError(err, "failed to update schedule")
	}

	uc.logger.Info("schedule updated", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return dto.SuccessScheduleResponse(dto.ScheduleDTOFromEntity(schedule)), nil
}

// updateScheduleFields updates schedule fields from request
func (uc *ScheduleUseCase) updateScheduleFields(schedule *entity.Schedule, req dto.UpdateScheduleRequest) error {
	// Update cron expression
	if req.CronExpression != "" {
		schedule.CronExpression = valueobject.MustNewCronExpression(req.CronExpression)
	}

	// Update input
	if req.Input != nil {
		inputJSON, err := dto.MapToString(req.Input)
		if err != nil {
			return fmt.Errorf("failed to marshal input: %w", err)
		}
		schedule.Input = inputJSON
	}

	// Update enabled status
	if req.Enabled != nil {
		if *req.Enabled {
			schedule.Enable()
		} else {
			schedule.Disable()
		}
	}

	return nil
}
