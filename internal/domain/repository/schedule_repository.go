package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// ScheduleRepository defines the interface for schedule data operations
type ScheduleRepository interface {
	// Create saves a new schedule
	Create(ctx context.Context, schedule *entity.Schedule) error

	// FindByID retrieves a schedule by ID
	FindByID(ctx context.Context, id string) (*entity.Schedule, error)

	// FindBySkill retrieves all schedules for a skill
	FindBySkill(ctx context.Context, skill string) ([]*entity.Schedule, error)

	// List retrieves all schedules
	List(ctx context.Context) ([]*entity.Schedule, error)

	// Update updates an existing schedule
	Update(ctx context.Context, schedule *entity.Schedule) error

	// Delete removes a schedule
	Delete(ctx context.Context, id string) error

	// FindEnabled retrieves all enabled schedules
	FindEnabled(ctx context.Context) ([]*entity.Schedule, error)
}
