package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// ToggleSchedule enables or disables a schedule
func (uc *ScheduleUseCase) ToggleSchedule(ctx context.Context, id string, req dto.ToggleScheduleRequest) (*dto.ScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return handleScheduleError(err, "schedule not found")
	}

	if req.Enabled {
		schedule.Enable()
	} else {
		schedule.Disable()
	}

	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return handleScheduleError(err, "failed to toggle schedule")
	}

	uc.logger.Info("schedule toggled", "schedule_id", schedule.ID, "enabled", schedule.Enabled)

	return dto.SuccessScheduleResponse(dto.ScheduleDTOFromEntity(schedule)), nil
}

// EnableSchedule enables a schedule
func (uc *ScheduleUseCase) EnableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: true})
}

// DisableSchedule disables a schedule
func (uc *ScheduleUseCase) DisableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: false})
}
