package usecase

import (
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// handleSendError handles errors in SendMessage use case
func handleSendError(err error, message string) (*dto.SendMessageResponse, error) {
	return dto.ErrorSendMessageResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleUserError handles errors in User use case
func handleUserError(err error, message string) (*dto.UserResponse, error) {
	return dto.ErrorUserResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleSessionError handles errors in Session use case
func handleSessionError(err error, message string) (*dto.SessionResponse, error) {
	return dto.ErrorSessionResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleSkillError handles errors in Skill use case
func handleSkillError(err error, message string) (*dto.SkillResponse, error) {
	return dto.ErrorSkillResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleScheduleError handles errors in Schedule use case
func handleScheduleError(err error, message string) (*dto.ScheduleResponse, error) {
	return dto.ErrorScheduleResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleMessagesError handles errors in Messages use case
func handleMessagesError(err error, message string) (*dto.MessagesResponse, error) {
	return dto.ErrorMessageResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleTasksError handles errors in Tasks use case
func handleTasksError(err error, message string) (*dto.TasksResponse, error) {
	return dto.ErrorTaskResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}

// handleSkillExecutionError handles errors in SkillExecution use case
func handleSkillExecutionError(err error, message string) (*dto.SkillExecutionResponse, error) {
	return dto.ErrorSkillExecutionResponse(fmt.Errorf("%s: %w", message, err)), fmt.Errorf("%s: %w", message, err)
}
