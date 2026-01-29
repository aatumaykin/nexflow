package main

import (
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/application/usecase"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	llmmock "github.com/atumaikin/nexflow/internal/infrastructure/llm/mock"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/sqlite"
	skillmock "github.com/atumaikin/nexflow/internal/infrastructure/skills/mock"
	"github.com/atumaikin/nexflow/internal/shared/config"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// DIContainer holds all application dependencies
type DIContainer struct {
	config *config.Config
	logger logging.Logger
	db     database.Database
	sqlDB  *sql.DB

	// Repositories
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	messageRepo  repository.MessageRepository
	taskRepo     repository.TaskRepository
	skillRepo    repository.SkillRepository
	scheduleRepo repository.ScheduleRepository

	// Ports
	llmProvider  ports.LLMProvider  // TODO: implement
	skillRuntime ports.SkillRuntime // TODO: implement

	// Use Cases
	chatUseCase     *usecase.ChatUseCase
	userUseCase     *usecase.UserUseCase
	skillUseCase    *usecase.SkillUseCase
	scheduleUseCase *usecase.ScheduleUseCase
}

// NewDIContainer creates and initializes the DI container
func NewDIContainer(cfg *config.Config, logger logging.Logger, db database.Database) (*DIContainer, error) {
	// Get underlying SQL DB for repository implementations
	// Type assertion to get *sql.DB from database.Database
	dbImpl, ok := db.(*database.DB)
	if !ok {
		return nil, fmt.Errorf("failed to assert database.Database to *database.DB")
	}
	sqlDB := dbImpl.GetDB()

	container := &DIContainer{
		config: cfg,
		logger: logger,
		db:     db,
		sqlDB:  sqlDB,
	}

	// Initialize repositories
	if err := container.initRepositories(); err != nil {
		return nil, err
	}

	// Initialize ports (TODO: implement actual implementations)
	if err := container.initPorts(); err != nil {
		return nil, err
	}

	// Initialize use cases
	if err := container.initUseCases(); err != nil {
		return nil, err
	}

	return container, nil
}

// initRepositories initializes all repository implementations
func (c *DIContainer) initRepositories() error {
	// User repository
	c.userRepo = sqlite.NewUserRepository(c.sqlDB)

	// Session repository
	c.sessionRepo = sqlite.NewSessionRepository(c.sqlDB)

	// Message repository
	c.messageRepo = sqlite.NewMessageRepository(c.sqlDB)

	// Task repository
	c.taskRepo = sqlite.NewTaskRepository(c.sqlDB)

	// Skill repository
	c.skillRepo = sqlite.NewSkillRepository(c.sqlDB)

	// Schedule repository
	c.scheduleRepo = sqlite.NewScheduleRepository(c.sqlDB)

	c.logger.Info("repositories initialized successfully")
	return nil
}

// initPorts initializes all port implementations
func (c *DIContainer) initPorts() error {
	// Initialize LLM provider mock
	c.llmProvider = llmmock.NewMockLLMProvider()

	// Initialize skill runtime mock
	c.skillRuntime = skillmock.NewMockSkillRuntime()

	c.logger.Info("ports initialized successfully (mock implementations)")
	return nil
}

// initUseCases initializes all use cases
func (c *DIContainer) initUseCases() error {
	// Chat use case
	c.chatUseCase = usecase.NewChatUseCase(
		c.userRepo,
		c.sessionRepo,
		c.messageRepo,
		c.taskRepo,
		c.llmProvider,
		c.skillRuntime,
		c.logger,
	)

	// User use case
	c.userUseCase = usecase.NewUserUseCase(
		c.userRepo,
		c.logger,
	)

	// Skill use case
	c.skillUseCase = usecase.NewSkillUseCase(
		c.skillRepo,
		c.skillRuntime,
		c.logger,
	)

	// Schedule use case
	c.scheduleUseCase = usecase.NewScheduleUseCase(
		c.scheduleRepo,
		c.logger,
	)

	c.logger.Info("use cases initialized successfully")
	return nil
}

// Getters for use cases
func (c *DIContainer) ChatUseCase() *usecase.ChatUseCase {
	return c.chatUseCase
}

func (c *DIContainer) UserUseCase() *usecase.UserUseCase {
	return c.userUseCase
}

func (c *DIContainer) SkillUseCase() *usecase.SkillUseCase {
	return c.skillUseCase
}

func (c *DIContainer) ScheduleUseCase() *usecase.ScheduleUseCase {
	return c.scheduleUseCase
}

// Shutdown performs cleanup operations
func (c *DIContainer) Shutdown() error {
	c.logger.Info("shutting down DI container")
	// Database is closed in main
	return nil
}
