package mappers

import (
	"database/sql"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/utils"
)

// TaskToDomain converts SQLC Task model to domain Task entity.
func TaskToDomain(dbTask *dbmodel.Task) *entity.Task {
	if dbTask == nil {
		return nil
	}

	output := ""
	if dbTask.Output.Valid {
		output = dbTask.Output.String
	}

	taskErr := ""
	if dbTask.Error.Valid {
		taskErr = dbTask.Error.String
	}

	return &entity.Task{
		ID:        valueobject.TaskID(dbTask.ID),
		SessionID: valueobject.MustNewSessionID(dbTask.SessionID),
		Skill:     dbTask.Skill,
		Input:     dbTask.Input,
		Output:    output,
		Status:    valueobject.MustNewTaskStatus(dbTask.Status),
		Error:     taskErr,
		CreatedAt: utils.ParseTimeRFC3339(dbTask.CreatedAt),
		UpdatedAt: utils.ParseTimeRFC3339(dbTask.UpdatedAt),
	}
}

// TaskToDB converts domain Task entity to SQLC Task model.
func TaskToDB(task *entity.Task) *dbmodel.Task {
	if task == nil {
		return nil
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

	return &dbmodel.Task{
		ID:        string(task.ID),
		SessionID: string(task.SessionID),
		Skill:     task.Skill,
		Input:     task.Input,
		Output:    output,
		Status:    string(task.Status),
		Error:     taskErr,
		CreatedAt: utils.FormatTimeRFC3339(task.CreatedAt),
		UpdatedAt: utils.FormatTimeRFC3339(task.UpdatedAt),
	}
}

// TasksToDomain converts slice of SQLC Task models to domain Task entities.
func TasksToDomain(dbTasks []dbmodel.Task) []*entity.Task {
	tasks := make([]*entity.Task, 0, len(dbTasks))
	for i := range dbTasks {
		tasks = append(tasks, TaskToDomain(&dbTasks[i]))
	}
	return tasks
}
