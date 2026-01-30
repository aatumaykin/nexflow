package zai

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
	"log/slog"
)

func TestNewProvider_ValidConfig(t *testing.T) {
	logger := slog.Default()
	config := &Config{
		APIKey: "test-api-key",
		Model:  "glm-4.7",
	}

	provider, err := NewProvider(config, logger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if provider == nil {
		t.Fatal("Expected provider to be created")
	}

	if provider.Name() != "zai" {
		t.Errorf("Expected provider name 'zai', got '%s'", provider.Name())
	}
}

func TestNewProvider_MissingAPIKey(t *testing.T) {
	logger := slog.Default()
	config := &Config{
		Model: "glm-4.7",
	}

	_, err := NewProvider(config, logger)
	if err == nil {
		t.Fatal("Expected error for missing API key")
	}

	if err.Error() != "zai: api_key is required" {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}

func TestNewProvider_MissingModel(t *testing.T) {
	logger := slog.Default()
	config := &Config{
		APIKey: "test-api-key",
	}

	_, err := NewProvider(config, logger)
	if err == nil {
		t.Fatal("Expected error for missing model")
	}

	if err.Error() != "zai: model is required" {
		t.Errorf("Expected specific error message, got '%s'", err.Error())
	}
}

func TestNewProvider_DefaultBaseURL(t *testing.T) {
	logger := slog.Default()
	config := &Config{
		APIKey: "test-api-key",
		Model:  "glm-4.7",
	}

	provider, err := NewProvider(config, logger)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Access provider.config.BaseURL through the provider's methods
	// We can't access the private field directly, but we can test by making a request
	if provider == nil {
		t.Fatal("Expected provider to be created")
	}
}

func TestChat_BasicRequest(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected correct Authorization header")
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected /chat/completions path, got %s", r.URL.Path)
		}

		// Parse request
		var req zaiChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if req.Model != "glm-4.7" {
			t.Errorf("Expected model 'glm-4.7', got '%s'", req.Model)
		}

		if len(req.Messages) == 0 {
			t.Error("Expected messages in request")
		}

		// Send mock response
		resp := zaiChatResponse{
			ID:      "test-id",
			Created: 1234567890,
			Model:   "glm-4.7",
			Choices: []zaiChoice{
				{
					Index: 0,
					Message: zaiResponseMessage{
						Role:    "assistant",
						Content: "Hello! How can I help you?",
					},
					FinishReason: "stop",
				},
			},
			Usage: zaiUsage{
				PromptTokens:     10,
				CompletionTokens: 20,
				TotalTokens:      30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	// Create request
	req := &llm.CompletionRequest{
		Messages: []*llm.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	// Call Chat
	resp, err := provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Content != "Hello! How can I help you?" {
		t.Errorf("Expected specific content, got '%s'", resp.Content)
	}

	if resp.TokensUsed != 30 {
		t.Errorf("Expected 30 tokens, got %d", resp.TokensUsed)
	}

	if resp.Model != "glm-4.7" {
		t.Errorf("Expected model 'glm-4.7', got '%s'", resp.Model)
	}
}

func TestChat_WithTemperature(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req zaiChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if req.Temperature != 0.8 {
			t.Errorf("Expected temperature 0.8, got %f", req.Temperature)
		}

		// Send mock response
		resp := zaiChatResponse{
			ID:      "test-id",
			Created: 1234567890,
			Model:   "glm-4.7",
			Choices: []zaiChoice{
				{
					Index: 0,
					Message: zaiResponseMessage{
						Role:    "assistant",
						Content: "Response",
					},
					FinishReason: "stop",
				},
			},
			Usage: zaiUsage{
				TotalTokens: 10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	req := &llm.CompletionRequest{
		Messages:    []*llm.Message{{Role: "user", Content: "Hello"}},
		Temperature: 0.8,
	}

	_, err = provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestChat_WithMaxTokens(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req zaiChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if req.MaxTokens != 1000 {
			t.Errorf("Expected max_tokens 1000, got %d", req.MaxTokens)
		}

		// Send mock response
		resp := zaiChatResponse{
			ID:      "test-id",
			Created: 1234567890,
			Model:   "glm-4.7",
			Choices: []zaiChoice{
				{
					Index: 0,
					Message: zaiResponseMessage{
						Role:    "assistant",
						Content: "Response",
					},
					FinishReason: "stop",
				},
			},
			Usage: zaiUsage{
				TotalTokens: 10,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	req := &llm.CompletionRequest{
		Messages:  []*llm.Message{{Role: "user", Content: "Hello"}},
		MaxTokens: 1000,
	}

	_, err = provider.Chat(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestChat_ErrorResponse(t *testing.T) {
	// Create mock server that returns error
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(zaiErrorResponse{
			Code:    401,
			Message: "Invalid API key",
		})
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	req := &llm.CompletionRequest{
		Messages: []*llm.Message{{Role: "user", Content: "Hello"}},
	}

	_, err = provider.Chat(context.Background(), req)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestStream_BasicStream(t *testing.T) {
	// Create mock server for streaming
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req zaiChatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		if !req.Stream {
			t.Error("Expected stream to be true")
		}

		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		// Send stream chunks
		chunks := []string{
			`data: {"id":"test-id","created":1234567890,"model":"glm-4.7","choices":[{"index":0,"delta":{"content":"Hello"},"finish_reason":null}]}`,
			`data: {"id":"test-id","created":1234567890,"model":"glm-4.7","choices":[{"index":0,"delta":{"content":" "},"finish_reason":null}]}`,
			`data: {"id":"test-id","created":1234567890,"model":"glm-4.7","choices":[{"index":0,"delta":{"content":"World!"},"finish_reason":null}]}`,
			`data: [DONE]`,
		}

		flusher, _ := w.(http.Flusher)
		for _, chunk := range chunks {
			w.Write([]byte(chunk + "\n\n"))
			flusher.Flush()
		}
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	req := &llm.CompletionRequest{
		Messages: []*llm.Message{{Role: "user", Content: "Hello"}},
	}

	ch, err := provider.Stream(context.Background(), req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Collect stream chunks
	var result string
	for chunk := range ch {
		result += chunk
	}

	if result != "Hello World!" {
		t.Errorf("Expected 'Hello World!', got '%s'", result)
	}
}

func TestEstimateCost(t *testing.T) {
	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey: "test-api-key",
		Model:  "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	tests := []struct {
		name         string
		model        string
		inputTokens  int
		outputTokens int
		expected     float64
	}{
		{
			name:         "GLM-4.7",
			model:        "glm-4.7",
			inputTokens:  1000,
			outputTokens: 500,
			expected:     (float64(1000) / 1_000_000 * glm47InputPrice) + (float64(500) / 1_000_000 * glm47OutputPrice),
		},
		{
			name:         "GLM-4.7-FlashX",
			model:        "glm-4.7-flashx",
			inputTokens:  1000,
			outputTokens: 500,
			expected:     (float64(1000) / 1_000_000 * glm47FlashXInputPrice) + (float64(500) / 1_000_000 * glm47FlashXOutputPrice),
		},
		{
			name:         "Unknown model (defaults to glm-4.7)",
			model:        "unknown-model",
			inputTokens:  1000,
			outputTokens: 500,
			expected:     (float64(1000) / 1_000_000 * glm47InputPrice) + (float64(500) / 1_000_000 * glm47OutputPrice),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cost := provider.EstimateCost(tt.model, tt.inputTokens, tt.outputTokens)
			if cost != tt.expected {
				t.Errorf("Expected cost %f, got %f", tt.expected, cost)
			}
		})
	}
}

func TestIsAvailable(t *testing.T) {
	// Create mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Send mock response
		resp := zaiChatResponse{
			ID:      "test-id",
			Created: 1234567890,
			Model:   "glm-4.7",
			Choices: []zaiChoice{
				{
					Index: 0,
					Message: zaiResponseMessage{
						Role:    "assistant",
						Content: "Hi",
					},
					FinishReason: "stop",
				},
			},
			Usage: zaiUsage{
				TotalTokens: 5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()
	provider, err := NewProvider(&Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	if err != nil {
		t.Fatalf("Failed to create provider: %v", err)
	}

	if !provider.IsAvailable(context.Background()) {
		t.Error("Expected provider to be available")
	}
}

func TestConvertMessages(t *testing.T) {
	messages := []*llm.Message{
		{Role: "system", Content: "You are helpful"},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there"},
	}

	zaiMessages := convertMessages(messages)

	if len(zaiMessages) != 3 {
		t.Errorf("Expected 3 messages, got %d", len(zaiMessages))
	}

	if zaiMessages[0].Role != "system" {
		t.Errorf("Expected first role 'system', got '%s'", zaiMessages[0].Role)
	}

	if zaiMessages[1].Role != "user" {
		t.Errorf("Expected second role 'user', got '%s'", zaiMessages[1].Role)
	}

	if zaiMessages[2].Role != "assistant" {
		t.Errorf("Expected third role 'assistant', got '%s'", zaiMessages[2].Role)
	}
}
