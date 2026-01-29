package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// CreateSchedule creates a new schedule
func (uc *ScheduleUseCase) CreateSchedule(ctx context.Context, req dto.CreateScheduleRequest) (*dto.ScheduleResponse, error) {
	// Convert input to JSON
	inputJSON, err := dto.MapToString(req.Input)
	if err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to marshal input: %v", err),
		}, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Create new schedule
	schedule := entity.NewSchedule(req.Skill, req.CronExpression, inputJSON)

	if err := uc.scheduleRepo.Create(ctx, schedule); err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create schedule: %v", err),
		}, fmt.Errorf("failed to create schedule: %w", err)
	}

	uc.logger.Info("schedule created", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return &dto.ScheduleResponse{
		Success:  true,
		Schedule: dto.ScheduleDTOFromEntity(schedule),
	}, nil
}
