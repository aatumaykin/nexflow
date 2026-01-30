package main

import (
	"context"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/application/router"
	channelmock "github.com/atumaikin/nexflow/internal/infrastructure/channels/mock"
	"github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/atumaikin/nexflow/internal/shared/config"
	"github.com/atumaikin/nexflow/internal/shared/eventbus"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

func TestDIContainerInitConnectors(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Channels: config.ChannelsConfig{
			Telegram: config.TelegramConfig{Enabled: true, BotToken: "test-token"},
			Discord:  config.DiscordConfig{Enabled: true, BotToken: "test-token"},
			Web:      config.WebConfig{Enabled: true},
		},
		Database: config.DatabaseConfig{
			Type:           "sqlite",
			Path:           ":memory:",
			MigrationsPath: "../../migrations",
		},
		LLM: config.LLMConfig{
			DefaultProvider: "mock",
			Providers: map[string]config.LLMProvider{
				"mock": {
					APIKey:  "",
					BaseURL: "",
					Model:   "mock",
				},
			},
		},
		Skills: config.SkillsConfig{
			Directory: "",
		},
		EventBus: config.EventBusConfig{
			Enabled: false,
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "text",
		},
	}

	logger := logging.NewNoopLogger()

	// Initialize in-memory database
	db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create DI container
	diContainer, err := NewDIContainer(cfg, logger, db)
	if err != nil {
		t.Fatalf("Failed to create DI container: %v", err)
	}
	defer diContainer.Shutdown()

	// Check that message router was initialized
	messageRouter := diContainer.MessageRouter()
	if messageRouter == nil {
		t.Error("Expected message router to be initialized")
	}

	// Check that all connectors are registered
	connectorNames := messageRouter.ListConnectors()
	expectedConnectors := []string{"telegram", "discord", "web"}

	if len(connectorNames) != len(expectedConnectors) {
		t.Errorf("Expected %d connectors, got %d", len(expectedConnectors), len(connectorNames))
	}

	for _, expected := range expectedConnectors {
		found := false
		for _, name := range connectorNames {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected connector '%s' to be registered", expected)
		}
	}
}

func TestDIContainerNoConnectorsEnabled(t *testing.T) {
	// Create test configuration with no connectors enabled
	cfg := &config.Config{
		Channels: config.ChannelsConfig{
			Telegram: config.TelegramConfig{Enabled: false},
			Discord:  config.DiscordConfig{Enabled: false},
			Web:      config.WebConfig{Enabled: false},
		},
		Database: config.DatabaseConfig{
			Type:           "sqlite",
			Path:           ":memory:",
			MigrationsPath: "../../migrations",
		},
		LLM: config.LLMConfig{
			DefaultProvider: "mock",
			Providers: map[string]config.LLMProvider{
				"mock": {
					APIKey:  "",
					BaseURL: "",
					Model:   "mock",
				},
			},
		},
		Skills: config.SkillsConfig{
			Directory: "",
		},
		EventBus: config.EventBusConfig{
			Enabled: false,
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "text",
		},
	}

	logger := logging.NewNoopLogger()

	// Initialize in-memory database
	db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create DI container
	diContainer, err := NewDIContainer(cfg, logger, db)
	if err != nil {
		t.Fatalf("Failed to create DI container: %v", err)
	}
	defer diContainer.Shutdown()

	// Check that message router was initialized
	messageRouter := diContainer.MessageRouter()
	if messageRouter == nil {
		t.Error("Expected message router to be initialized")
	}

	// Check that no connectors are registered
	connectorNames := messageRouter.ListConnectors()
	if len(connectorNames) != 0 {
		t.Errorf("Expected no connectors, got %d", len(connectorNames))
	}
}

func TestMessageRouterIntegration(t *testing.T) {
	// Create a message router with mock connectors
	logger := logging.NewNoopLogger()
	eventBus := eventbus.NewEventBus(nil)
	router := router.NewMessageRouter(nil, eventBus, logger)

	// Create mock connectors
	telegramConn := channelmock.NewTelegramConnector()
	discordConn := channelmock.NewDiscordConnector()
	webConn := channelmock.NewWebConnector()

	// Register connectors
	router.RegisterConnector(telegramConn)
	router.RegisterConnector(discordConn)
	router.RegisterConnector(webConn)

	// Verify registration
	connectorNames := router.ListConnectors()
	if len(connectorNames) != 3 {
		t.Errorf("Expected 3 connectors, got %d", len(connectorNames))
	}

	// Start the router
	if err := router.Start(); err != nil {
		t.Fatalf("Failed to start router: %v", err)
	}
	defer router.Stop()

	// Verify connectors are running
	if !telegramConn.IsRunning() {
		t.Error("Expected telegram connector to be running")
	}
	if !discordConn.IsRunning() {
		t.Error("Expected discord connector to be running")
	}
	if !webConn.IsRunning() {
		t.Error("Expected web connector to be running")
	}

	// Send test messages
	telegramConn.SendTestMessage("user-1", "chat-1", "Hello from Telegram")
	discordConn.SendTestMessage("user-2", "server-1", "Hello from Discord")
	webConn.SendTestMessage("user-3", "session-1", "Hello from Web")

	// Wait a bit for message processing
	time.Sleep(100 * time.Millisecond)

	t.Log("Integration test completed successfully")
}

func TestConnectorGetters(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Channels: config.ChannelsConfig{
			Telegram: config.TelegramConfig{Enabled: true, BotToken: "test-token"},
		},
		Database: config.DatabaseConfig{
			Type:           "sqlite",
			Path:           ":memory:",
			MigrationsPath: "../../migrations",
		},
		LLM: config.LLMConfig{
			DefaultProvider: "mock",
			Providers: map[string]config.LLMProvider{
				"mock": {
					APIKey:  "",
					BaseURL: "",
					Model:   "mock",
				},
			},
		},
		Skills: config.SkillsConfig{
			Directory: "",
		},
		EventBus: config.EventBusConfig{
			Enabled: false,
		},
		Logging: config.LoggingConfig{
			Level:  "debug",
			Format: "text",
		},
	}

	logger := logging.NewNoopLogger()

	// Initialize in-memory database
	db, err := database.NewDatabase(&cfg.Database, database.WithLogger(logger))
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// Run migrations
	ctx := context.Background()
	if err := db.Migrate(ctx); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Create DI container
	diContainer, err := NewDIContainer(cfg, logger, db)
	if err != nil {
		t.Fatalf("Failed to create DI container: %v", err)
	}
	defer diContainer.Shutdown()

	messageRouter := diContainer.MessageRouter()

	// Test GetConnector for each registered connector
	connectors := messageRouter.ListConnectors()
	for _, name := range connectors {
		conn, exists := messageRouter.GetConnector(name)
		if !exists {
			t.Errorf("Expected connector '%s' to exist", name)
		}
		if conn == nil {
			t.Errorf("Expected non-nil connector for '%s'", name)
		}
		if conn.Name() != name {
			t.Errorf("Expected connector name '%s', got '%s'", name, conn.Name())
		}
	}

	// Test GetConnector for non-existent connector
	_, exists := messageRouter.GetConnector("non-existent")
	if exists {
		t.Error("Expected non-existent connector to not exist")
	}
}
