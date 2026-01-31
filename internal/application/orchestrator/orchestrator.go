package orchestrator

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/application/usecase"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// Orchestrator coordinates message processing, LLM interaction, and skill execution.
// It serves as the entry point for MessageRouter and delegates to ChatUseCase for business logic.
type Orchestrator struct {
	chatUseCase *usecase.ChatUseCase
	logger      logging.Logger
}

// NewOrchestrator creates a new Orchestrator instance.
//
// Parameters:
//   - chatUseCase: ChatUseCase for message processing
//   - logger: Structured logger for logging
//
// Returns:
//   - ports.Orchestrator: Initialized orchestrator
func NewOrchestrator(chatUseCase *usecase.ChatUseCase, logger logging.Logger) ports.Orchestrator {
	return &Orchestrator{
		chatUseCase: chatUseCase,
		logger:      logger,
	}
}

// ProcessMessage processes an incoming message and returns AI response.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: User ID
//   - content: Message content
//   - options: Message options (model, max tokens, etc.)
//
// Returns:
//   - *dto.SendMessageResponse: Response containing AI message and conversation history
//   - error: Error if operation failed
func (o *Orchestrator) ProcessMessage(ctx context.Context, userID, content string, options dto.MessageOptions) (*dto.SendMessageResponse, error) {
	o.logger.Info("orchestrator: processing message", "user_id", userID, "content_length", len(content))

	// Create send message request
	req := dto.SendMessageRequest{
		UserID: userID,
		Message: dto.ChatMessage{
			Role:    "user",
			Content: content,
		},
		Options: options,
	}

	// Delegate to ChatUseCase for message processing
	resp, err := o.chatUseCase.SendMessage(ctx, req)
	if err != nil {
		o.logger.Error("orchestrator: failed to process message", "user_id", userID, "error", err)
		return nil, err
	}

	if resp.Message != nil {
		o.logger.Info("orchestrator: message processed successfully", "user_id", userID, "session_id", resp.Message.SessionID)
	} else {
		o.logger.Info("orchestrator: message processed successfully", "user_id", userID)
	}
	return resp, nil
}

// GetConversation retrieves conversation history for a session.
//
// Parameters:
//   - ctx: Context for the operation
//   - sessionID: Session ID
//
// Returns:
//   - *dto.MessagesResponse: Response containing conversation messages
//   - error: Error if operation failed
func (o *Orchestrator) GetConversation(ctx context.Context, sessionID string) (*dto.MessagesResponse, error) {
	o.logger.Info("orchestrator: getting conversation", "session_id", sessionID)

	resp, err := o.chatUseCase.GetConversation(ctx, sessionID)
	if err != nil {
		o.logger.Error("orchestrator: failed to get conversation", "session_id", sessionID, "error", err)
		return nil, err
	}

	return resp, nil
}

// GetUserSessions retrieves all sessions for a user.
//
// Parameters:
//   - ctx: Context for the operation
//   - userID: User ID
//
// Returns:
//   - *dto.SessionsResponse: Response containing user sessions
//   - error: Error if operation failed
func (o *Orchestrator) GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error) {
	o.logger.Info("orchestrator: getting user sessions", "user_id", userID)

	resp, err := o.chatUseCase.GetUserSessions(ctx, userID)
	if err != nil {
		o.logger.Error("orchestrator: failed to get user sessions", "user_id", userID, "error", err)
		return nil, err
	}

	return resp, nil
}

// CreateSession creates a new session for a user.
//
// Parameters:
//   - ctx: Context for the operation
//   - req: CreateSessionRequest containing session details
//
// Returns:
//   - *dto.SessionResponse: Response containing created session
//   - error: Error if operation failed
func (o *Orchestrator) CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error) {
	o.logger.Info("orchestrator: creating session", "user_id", req.UserID)

	resp, err := o.chatUseCase.CreateSession(ctx, req)
	if err != nil {
		o.logger.Error("orchestrator: failed to create session", "user_id", req.UserID, "error", err)
		return nil, err
	}

	return resp, nil
}

// ExecuteSkill executes a skill for a session.
//
// Parameters:
//   - ctx: Context for the operation
//   - sessionID: Session ID
//   - skillName: Name of skill to execute
//   - input: Input parameters for the skill
//
// Returns:
//   - *dto.SkillExecutionResponse: Response containing skill execution result
//   - error: Error if operation failed
func (o *Orchestrator) ExecuteSkill(ctx context.Context, sessionID, skillName string, input map[string]interface{}) (*dto.SkillExecutionResponse, error) {
	o.logger.Info("orchestrator: executing skill", "session_id", sessionID, "skill", skillName)

	resp, err := o.chatUseCase.ExecuteSkill(ctx, sessionID, skillName, input)
	if err != nil {
		o.logger.Error("orchestrator: failed to execute skill", "session_id", sessionID, "skill", skillName, "error", err)
		return nil, err
	}

	return resp, nil
}

// GetSessionTasks retrieves all tasks for a session.
//
// Parameters:
//   - ctx: Context for the operation
//   - sessionID: Session ID
//
// Returns:
//   - *dto.TasksResponse: Response containing session tasks
//   - error: Error if operation failed
func (o *Orchestrator) GetSessionTasks(ctx context.Context, sessionID string) (*dto.TasksResponse, error) {
	o.logger.Info("orchestrator: getting session tasks", "session_id", sessionID)

	resp, err := o.chatUseCase.GetSessionTasks(ctx, sessionID)
	if err != nil {
		o.logger.Error("orchestrator: failed to get session tasks", "session_id", sessionID, "error", err)
		return nil, err
	}

	return resp, nil
}
