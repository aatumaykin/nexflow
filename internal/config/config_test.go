package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadYAML(t *testing.T) {
	// Create temporary YAML config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "sk-test"
      model: "gpt-4"

channels:
  telegram:
    bot_token: "test-token"
    allowed_users: []
  web:
    enabled: true

skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true

logging:
  level: "info"
  format: "json"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if config.Server.Host != "127.0.0.1" {
		t.Errorf("Expected host 127.0.0.1, got %s", config.Server.Host)
	}
	if config.Server.Port != 8080 {
		t.Errorf("Expected port 8080, got %d", config.Server.Port)
	}
	if config.Database.Type != "sqlite" {
		t.Errorf("Expected database type sqlite, got %s", config.Database.Type)
	}
	if config.LLM.DefaultProvider != "openai" {
		t.Errorf("Expected default_provider openai, got %s", config.LLM.DefaultProvider)
	}
	if config.Skills.TimeoutSec != 30 {
		t.Errorf("Expected timeout 30, got %d", config.Skills.TimeoutSec)
	}
}

func TestLoadJSON(t *testing.T) {
	// Create temporary JSON config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	jsonContent := `{
  "server": {
    "host": "0.0.0.0",
    "port": 9090
  },
  "database": {
    "type": "postgres",
    "path": "/var/lib/nexflow.db"
  },
  "llm": {
    "default_provider": "anthropic",
    "providers": {
      "anthropic": {
        "api_key": "sk-ant-test",
        "model": "claude-opus-4"
      }
    }
  },
  "channels": {
    "telegram": {
      "bot_token": "token123",
      "allowed_users": []
    },
    "web": {
      "enabled": false
    }
  },
  "skills": {
    "directory": "/tmp/skills",
    "timeout_sec": 60,
    "sandbox_enabled": false
  },
  "logging": {
    "level": "debug",
    "format": "text"
  }
}`

	if err := os.WriteFile(configPath, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify values
	if config.Server.Host != "0.0.0.0" {
		t.Errorf("Expected host 0.0.0.0, got %s", config.Server.Host)
	}
	if config.Server.Port != 9090 {
		t.Errorf("Expected port 9090, got %d", config.Server.Port)
	}
	if config.Database.Type != "postgres" {
		t.Errorf("Expected database type postgres, got %s", config.Database.Type)
	}
	if config.LLM.DefaultProvider != "anthropic" {
		t.Errorf("Expected default_provider anthropic, got %s", config.LLM.DefaultProvider)
	}
	if config.Skills.TimeoutSec != 60 {
		t.Errorf("Expected timeout 60, got %d", config.Skills.TimeoutSec)
	}
	if config.Logging.Level != "debug" {
		t.Errorf("Expected logging level debug, got %s", config.Logging.Level)
	}
}

func TestEnvVarExpansion(t *testing.T) {
	// Set environment variables
	os.Setenv("TEST_API_KEY", "secret-key-123")
	os.Setenv("TEST_BOT_TOKEN", "token-456")
	defer func() {
		os.Unsetenv("TEST_API_KEY")
		os.Unsetenv("TEST_BOT_TOKEN")
	}()

	// Create temporary config file with env vars
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${TEST_API_KEY}"
      model: "gpt-4"

channels:
  telegram:
    bot_token: "${TEST_BOT_TOKEN}"
    allowed_users: []
  web:
    enabled: true

skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true

logging:
  level: "info"
  format: "json"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify env var expansion
	if config.LLM.Providers["openai"].APIKey != "secret-key-123" {
		t.Errorf("Expected API key to be expanded, got %s", config.LLM.Providers["openai"].APIKey)
	}
	if config.Channels.Telegram.BotToken != "token-456" {
		t.Errorf("Expected bot token to be expanded, got %s", config.Channels.Telegram.BotToken)
	}
}

func TestEnvVarExpansionWithCustomProvider(t *testing.T) {
	// Set environment variables
	os.Setenv("CUSTOM_LLM_URL", "https://custom-llm.example.com")
	os.Setenv("CUSTOM_LLM_KEY", "custom-key-789")
	defer func() {
		os.Unsetenv("CUSTOM_LLM_URL")
		os.Unsetenv("CUSTOM_LLM_KEY")
	}()

	// Create temporary config file with env vars
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"

llm:
  default_provider: "custom"
  providers:
    custom:
      base_url: "${CUSTOM_LLM_URL}"
      api_key: "${CUSTOM_LLM_KEY}"
      model: "custom-model"

channels:
  telegram:
    bot_token: "test-token"
    allowed_users: []
  web:
    enabled: true

skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true

logging:
  level: "info"
  format: "json"
`

	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	config, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify env var expansion
	if config.LLM.Providers["custom"].BaseURL != "https://custom-llm.example.com" {
		t.Errorf("Expected base_url to be expanded, got %s", config.LLM.Providers["custom"].BaseURL)
	}
	if config.LLM.Providers["custom"].APIKey != "custom-key-789" {
		t.Errorf("Expected api_key to be expanded, got %s", config.LLM.Providers["custom"].APIKey)
	}
}

func TestValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		config    string
		wantError bool
		errorMsg  string
	}{
		{
			name: "missing server host",
			config: `server:
  port: 8080
database:
  type: "sqlite"
  path: "./data/nexflow.db"
llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "test"
      model: "gpt-4"
channels:
  telegram:
    bot_token: "test"
    allowed_users: []
  web:
    enabled: true
skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true
logging:
  level: "info"
  format: "json"`,
			wantError: true,
			errorMsg:  "server.host",
		},
		{
			name: "invalid server port",
			config: `server:
  host: "127.0.0.1"
  port: 0
database:
  type: "sqlite"
  path: "./data/nexflow.db"
llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "test"
      model: "gpt-4"
channels:
  telegram:
    bot_token: "test"
    allowed_users: []
  web:
    enabled: true
skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true
logging:
  level: "info"
  format: "json"`,
			wantError: true,
			errorMsg:  "server.port",
		},
		{
			name: "invalid logging level",
			config: `server:
  host: "127.0.0.1"
  port: 8080
database:
  type: "sqlite"
  path: "./data/nexflow.db"
llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "test"
      model: "gpt-4"
channels:
  telegram:
    bot_token: "test"
    allowed_users: []
  web:
    enabled: true
skills:
  directory: "./skills"
  timeout_sec: 30
  sandbox_enabled: true
logging:
  level: "invalid"
  format: "json"`,
			wantError: true,
			errorMsg:  "logging.level",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configPath := filepath.Join(tmpDir, "config.yml")

			if err := os.WriteFile(configPath, []byte(tt.config), 0644); err != nil {
				t.Fatalf("Failed to write test config: %v", err)
			}

			_, err := Load(configPath)
			if (err != nil) != tt.wantError {
				t.Errorf("Load() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if tt.wantError && err != nil {
				if !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Error message should contain %q, got %q", tt.errorMsg, err.Error())
				}
			}
		})
	}
}

func TestUnsupportedFormat(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.xml")

	if err := os.WriteFile(configPath, []byte("<xml></xml>"), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	_, err := Load(configPath)
	if err == nil {
		t.Error("Expected error for unsupported format, got nil")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
