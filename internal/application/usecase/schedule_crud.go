package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// DeleteSchedule deletes a schedule by ID
func (uc *ScheduleUseCase) DeleteSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	// Check if schedule exists
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("schedule not found: %v", err),
		}, fmt.Errorf("schedule not found: %w", err)
	}

	// Delete schedule
	if err := uc.scheduleRepo.Delete(ctx, id); err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to delete schedule: %v", err),
		}, fmt.Errorf("failed to delete schedule: %w", err)
	}

	uc.logger.Info("schedule deleted", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return &dto.ScheduleResponse{
		Success:  true,
		Schedule: dto.ScheduleDTOFromEntity(schedule),
	}, nil
}
