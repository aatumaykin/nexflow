package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// ChatUseCase handles chat-related business logic
type ChatUseCase struct {
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	messageRepo  repository.MessageRepository
	taskRepo     repository.TaskRepository
	llmProvider  ports.LLMProvider
	skillRuntime ports.SkillRuntime
	logger       logging.Logger
}

// NewChatUseCase creates a new ChatUseCase
func NewChatUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	messageRepo repository.MessageRepository,
	taskRepo repository.TaskRepository,
	llmProvider ports.LLMProvider,
	skillRuntime ports.SkillRuntime,
	logger logging.Logger,
) *ChatUseCase {
	return &ChatUseCase{
		userRepo:     userRepo,
		sessionRepo:  sessionRepo,
		messageRepo:  messageRepo,
		taskRepo:     taskRepo,
		llmProvider:  llmProvider,
		skillRuntime: skillRuntime,
		logger:       logger,
	}
}

// SendMessage processes a user message and returns AI response
func (uc *ChatUseCase) SendMessage(ctx context.Context, req dto.SendMessageRequest) (*dto.SendMessageResponse, error) {
	// Find or create user
	user, err := uc.userRepo.FindByChannel(ctx, "web", req.UserID) // Default to web channel
	if err != nil {
		// User not found, create new user
		newUser := entity.NewUser("web", req.UserID)
		if err := uc.userRepo.Create(ctx, newUser); err != nil {
			return &dto.SendMessageResponse{
				Success: false,
				Error:   fmt.Sprintf("failed to create user: %v", err),
			}, fmt.Errorf("failed to create user: %w", err)
		}
		user = newUser
	}

	// Create new session (simplified - in real app, we'd manage active sessions)
	session := entity.NewSession(string(user.ID))
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create session: %v", err),
		}, fmt.Errorf("failed to create session: %w", err)
	}

	// Save user message
	userMessage := entity.NewUserMessage(string(session.ID), req.Message.Content)
	if err := uc.messageRepo.Create(ctx, userMessage); err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to save user message: %v", err),
		}, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history
	messages, err := uc.messageRepo.FindBySessionID(ctx, string(session.ID))
	if err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get conversation history: %v", err),
		}, fmt.Errorf("failed to get conversation history: %w", err)
	}

	// Convert to LLM format
	llmMessages := make([]ports.Message, 0, len(messages))
	for _, msg := range messages {
		llmMessages = append(llmMessages, ports.Message{
			Role:    string(msg.Role),
			Content: msg.Content,
		})
	}

	// Call LLM
	llmReq := ports.CompletionRequest{
		Messages:  llmMessages,
		Model:     req.Options.Model,
		MaxTokens: req.Options.MaxTokens,
	}
	llmResp, err := uc.llmProvider.Generate(ctx, llmReq)
	if err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to generate response: %v", err),
		}, fmt.Errorf("failed to generate response: %w", err)
	}

	// Save assistant message
	assistantMessage := entity.NewAssistantMessage(string(session.ID), llmResp.Message.Content)
	if err := uc.messageRepo.Create(ctx, assistantMessage); err != nil {
		uc.logger.Error("failed to save assistant message", "error", err)
	}

	// Update session timestamp
	session.UpdateTimestamp()
	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		uc.logger.Error("failed to update session", "error", err)
	}

	// Get updated messages for response
	updatedMessages, err := uc.messageRepo.FindBySessionID(ctx, string(session.ID))
	if err != nil {
		uc.logger.Error("failed to get updated messages", "error", err)
	}

	// Convert messages to DTOs
	messageDTOs := make([]*dto.MessageDTO, 0, len(updatedMessages))
	for _, msg := range updatedMessages {
		messageDTOs = append(messageDTOs, dto.MessageDTOFromEntity(msg))
	}

	return &dto.SendMessageResponse{
		Success:  true,
		Message:  dto.MessageDTOFromEntity(assistantMessage),
		Messages: messageDTOs,
	}, nil
}

// GetConversation retrieves conversation history for a session
func (uc *ChatUseCase) GetConversation(ctx context.Context, sessionID string) (*dto.MessagesResponse, error) {
	messages, err := uc.messageRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return dto.ErrorMessageResponse(fmt.Errorf("failed to get conversation: %w", err)), fmt.Errorf("failed to get conversation: %w", err)
	}

	messageDTOs := make([]*dto.MessageDTO, 0, len(messages))
	for _, msg := range messages {
		messageDTOs = append(messageDTOs, dto.MessageDTOFromEntity(msg))
	}

	return dto.SuccessMessagesResponse(messageDTOs), nil
}

// GetUserSessions retrieves all sessions for a user
func (uc *ChatUseCase) GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return dto.ErrorSessionsResponse(fmt.Errorf("failed to get user sessions: %w", err)), fmt.Errorf("failed to get user sessions: %w", err)
	}

	sessionDTOs := make([]*dto.SessionDTO, 0, len(sessions))
	for _, session := range sessions {
		sessionDTOs = append(sessionDTOs, dto.SessionDTOFromEntity(session))
	}

	return dto.SuccessSessionsResponse(sessionDTOs), nil
}

// CreateSession creates a new session for a user
func (uc *ChatUseCase) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	session := entity.NewSession(req.UserID)
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return dto.ErrorSessionResponse(fmt.Errorf("failed to create session: %w", err)), fmt.Errorf("failed to create session: %w", err)
	}

	return dto.SuccessSessionResponse(dto.SessionDTOFromEntity(session)), nil
}

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
	task := entity.NewTask(string(sessionID), skillName, string(inputJSON))
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
