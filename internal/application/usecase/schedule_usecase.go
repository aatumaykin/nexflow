package usecase

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// ScheduleUseCase handles schedule-related business logic
type ScheduleUseCase struct {
	scheduleRepo repository.ScheduleRepository
	logger       logging.Logger
}

// NewScheduleUseCase creates a new ScheduleUseCase
func NewScheduleUseCase(
	scheduleRepo repository.ScheduleRepository,
	logger logging.Logger,
) *ScheduleUseCase {
	return &ScheduleUseCase{
		scheduleRepo: scheduleRepo,
		logger:       logger,
	}
}

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

// UpdateSchedule updates an existing schedule
func (uc *ScheduleUseCase) UpdateSchedule(ctx context.Context, id string, req dto.UpdateScheduleRequest) (*dto.ScheduleResponse, error) {
	// Get existing schedule
	schedule, err := uc.scheduleRepo.FindByID(ctx, id)
	if err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("schedule not found: %v", err),
		}, fmt.Errorf("schedule not found: %w", err)
	}

	// Update fields
	if req.CronExpression != "" {
		schedule.CronExpression = valueobject.MustNewCronExpression(req.CronExpression)
	}
	if req.Input != nil {
		inputJSON, err := dto.MapToString(req.Input)
		if err != nil {
			return &dto.ScheduleResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to marshal input: %v", err),
			}, fmt.Errorf("failed to marshal input: %w", err)
		}
		schedule.Input = inputJSON
	}
	if req.Enabled != nil {
		if *req.Enabled {
			schedule.Enable()
		} else {
			schedule.Disable()
		}
	}

	// Save updated schedule
	if err := uc.scheduleRepo.Update(ctx, schedule); err != nil {
		return &dto.ScheduleResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to update schedule: %v", err),
		}, fmt.Errorf("failed to update schedule: %w", err)
	}

	uc.logger.Info("schedule updated", "schedule_id", schedule.ID, "skill", schedule.Skill)

	return &dto.ScheduleResponse{
		Success:  true,
		Schedule: dto.ScheduleDTOFromEntity(schedule),
	}, nil
}

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

// EnableSchedule enables a schedule
func (uc *ScheduleUseCase) EnableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: true})
}

// DisableSchedule disables a schedule
func (uc *ScheduleUseCase) DisableSchedule(ctx context.Context, id string) (*dto.ScheduleResponse, error) {
	return uc.ToggleSchedule(ctx, id, dto.ToggleScheduleRequest{Enabled: false})
}
