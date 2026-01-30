package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// ExecuteSkill executes a skill based on LLM response
func (uc *ChatUseCase) ExecuteSkill(ctx context.Context, sessionID, skillName string, input map[string]interface{}) (*dto.SkillExecutionResponse, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return handleSkillExecutionError(err, "failed to marshal skill input")
	}

	task := entity.NewTask(sessionID, skillName, string(inputJSON))
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return handleSkillExecutionError(err, "failed to create task")
	}

	execution, err := uc.skillRuntime.Execute(ctx, skillName, input)
	if err != nil {
		task.SetFailed(fmt.Sprintf("skill execution failed: %v", err))
		if err := uc.taskRepo.Update(ctx, task); err != nil {
			uc.logger.Error("failed to update task status", "error", err)
		}
		return handleSkillExecutionError(err, "skill execution failed")
	}

	task.SetRunning()
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Error("failed to update task status", "error", err)
	}

	if execution.Success {
		task.SetCompleted(execution.Output)
	} else {
		task.SetFailed(execution.Error)
	}

	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Error("failed to update task completion", "error", err)
	}

	return &dto.SkillExecutionResponse{
		Success: execution.Success,
		Output:  execution.Output,
		Error:   execution.Error,
	}, nil
}

// GetSessionTasks retrieves all tasks for a session
func (uc *ChatUseCase) GetSessionTasks(ctx context.Context, sessionID string) (*dto.TasksResponse, error) {
	tasks, err := uc.taskRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return handleTasksError(err, "failed to get session tasks")
	}

	taskDTOs := make([]*dto.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, dto.TaskDTOFromEntity(task))
	}

	return dto.SuccessTasksResponse(taskDTOs), nil
}
