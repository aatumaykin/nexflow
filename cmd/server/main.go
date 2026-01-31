package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpinf "github.com/atumaikin/nexflow/internal/infrastructure/http"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/config"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

func main() {
	// Load configuration
	cfg, err := config.Load("config.yml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger, err := logging.New(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	logger.Info("Starting Nexflow server",
		"version", "0.1.0",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
	)

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Migrations completed successfully")

	// Initialize DI container
	diContainer, err := NewDIContainer(cfg, logger, db)
	if err != nil {
		logger.Error("Failed to initialize DI container", "error", err)
		os.Exit(1)
	}
	logger.Info("DI container initialized successfully")

	// Start message router to begin receiving messages from connectors
	if err := diContainer.MessageRouter().Start(); err != nil {
		logger.Error("Failed to start message router", "error", err)
		os.Exit(1)
	}
	logger.Info("Message router started successfully")

	// Access use cases from DI container
	// chatUseCase := diContainer.ChatUseCase()
	// userUseCase := diContainer.UserUseCase()
	// skillUseCase := diContainer.SkillUseCase()
	// scheduleUseCase := diContainer.ScheduleUseCase()

	// Initialize HTTP server
	router := httpinf.NewRouter()

	// Register routes
	httpinf.RegisterUserRoutes(router, diContainer.UserHandler())
	httpinf.RegisterSessionRoutes(router, diContainer.SessionHandler())
	httpinf.RegisterMessageRoutes(router, diContainer.MessageHandler())
	httpinf.RegisterTaskRoutes(router, diContainer.TaskHandler())
	httpinf.RegisterSkillRoutes(router, diContainer.SkillHandler())
	httpinf.RegisterScheduleRoutes(router, diContainer.ScheduleHandler())
	httpinf.RegisterLogRoutes(router, diContainer.LogHandler())

	// Apply middleware
	handler := httpinf.NewHandlerBuilder(router.Handler()).
		Use(httpinf.Logging).
		Use(httpinf.Recovery).
		Use(httpinf.CORS).
		Use(httpinf.RequestID).
		Build()

	// Create HTTP server
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer := httpinf.NewServer(&httpinf.ServerConfig{
		Addr:         addr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	})

	// Start HTTP server in goroutine
	serverCtx := context.Background()
	go func() {
		if err := httpServer.Start(serverCtx); err != nil {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	logger.Info("Application started successfully, waiting for shutdown signal")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan

	logger.Info("Shutting down gracefully...")

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("Failed to shutdown HTTP server", "error", err)
	}

	// Cleanup DI container
	if err := diContainer.Shutdown(); err != nil {
		logger.Error("Failed to shutdown DI container", "error", err)
	}

	logger.Info("Shutdown complete")
}

// TODO: Create DI container for better dependency management
// type DIContainer struct {
// 	config  *config.Config
// 	logger  logging.Logger
// 	db      database.Database
// 	// Add other dependencies
// }
//
// func NewDIContainer(cfg *config.Config, logger logging.Logger) (*DIContainer, error) {
// 	// Initialize all dependencies
// 	return &DIContainer{
// 		config: cfg,
// 		logger: logger,
// 	}, nil
// }
