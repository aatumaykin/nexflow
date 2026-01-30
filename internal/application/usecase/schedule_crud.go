package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// DeleteSchedule deletes a schedule by ID
func (uc *ScheduleUseCase) DeleteSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return handleScheduleError(err, "schedule not found")
	}

	if err := uc.scheduleRepo.Delete(ctx, id); err != nil {
		return handleScheduleError(err, "failed to delete schedule")
	}

	uc.logger.Info("schedule deleted", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return dto.SuccessScheduleResponse(dto.ScheduleDTOFromEntity(schedule)), nil
}
