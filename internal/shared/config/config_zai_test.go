package config

import (
	"os"
	"testing"
)

func TestLoad_ZAIConfig(t *testing.T) {
	// Create test config with zai provider
	testYAML := `
server:
  host: "0.0.0.0"
  port: 8080

database:
  type: "sqlite"
  path: ":memory:"
  migrations_path: "./migrations"
  max_open_conns: 25
  max_idle_conns: 25
  conn_max_lifetime: 5m

eventbus:
  enabled: false

logging:
  level: "info"
  format: "text"

skills:
  directory: "./skills"
  timeout_sec: 60
  sandbox_enabled: false

channels:
  telegram:
    enabled: false
  discord:
    enabled: false
  web:
    enabled: false

llm:
  default_provider: "zai"
  providers:
    zai:
      api_key: "test-key"
      base_url: "https://api.z.ai/api/paas/v4"
      model: "glm-4"
      temperature: 0.7
      max_tokens: 2000
`
	// Create temp file
	tmpDir := t.TempDir()
	cfgPath := tmpDir + "/test.yml"
	if err := writeFile(cfgPath, testYAML); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load config
	cfg, err := Load(cfgPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify default provider
	if cfg.LLM.DefaultProvider != "zai" {
		t.Errorf("Expected default provider 'zai', got '%s'", cfg.LLM.DefaultProvider)
	}

	// Verify zai provider exists
	zaiConfig, ok := cfg.LLM.Providers["zai"]
	if !ok {
		t.Fatal("zai provider not found in config")
	}

	// Verify all fields
	if zaiConfig.APIKey != "test-key" {
		t.Errorf("Expected api_key 'test-key', got '%s'", zaiConfig.APIKey)
	}
	if zaiConfig.BaseURL != "https://api.z.ai/api/paas/v4" {
		t.Errorf("Expected base_url 'https://api.z.ai/api/paas/v4', got '%s'", zaiConfig.BaseURL)
	}
	if zaiConfig.Model != "glm-4" {
		t.Errorf("Expected model 'glm-4', got '%s'", zaiConfig.Model)
	}
	if zaiConfig.Temperature != 0.7 {
		t.Errorf("Expected temperature 0.7, got %f", zaiConfig.Temperature)
	}
	if zaiConfig.MaxTokens != 2000 {
		t.Errorf("Expected max_tokens 2000, got %d", zaiConfig.MaxTokens)
	}
}

// Helper function to write file
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
