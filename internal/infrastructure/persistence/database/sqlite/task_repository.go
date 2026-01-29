package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/mappers"
)

var _ repository.TaskRepository = (*TaskRepository)(nil)

// TaskRepository implements repository.TaskRepository using SQLC-generated queries
type TaskRepository struct {
	db *sql.DB
}

// NewTaskRepository creates a new TaskRepository instance
func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Create saves a new task
func (r *TaskRepository) Create(ctx context.Context, task *entity.Task) error {
	dbTask := mappers.TaskToDB(task)
	if dbTask == nil {
		return fmt.Errorf("failed to convert task to db model")
	}

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO tasks (id, session_id, skill, input, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		dbTask.ID, dbTask.SessionID, dbTask.Skill, dbTask.Input, dbTask.Status, dbTask.CreatedAt, dbTask.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// FindByID retrieves a task by ID
func (r *TaskRepository) FindByID(ctx context.Context, id string) (*entity.Task, error) {
	var dbTask dbmodel.Task

	err := r.db.QueryRowContext(ctx,
		`SELECT id, session_id, skill, input, output, status, error, created_at, updated_at FROM tasks WHERE id = ? LIMIT 1`,
		id,
	).Scan(
		&dbTask.ID, &dbTask.SessionID, &dbTask.Skill, &dbTask.Input,
		&dbTask.Output, &dbTask.Status, &dbTask.Error,
		&dbTask.CreatedAt, &dbTask.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found: %s", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	return mappers.TaskToDomain(&dbTask), nil
}

// FindBySessionID retrieves all tasks for a session
func (r *TaskRepository) FindBySessionID(ctx context.Context, sessionID string) ([]*entity.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, session_id, skill, input, output, status, error, created_at, updated_at FROM tasks WHERE session_id = ? ORDER BY created_at DESC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to find tasks by session id: %w", err)
	}
	defer rows.Close()

	var dbTasks []dbmodel.Task
	for rows.Next() {
		var dbTask dbmodel.Task

		if err := rows.Scan(
			&dbTask.ID, &dbTask.SessionID, &dbTask.Skill, &dbTask.Input,
			&dbTask.Output, &dbTask.Status, &dbTask.Error,
			&dbTask.CreatedAt, &dbTask.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}

		dbTasks = append(dbTasks, dbTask)
	}

	return mappers.TasksToDomain(dbTasks), nil
}

// Update updates an existing task
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

	now := time.Now().Format(time.RFC3339)
	_, err := r.db.ExecContext(ctx,
		`UPDATE tasks SET output = ?, status = ?, error = ?, updated_at = ? WHERE id = ?`,
		output, task.Status, taskErr, now, dbTask.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	return nil
}

// Delete removes a task
func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task not found: %s", id)
	}

	return nil
}
