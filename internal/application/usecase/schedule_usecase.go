package usecase

import (
	"github.com/atumaikin/nexflow/internal/domain/repository"
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
