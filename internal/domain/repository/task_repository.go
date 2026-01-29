package repository

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	// Create saves a new task
	Create(ctx context.Context, task *entity.Task) error

	// FindByID retrieves a task by ID
	FindByID(ctx context.Context, id string) (*entity.Task, error)

	// FindBySessionID retrieves all tasks for a session
	FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Task, error)

	// Update updates an existing task
	Update(ctx context.Context, task *entity.Task) error

	// Delete removes a task
	Delete(ctx context.Context, id string) error
}
