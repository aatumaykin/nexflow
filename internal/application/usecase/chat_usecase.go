package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// ChatUseCase handles chat-related business logic
type ChatUseCase struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	messageRepo repository.MessageRepository
	taskRepo    repository.TaskRepository
	llmProvider ports.LLMProvider
	logger      logging.Logger
}

// NewChatUseCase creates a new ChatUseCase
func NewChatUseCase(
	userRepo repository.UserRepository,
	sessionRepo repository.SessionRepository,
	messageRepo repository.MessageRepository,
	taskRepo repository.TaskRepository,
	llmProvider ports.LLMProvider,
	logger logging.Logger,
) *ChatUseCase {
	return &ChatUseCase{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		messageRepo: messageRepo,
		taskRepo:    taskRepo,
		llmProvider: llmProvider,
		logger:      logger,
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
	session := entity.NewSession(user.ID)
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create session: %v", err),
		}, fmt.Errorf("failed to create session: %w", err)
	}

	// Save user message
	userMessage := entity.NewUserMessage(session.ID, req.Message.Content)
	if err := uc.messageRepo.Create(ctx, userMessage); err != nil {
		return &dto.SendMessageResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to save user message: %v", err),
		}, fmt.Errorf("failed to save user message: %w", err)
	}

	// Get conversation history
	messages, err := uc.messageRepo.FindBySessionID(ctx, session.ID)
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
			Role:    msg.Role,
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
	assistantMessage := entity.NewAssistantMessage(session.ID, llmResp.Message.Content)
	if err := uc.messageRepo.Create(ctx, assistantMessage); err != nil {
		uc.logger.Error("failed to save assistant message", "error", err)
	}

	// Update session timestamp
	session.UpdateTimestamp()
	if err := uc.sessionRepo.Update(ctx, session); err != nil {
		uc.logger.Error("failed to update session", "error", err)
	}

	// Get updated messages for response
	updatedMessages, err := uc.messageRepo.FindBySessionID(ctx, session.ID)
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
		return &dto.MessagesResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get conversation: %v", err),
		}, fmt.Errorf("failed to get conversation: %w", err)
	}

	messageDTOs := make([]*dto.MessageDTO, 0, len(messages))
	for _, msg := range messages {
		messageDTOs = append(messageDTOs, dto.MessageDTOFromEntity(msg))
	}

	return &dto.MessagesResponse{
		Success:  true,
		Messages: messageDTOs,
	}, nil
}

// GetUserSessions retrieves all sessions for a user
func (uc *ChatUseCase) GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error) {
	sessions, err := uc.sessionRepo.FindByUserID(ctx, userID)
	if err != nil {
		return &dto.SessionsResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get user sessions: %v", err),
		}, fmt.Errorf("failed to get user sessions: %w", err)
	}

	sessionDTOs := make([]*dto.SessionDTO, 0, len(sessions))
	for _, session := range sessions {
		sessionDTOs = append(sessionDTOs, dto.SessionDTOFromEntity(session))
	}

	return &dto.SessionsResponse{
		Success:  true,
		Sessions: sessionDTOs,
	}, nil
}

// CreateSession creates a new session for a user
func (uc *ChatUseCase) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	session := entity.NewSession(req.UserID)
	if err := uc.sessionRepo.Create(ctx, session); err != nil {
		return &dto.SessionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create session: %v", err),
		}, fmt.Errorf("failed to create session: %w", err)
	}

	return &dto.SessionResponse{
		Success: true,
		Session: dto.SessionDTOFromEntity(session),
	}, nil
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
	task := entity.NewTask(sessionID, skillName, string(inputJSON))
	if err := uc.taskRepo.Create(ctx, task); err != nil {
		return &dto.SkillExecutionResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to create task: %v", err),
		}, fmt.Errorf("failed to create task: %w", err)
	}

	// Execute skill (simplified - in real app, use SkillRuntime port)
	task.SetRunning()
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Error("failed to update task status", "error", err)
	}

	// Simulate skill execution
	time.Sleep(100 * time.Millisecond)

	// Set task as completed
	output := fmt.Sprintf(`{"result": "skill %s executed successfully", "input": %s}`, skillName, string(inputJSON))
	task.SetCompleted(output)
	if err := uc.taskRepo.Update(ctx, task); err != nil {
		uc.logger.Error("failed to update task completion", "error", err)
	}

	return &dto.SkillExecutionResponse{
		Success: true,
		Output:  output,
	}, nil
}

// GetSessionTasks retrieves all tasks for a session
func (uc *ChatUseCase) GetSessionTasks(ctx context.Context, sessionID string) (*dto.TasksResponse, error) {
	tasks, err := uc.taskRepo.FindBySessionID(ctx, sessionID)
	if err != nil {
		return &dto.TasksResponse{
			Success: false,
			Error:   fmt.Sprintf("failed to get session tasks: %v", err),
		}, fmt.Errorf("failed to get session tasks: %w", err)
	}

	taskDTOs := make([]*dto.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, dto.TaskDTOFromEntity(task))
	}

	return &dto.TasksResponse{
		Success: true,
		Tasks:   taskDTOs,
	}, nil
}
