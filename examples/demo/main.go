package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/atumaikin/nexflow/internal/config"
	"github.com/atumaikin/nexflow/internal/database"
	"github.com/atumaikin/nexflow/internal/logging"
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

	logger.Info("Starting Nexflow application",
		"version", "1.0.0",
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

	// Example: Create a user (with secret masking)
	logger.Info("Creating example user", "user_id", "12345")

	// Example: Log with secret (will be automatically masked)
	logger.Info("API request sent",
		"endpoint", "/users",
		"api_key", "sk-1234567890abcdef",
		"user_id", "12345",
	)

	// Example: Log with context
	logger.With("service", "auth").Info("User authenticated", "user_id", "12345")

	// Example: Log with different levels
	logger.Debug("Debug information", "detail", "This is a debug message")
	logger.Info("Application is running", "status", "healthy")
	logger.Warn("Rate limit approaching", "current", 90, "max", 100)

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("Application started successfully, waiting for shutdown signal")

	// Wait for shutdown signal
	<-sigChan

	logger.Info("Shutting down gracefully...")

	// Cleanup
	logger.Info("Shutdown complete")
}
