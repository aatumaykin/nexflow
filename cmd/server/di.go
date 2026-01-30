package main

import (
	"database/sql"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/application/usecase"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	httpinf "github.com/atumaikin/nexflow/internal/infrastructure/http"
	llmadapter "github.com/atumaikin/nexflow/internal/infrastructure/llm"
	anthropic "github.com/atumaikin/nexflow/internal/infrastructure/llm/anthropic"
	llmmock "github.com/atumaikin/nexflow/internal/infrastructure/llm/mock"
	ollama "github.com/atumaikin/nexflow/internal/infrastructure/llm/ollama"
	openai "github.com/atumaikin/nexflow/internal/infrastructure/llm/openai"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database/sqlite"
	"github.com/atumaikin/nexflow/internal/infrastructure/skills"
	skillmock "github.com/atumaikin/nexflow/internal/infrastructure/skills/mock"
	"github.com/atumaikin/nexflow/internal/shared/config"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// DIContainer holds all application dependencies
type DIContainer struct {
	config  *config.Config
	logger  logging.Logger
	db      database.Database
	sqlDB   *sql.DB
	queries *database.Queries

	// Repositories
	userRepo     repository.UserRepository
	sessionRepo  repository.SessionRepository
	messageRepo  repository.MessageRepository
	taskRepo     repository.TaskRepository
	skillRepo    repository.SkillRepository
	scheduleRepo repository.ScheduleRepository

	// Ports
	llmProvider  ports.LLMProvider
	skillRuntime ports.SkillRuntime

	// Use Cases
	chatUseCase     *usecase.ChatUseCase
	userUseCase     *usecase.UserUseCase
	skillUseCase    *usecase.SkillUseCase
	scheduleUseCase *usecase.ScheduleUseCase

	// HTTP Handlers
	userHandler     *httpinf.UserHandler
	sessionHandler  *httpinf.SessionHandler
	messageHandler  *httpinf.MessageHandler
	taskHandler     *httpinf.TaskHandler
	skillHandler    *httpinf.SkillHandler
	scheduleHandler *httpinf.ScheduleHandler
	logHandler      *httpinf.LogHandler
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
		config:  cfg,
		logger:  logger,
		db:      db,
		sqlDB:   sqlDB,
		queries: database.New(sqlDB),
	}

	// Initialize repositories
	if err := container.initRepositories(); err != nil {
		return nil, err
	}

	// Initialize ports
	if err := container.initPorts(); err != nil {
		return nil, err
	}

	// Initialize use cases
	if err := container.initUseCases(); err != nil {
		return nil, err
	}

	// Initialize HTTP handlers
	if err := container.initHandlers(); err != nil {
		return nil, err
	}

	return container, nil
}

// initRepositories initializes all repository implementations
func (c *DIContainer) initRepositories() error {
	// User repository
	c.userRepo = sqlite.NewUserRepository(c.queries)

	// Session repository
	c.sessionRepo = sqlite.NewSessionRepository(c.queries)

	// Message repository
	c.messageRepo = sqlite.NewMessageRepository(c.queries)

	// Task repository
	c.taskRepo = sqlite.NewTaskRepository(c.queries)

	// Skill repository
	c.skillRepo = sqlite.NewSkillRepository(c.queries)

	// Schedule repository
	c.scheduleRepo = sqlite.NewScheduleRepository(c.queries)

	c.logger.Info("repositories initialized successfully")
	return nil
}

// initPorts initializes all port implementations
func (c *DIContainer) initPorts() error {
	// Initialize LLM provider
	if err := c.initLLMProvider(); err != nil {
		return fmt.Errorf("failed to initialize LLM provider: %w", err)
	}

	// Initialize skill runtime
	if err := c.initSkillRuntime(); err != nil {
		return fmt.Errorf("failed to initialize skill runtime: %w", err)
	}

	c.logger.Info("ports initialized successfully")
	return nil
}

