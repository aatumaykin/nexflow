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
	// Convert input to JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to marshal skill input: %v", err),
		}, fmt.Errorf("failed to marshal skill input: %w", err)
	}

	// Create task
	task := entity.NewTask(sessionID, skillName, string(inputJSON))
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create task: %v", err),
		}, fmt.Errorf("failed to create task: %w", err)
	}

	// Execute skill using SkillRuntime port
	execution, err := uc.skillRuntime.Execute(ctx, skillName, input)
	if err != nil {
		task.SetFailed(fmt.Sprintf("skill execution failed: %v", err))
		if err := uc.taskRepo.Update(ctx, task); err != nil {
			uc.logger.Error("failed to update task status", "error", err)
		}
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   execution.Error,
		}, fmt.Errorf("skill execution failed: %w", err)
	}

	// Update task with execution result
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
		return dto.ErrorTaskResponse(fmt.Errorf("failed to get session tasks: %w", err)), fmt.Errorf("failed to get session tasks: %w", err)
	}

	taskDTOs := make([]*dto.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, dto.TaskDTOFromEntity(task))
	}

	return dto.SuccessTasksResponse(taskDTOs), nil
}
