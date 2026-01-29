package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	// Access use cases from DI container
	// chatUseCase := diContainer.ChatUseCase()
	// userUseCase := diContainer.UserUseCase()
	// skillUseCase := diContainer.SkillUseCase()
	// scheduleUseCase := diContainer.ScheduleUseCase()

	// TODO: Initialize HTTP server (when ready)
	// server := http.NewServer(&cfg.Server, chatUseCase, logger)
	// go server.Run()

	logger.Info("Application started successfully, waiting for shutdown signal")

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for shutdown signal
	<-sigChan

	logger.Info("Shutting down gracefully...")

	// TODO: Shutdown HTTP server
	// server.Shutdown()

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