// initLLMProvider initializes the LLM provider based on configuration
func (c *DIContainer) initLLMProvider() error {
	// Check if LLM config is available
	if c.config.LLM.DefaultProvider == "" {
		c.logger.Info("no LLM provider configured, using mock")
		c.llmProvider = llmmock.NewMockLLMProvider()
		return nil
	}

	// Get provider config
	providerName := c.config.LLM.DefaultProvider
	providerConfig, ok := c.config.LLM.Providers[providerName]
	if !ok {
		return fmt.Errorf("LLM provider '%s' not found in configuration", providerName)
	}

	// Get slog logger from interface
	slogLogger, ok := c.logger.(*logging.SlogLogger)
	if !ok {
		c.logger.Warn("logger is not SlogLogger, using mock LLM provider")
		c.llmProvider = llmmock.NewMockLLMProvider()
		return nil
	}

	// Create appropriate provider based on name
	var provider llmadapter.Provider
	var err error

	switch providerName {
	case "openai":
		provider, err = openai.NewProvider(&openai.Config{
			APIKey:  providerConfig.APIKey,
			BaseURL: providerConfig.BaseURL,
			Model:   providerConfig.Model,
		}, slogLogger.GetSlogLogger())
	case "anthropic":
		provider, err = anthropic.NewProvider(&anthropic.Config{
			APIKey:  providerConfig.APIKey,
			BaseURL: providerConfig.BaseURL,
			Model:   providerConfig.Model,
		}, slogLogger.GetSlogLogger())
	case "ollama":
		provider, err = ollama.NewProvider(&ollama.Config{
			BaseURL: providerConfig.BaseURL,
			Model:   providerConfig.Model,
		}, slogLogger.GetSlogLogger())
	default:
		c.logger.Warn("Unknown LLM provider, using mock", "provider", providerName)
		c.llmProvider = llmmock.NewMockLLMProvider()
		return nil
	}

	if err != nil {
		return fmt.Errorf("failed to create LLM provider: %w", err)
	}

	// Wrap provider with adapter
	c.llmProvider = llmadapter.NewProviderAdapter(provider)

	c.logger.Info("LLM provider initialized",
		"provider", providerName,
		"model", providerConfig.Model)

	return nil
}

// initSkillRuntime initializes the skill runtime based on configuration
func (c *DIContainer) initSkillRuntime() error {
	// Check if skills config is available
	if c.config.Skills.Directory == "" {
		c.logger.Info("no skills directory configured, using mock")
		c.skillRuntime = skillmock.NewMockSkillRuntime()
		return nil
	}

	// Get slog logger from interface
	slogLogger, ok := c.logger.(*logging.SlogLogger)
	if !ok {
		c.logger.Warn("logger is not SlogLogger, using mock skill runtime")
		c.skillRuntime = skillmock.NewMockSkillRuntime()
		return nil
	}

	// Create local skill runtime
	localRuntime, err := skills.NewLocalRuntime(&skills.Config{
		Directory:      c.config.Skills.Directory,
		TimeoutSeconds: c.config.Skills.TimeoutSec,
		SandboxEnabled: c.config.Skills.SandboxEnabled,
	}, slogLogger.GetSlogLogger())
	if err != nil {
		return fmt.Errorf("failed to create local skill runtime: %w", err)
	}

	// Wrap runtime with adapter
	c.skillRuntime = skills.NewRuntimeAdapter(localRuntime)

	c.logger.Info("skill runtime initialized",
		"directory", c.config.Skills.Directory,
		"timeout", c.config.Skills.TimeoutSec,
		"sandbox", c.config.Skills.SandboxEnabled)

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

// initHandlers initializes all HTTP handlers
func (c *DIContainer) initHandlers() error {
	// User handler
	c.userHandler = httpinf.NewUserHandler(c.userUseCase, c.logger)

	// Session handler
	c.sessionHandler = httpinf.NewSessionHandler(c.chatUseCase, c.logger)

	// Message handler
	c.messageHandler = httpinf.NewMessageHandler(c.chatUseCase, c.logger)

	// Task handler
	c.taskHandler = httpinf.NewTaskHandler(c.chatUseCase, c.logger)

	// Skill handler
	c.skillHandler = httpinf.NewSkillHandler(c.skillUseCase, c.logger)

	// Schedule handler
	c.scheduleHandler = httpinf.NewScheduleHandler(c.scheduleUseCase, c.logger)

	// Log handler
	c.logHandler = httpinf.NewLogHandler(c.logger)

	c.logger.Info("HTTP handlers initialized successfully")
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

// Getters for HTTP handlers
func (c *DIContainer) UserHandler() *httpinf.UserHandler {
	return c.userHandler
}

func (c *DIContainer) SessionHandler() *httpinf.SessionHandler {
	return c.sessionHandler
}

func (c *DIContainer) MessageHandler() *httpinf.MessageHandler {
	return c.messageHandler
}

func (c *DIContainer) TaskHandler() *httpinf.TaskHandler {
	return c.taskHandler
}

func (c *DIContainer) SkillHandler() *httpinf.SkillHandler {
	return c.skillHandler
}

func (c *DIContainer) ScheduleHandler() *httpinf.ScheduleHandler {
	return c.scheduleHandler
}

func (c *DIContainer) LogHandler() *httpinf.LogHandler {
	return c.logHandler
}

// Shutdown performs cleanup operations
func (c *DIContainer) Shutdown() error {
	c.logger.Info("shutting down DI container")
	// Database is closed in main
	return nil
}
