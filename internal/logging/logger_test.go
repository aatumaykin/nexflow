package logging

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		level   string
		format  string
		wantErr bool
	}{
		{
			name:    "Valid JSON logger with info level",
			level:   "info",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "Valid text logger with debug level",
			level:   "debug",
			format:  "text",
			wantErr: false,
		},
		{
			name:    "Invalid level",
			level:   "invalid",
			format:  "json",
			wantErr: true,
		},
		{
			name:    "All valid levels",
			level:   "warn",
			format:  "json",
			wantErr: false,
		},
		{
			name:    "Error level",
			level:   "error",
			format:  "text",
			wantErr: false,
		},
		{
			name:    "Fatal level (maps to error)",
			level:   "fatal",
			format:  "json",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level, tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && logger == nil {
				t.Error("New() returned nil logger without error")
			}
		})
	}
}

func TestSecretMasking(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		value      string
		wantMasked bool
	}{
		{
			name:       "API key should be masked",
			key:        "api_key",
			value:      "sk-1234567890abcdef",
			wantMasked: true,
		},
		{
			name:       "Password should be masked",
			key:        "password",
			value:      "secret123",
			wantMasked: true,
		},
		{
			name:       "Token should be masked",
			key:        "token",
			value:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
			wantMasked: true,
		},
		{
			name:       "Regular field should not be masked",
			key:        "user_id",
			value:      "12345",
			wantMasked: false,
		},
		{
			name:       "Bot token should be masked",
			key:        "bot_token",
			value:      "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11",
			wantMasked: true,
		},
		{
			name:       "Secret should be masked",
			key:        "secret",
			value:      "mysecretvalue",
			wantMasked: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isMasked := shouldMask(tt.key)
			if isMasked != tt.wantMasked {
				t.Errorf("shouldMask(%q) = %v, want %v", tt.key, isMasked, tt.wantMasked)
			}
		})
	}
}

func TestMaskValue(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "Long value",
			value:    "abcdefghijklmnopqrstuvwxyz",
			expected: "ab**********************yz",
		},
		{
			name:     "Short value (<= 4 chars)",
			value:    "abc",
			expected: "***",
		},
		{
			name:     "Exactly 4 chars",
			value:    "abcd",
			expected: "***",
		},
		{
			name:     "API key format",
			value:    "sk-1234567890abcdef",
			expected: "sk***************ef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskValue(tt.value)
			if result != tt.expected {
				t.Errorf("maskValue(%q) = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestLoggerOutput(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger, err := New("info", "json")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Test logging
	logger.Info("Test message", "key", "value")

	// Restore stdout
	w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)

	// Verify JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	// Check required fields
	if logEntry["level"] != "INFO" {
		t.Errorf("Expected level INFO, got %v", logEntry["level"])
	}
	if logEntry["msg"] != "Test message" {
		t.Errorf("Expected msg 'Test message', got %v", logEntry["msg"])
	}
	if logEntry["source"] != "nexflow" {
		t.Errorf("Expected source 'nexflow', got %v", logEntry["source"])
	}
}

func TestSecretMaskingInLogs(t *testing.T) {
	// Capture stdout
	var buf bytes.Buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	logger, err := New("info", "json")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Log with secret
	logger.Info("Test with secret", "api_key", "sk-1234567890abcdef", "user", "test")

	// Restore stdout and capture output
	w.Close()
	os.Stdout = old
	_, _ = buf.ReadFrom(r)

	// Parse JSON output
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Fatalf("Failed to parse log output as JSON: %v", err)
	}

	// Check that API key is masked
	apiKey, ok := logEntry["api_key"].(string)
	if !ok {
		t.Fatal("api_key field not found in log output")
	}
	if apiKey == "sk-1234567890abcdef" {
		t.Error("API key was not masked in log output")
	}
	if !strings.Contains(apiKey, "*") {
		t.Error("Masked API key does not contain asterisks")
	}
}

func TestLoggerWith(t *testing.T) {
	logger, err := New("info", "json")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create logger with additional fields
	loggerWithFields := logger.With("service", "test", "version", "1.0.0")

	if loggerWithFields == nil {
		t.Error("With() returned nil logger")
	}

	// This is a basic test - more detailed testing would require capturing output
	loggerWithFields.Info("Message with fields")
}

func TestLoggerLevels(t *testing.T) {
	tests := []struct {
		name         string
		level        string
		shouldLog    []string // levels that should log
		shouldNotLog []string // levels that should not log
	}{
		{
			name:         "Info level",
			level:        "info",
			shouldLog:    []string{"info", "warn", "error"},
			shouldNotLog: []string{"debug"},
		},
		{
			name:         "Debug level",
			level:        "debug",
			shouldLog:    []string{"debug", "info", "warn", "error"},
			shouldNotLog: []string{},
		},
		{
			name:         "Warn level",
			level:        "warn",
			shouldLog:    []string{"warn", "error"},
			shouldNotLog: []string{"debug", "info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := New(tt.level, "json")
			if err != nil {
				t.Fatalf("Failed to create logger: %v", err)
			}

			// Test that correct levels log (basic smoke test)
			// In real scenario, we would capture output and verify
			for _, level := range tt.shouldLog {
				switch level {
				case "debug":
					logger.Debug("debug message")
				case "info":
					logger.Info("info message")
				case "warn":
					logger.Warn("warn message")
				case "error":
					logger.Error("error message")
				}
			}
		})
	}
}

func TestAllSecretFields(t *testing.T) {
	fields := []string{
		"api_key", "apikey", "apiKey",
		"token", "access_token", "accessToken",
		"password", "pass",
		"secret", "private_key", "privateKey",
		"bot_token", "botToken",
		"auth_token", "authToken",
		"bearer_token", "bearerToken",
	}

	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			if !shouldMask(field) {
				t.Errorf("Field %q should be masked", field)
			}
		})
	}
}
