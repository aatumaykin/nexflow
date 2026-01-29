package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetScheduleByID retrieves a schedule by ID
func (uc *ScheduleUseCase) GetScheduleByID(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to find schedule: %v", err),
		}, fmt.Errorf("failed to find schedule: %w", err)
	}

	return &dto.ScheduleResponse{
		Success:  true,
		Schedule: dto.ScheduleDTOFromEntity(schedule),
	}, nil
}

// ListSchedules retrieves all schedules
func (uc *ScheduleUseCase) ListSchedules(ctx context.Context) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.List(ctx)
	if err != nil {
		return &dto.SchedulesResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to list schedules: %v", err),
		}, fmt.Errorf("failed to list schedules: %w", err)
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return &dto.SchedulesResponse{
		Success:   true,
		Schedules: scheduleDTOs,
	}, nil
}

// ListEnabledSchedules retrieves all enabled schedules
func (uc *ScheduleUseCase) ListEnabledSchedules(ctx context.Context) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.FindEnabled(ctx)
	if err != nil {
		return &dto.SchedulesResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to list enabled schedules: %v", err),
		}, fmt.Errorf("failed to list enabled schedules: %w", err)
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return &dto.SchedulesResponse{
		Success:   true,
		Schedules: scheduleDTOs,
	}, nil
}

// GetSchedulesBySkill retrieves all schedules for a specific skill
func (uc *ScheduleUseCase) GetSchedulesBySkill(ctx context.Context, skill string) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.FindBySkill(ctx, skill)
	if err != nil {
		return &dto.SchedulesResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to find schedules by skill: %v", err),
		}, fmt.Errorf("failed to find schedules by skill: %w", err)
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return &dto.SchedulesResponse{
		Success:   true,
		Schedules: scheduleDTOs,
	}, nil
}
