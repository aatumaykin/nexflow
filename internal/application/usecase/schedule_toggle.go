package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// ToggleSchedule enables or disables a schedule
func (uc *ScheduleUseCase) ToggleSchedule(ctx context.Context, id string, req dto.ToggleScheduleRequest) (*dto.ScheduleResponse, error) {
	// Get existing schedule
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("schedule not found: %v", err),
		}, fmt.Errorf("schedule not found: %w", err)
	}

	// Toggle enabled status
	if req.Enabled {
		schedule.Enable()
	} else {
		schedule.Disable()
	}

	// Save updated schedule
	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to toggle schedule: %v", err),
		}, fmt.Errorf("failed to toggle schedule: %w", err)
	}

	uc.logger.Info("schedule toggled", "schedule_id", schedule.ID, "enabled", schedule.Enabled)

	return &dto.ScheduleResponse{
		Success:  true,
		Schedule: dto.ScheduleDTOFromEntity(schedule),
	}, nil
}

// EnableSchedule enables a schedule
func (uc *ScheduleUseCase) EnableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: true})
}

// DisableSchedule disables a schedule
func (uc *ScheduleUseCase) DisableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: false})
}
