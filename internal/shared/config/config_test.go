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
  migrations_path: "./migrations"

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
  migrations_path: "./migrations"

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
  migrations_path: "./migrations"

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
  migrations_path: "./migrations"
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
  migrations_path: "./migrations"
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
  migrations_path: "./migrations"
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

// TestUniversalEnvVarExpansion tests that env vars work in any string field
func TestUniversalEnvVarExpansion(t *testing.T) {
	// Set multiple environment variables
	os.Setenv("TEST_HOST", "example.com")
	os.Setenv("TEST_DB_PATH", "/custom/path/db.sqlite")
	os.Setenv("TEST_SKILLS_DIR", "/custom/skills")
	os.Setenv("TEST_LOG_LEVEL", "debug")
	defer func() {
		os.Unsetenv("TEST_HOST")
		os.Unsetenv("TEST_DB_PATH")
		os.Unsetenv("TEST_SKILLS_DIR")
		os.Unsetenv("TEST_LOG_LEVEL")
	}()

	// Create config with env vars in various fields
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
   host: "${TEST_HOST}"
   port: 8080

database:
   type: "sqlite"
   path: "${TEST_DB_PATH}"
   migrations_path: "./migrations"

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
   directory: "${TEST_SKILLS_DIR}"
   timeout_sec: 30
   sandbox_enabled: true

logging:
   level: "${TEST_LOG_LEVEL}"
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

	// Verify env var expansion in all fields
	if config.Server.Host != "example.com" {
		t.Errorf("Expected host example.com, got %s", config.Server.Host)
	}
	if config.Database.Path != "/custom/path/db.sqlite" {
		t.Errorf("Expected database path /custom/path/db.sqlite, got %s", config.Database.Path)
	}
	if config.Skills.Directory != "/custom/skills" {
		t.Errorf("Expected skills dir /custom/skills, got %s", config.Skills.Directory)
	}
	if config.Logging.Level != "debug" {
		t.Errorf("Expected log level debug, got %s", config.Logging.Level)
	}
}

// TestMultipleEnvVarsInOneString tests multiple env vars in a single string
func TestMultipleEnvVarsInOneString(t *testing.T) {
	// Set environment variables
	os.Setenv("PROTOCOL", "https")
	os.Setenv("DOMAIN", "api.example.com")
	os.Setenv("PORT", "8443")
	defer func() {
		os.Unsetenv("PROTOCOL")
		os.Unsetenv("DOMAIN")
		os.Unsetenv("PORT")
	}()

	// Create config with multiple env vars in one string
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"
  migrations_path: "./migrations"

llm:
  default_provider: "custom"
  providers:
    custom:
      base_url: "${PROTOCOL}://${DOMAIN}:${PORT}/v1"
      api_key: "test-key"
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

	// Verify multiple env vars expansion
	expectedURL := "https://api.example.com:8443/v1"
	if config.LLM.Providers["custom"].BaseURL != expectedURL {
		t.Errorf("Expected base_url %s, got %s", expectedURL, config.LLM.Providers["custom"].BaseURL)
	}
}

// TestEnvVarsInSlice tests env var expansion in string slices
func TestEnvVarsInSlice(t *testing.T) {
	// Set environment variable
	os.Setenv("USER_ID_1", "user123")
	os.Setenv("USER_ID_2", "user456")
	defer func() {
		os.Unsetenv("USER_ID_1")
		os.Unsetenv("USER_ID_2")
	}()

	// Create config with env vars in slice
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"
  migrations_path: "./migrations"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "sk-test"
      model: "gpt-4"

channels:
  telegram:
    bot_token: "test-token"
    allowed_users:
      - "${USER_ID_1}"
      - "${USER_ID_2}"
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

	// Verify env var expansion in slice
	if len(config.Channels.Telegram.AllowedUsers) != 2 {
		t.Errorf("Expected 2 allowed users, got %d", len(config.Channels.Telegram.AllowedUsers))
	}
	if config.Channels.Telegram.AllowedUsers[0] != "user123" {
		t.Errorf("Expected user123, got %s", config.Channels.Telegram.AllowedUsers[0])
	}
	if config.Channels.Telegram.AllowedUsers[1] != "user456" {
		t.Errorf("Expected user456, got %s", config.Channels.Telegram.AllowedUsers[1])
	}
}

// TestEnvVarNotFound keeps template if env var not set
func TestEnvVarNotFound(t *testing.T) {
	// Don't set the environment variable

	// Create config with env var that doesn't exist
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yml")

	yamlContent := `server:
  host: "127.0.0.1"
  port: 8080

database:
  type: "sqlite"
  path: "./data/nexflow.db"
  migrations_path: "./migrations"

llm:
  default_provider: "openai"
  providers:
    openai:
      api_key: "${MISSING_VAR}"
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

	// Verify that template is kept when env var is not set
	if config.LLM.Providers["openai"].APIKey != "${MISSING_VAR}" {
		t.Errorf("Expected template ${MISSING_VAR} to be kept, got %s", config.LLM.Providers["openai"].APIKey)
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
