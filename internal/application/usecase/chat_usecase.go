package usecase

import (
	"github.com/atumaikin/nexflow/internal/application/ports"
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

// NewChatUseCase creates a new ChatUseCase with all required dependencies
//
// Parameters:
//   - userRepo: Repository for user data access
//   - sessionRepo: Repository for session data access
//   - messageRepo: Repository for message data access
//   - taskRepo: Repository for task data access
//   - llmProvider: LLM provider for generating responses
//   - skillRuntime: Skill runtime for executing skills
//   - logger: Structured logger for logging
//
// Returns:
//   - *ChatUseCase: Initialized chat use case
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
