package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.TaskRepository = (*TaskRepository)(nil)

type TaskRepository struct {
	queries *database.Queries
}

func NewTaskRepository(queries *database.Queries) *TaskRepository {
	return &TaskRepository{queries: queries}
}

func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	dbTask := mappers.TaskToDB(task)
	if dbTask == nil {
		return fmt.Errorf("failed to convert task to db model")
	}

	_, err := r.queries.CreateTask(ctx, database.CreateTaskParams{
		ID:        dbTask.ID,
		SessionID: dbTask.SessionID,
		Skill:     dbTask.Skill,
		Input:     dbTask.Input,
		Status:    dbTask.Status,
		CreatedAt: dbTask.CreatedAt,
		UpdatedAt: dbTask.UpdatedAt,
	})

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	dbTask, err := r.queries.GetTaskByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	return mappers.TaskToDomain(&dbTask), nil
}

func (r *TaskRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Task, error) {
	dbTasks, err := r.queries.GetTasksBySessionID(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks by session id: %w", err)
	}

	return mappers.TasksToDomain(dbTasks), nil
}

func (r *TaskRepository) Update(ctx context.Context, task *entity.Task) error {
	dbTask := mappers.TaskToDB(task)
	if dbTask == nil {
		return fmt.Errorf("failed to convert task to db model")
	}

	var output sql.NullString
	if task.Output != "" {
		output.Valid = true
		output.String = task.Output
	}

	var taskErr sql.NullString
	if task.Error != "" {
		taskErr.Valid = true
		taskErr.String = task.Error
	}

	_, err := r.queries.UpdateTask(ctx, database.UpdateTaskParams{
		Output:    output,
		Status:    dbTask.Status,
		Error:     taskErr,
		UpdatedAt: time.Now().Format(time.RFC3339),
		ID:        dbTask.ID,
	})

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	_, err := r.queries.GetTaskByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("task not found: %s", id)
		}
		return fmt.Errorf("failed to check task existence: %w", err)
	}

	err = r.queries.DeleteTask(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
