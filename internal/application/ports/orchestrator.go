package ports

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/dto"
)

// Orchestrator defines the interface for message orchestration.
// It coordinates message processing, LLM interaction, and skill execution.
type Orchestrator interface {
	// ProcessMessage processes an incoming message and returns the AI response.
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
	ProcessMessage(ctx context.Context, userID, content string, options dto.MessageOptions) (*dto.SendMessageResponse, error)

	// GetConversation retrieves the conversation history for a session.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - sessionID: Session ID
	//
	// Returns:
	//   - *dto.MessagesResponse: Response containing conversation messages
	//   - error: Error if operation failed
	GetConversation(ctx context.Context, sessionID string) (*dto.MessagesResponse, error)

	// GetUserSessions retrieves all sessions for a user.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - userID: User ID
	//
	// Returns:
	//   - *dto.SessionsResponse: Response containing user sessions
	//   - error: Error if operation failed
	GetUserSessions(ctx context.Context, userID string) (*dto.SessionsResponse, error)

	// CreateSession creates a new session for a user.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - req: CreateSessionRequest containing session details
	//
	// Returns:
	//   - *dto.SessionResponse: Response containing created session
	//   - error: Error if operation failed
	CreateSession(ctx context.Context, req dto.CreateSessionRequest) (*dto.SessionResponse, error)

	// ExecuteSkill executes a skill for a session.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - sessionID: Session ID
	//   - skillName: Name of the skill to execute
	//   - input: Input parameters for the skill
	//
	// Returns:
	//   - *dto.SkillExecutionResponse: Response containing skill execution result
	//   - error: Error if operation failed
	ExecuteSkill(ctx context.Context, sessionID, skillName string, input map[string]interface{}) (*dto.SkillExecutionResponse, error)

	// GetSessionTasks retrieves all tasks for a session.
	//
	// Parameters:
	//   - ctx: Context for the operation
	//   - sessionID: Session ID
	//
	// Returns:
	//   - *dto.TasksResponse: Response containing session tasks
	//   - error: Error if operation failed
	GetSessionTasks(ctx context.Context, sessionID string) (*dto.TasksResponse, error)
}
