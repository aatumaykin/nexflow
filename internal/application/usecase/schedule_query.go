package usecase

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// GetScheduleByID retrieves a schedule by ID
func (uc *ScheduleUseCase) GetScheduleByID(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return handleScheduleError(err, "failed to find schedule")
	}

	return dto.SuccessScheduleResponse(dto.ScheduleDTOFromEntity(schedule)), nil
}

// ListSchedules retrieves all schedules
func (uc *ScheduleUseCase) ListSchedules(ctx context.Context) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.List(ctx)
	if err != nil {
		return dto.ErrorSchedulesResponse(err), err
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return dto.SuccessSchedulesResponse(scheduleDTOs), nil
}

// ListEnabledSchedules retrieves all enabled schedules
func (uc *ScheduleUseCase) ListEnabledSchedules(ctx context.Context) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.FindEnabled(ctx)
	if err != nil {
		return dto.ErrorSchedulesResponse(err), err
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return dto.SuccessSchedulesResponse(scheduleDTOs), nil
}

// GetSchedulesBySkill retrieves all schedules for a specific skill
func (uc *ScheduleUseCase) GetSchedulesBySkill(ctx context.Context, skill string) (*dto.SchedulesResponse, error) {
	schedules, err := uc.scheduleRepo.FindBySkill(ctx, skill)
	if err != nil {
		return dto.ErrorSchedulesResponse(err), err
	}

	scheduleDTOs := make([]*dto.ScheduleDTO, 0, len(schedules))
	for _, schedule := range schedules {
		scheduleDTOs = append(scheduleDTOs, dto.ScheduleDTOFromEntity(schedule))
	}

	return dto.SuccessSchedulesResponse(scheduleDTOs), nil
}
