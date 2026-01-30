package integration

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atumaikin/nexflow/internal/application/ports"
	llmadapter "github.com/atumaikin/nexflow/internal/infrastructure/llm"
	"github.com/atumaikin/nexflow/internal/infrastructure/llm/zai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
)

// TestZaiProviderAdapterIntegration tests the full integration of zai provider through the adapter
// This test verifies that zai.Provider implements llm.Provider interface correctly
// and that ProviderAdapter properly adapts it to ports.LLMProvider interface
func TestZaiProviderAdapterIntegration(t *testing.T) {
	// Create mock zai server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request format
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected correct Authorization header")
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected /chat/completions path, got %s", r.URL.Path)
		}

		// Send mock response
		resp := map[string]interface{}{
			"id":      "test-id",
			"created": 1234567890,
			"model":   "glm-4.7",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hello! How can I help you?",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]int{
				"prompt_tokens":     10,
				"completion_tokens": 20,
				"total_tokens":      30,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()

	// Step 1: Create zai provider (implements llm.Provider)
	zaiProvider, err := zai.NewProvider(&zai.Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	require.NoError(t, err)
	require.NotNil(t, zaiProvider)

	// Step 2: Verify zai.Provider implements llm.Provider interface
	var infraProvider llmadapter.Provider = zaiProvider
	assert.Equal(t, "zai", infraProvider.Name())

	// Step 3: Test zai provider directly through llm.Provider interface
	infraReq := &llmadapter.CompletionRequest{
		Messages: []*llmadapter.Message{
			{Role: "user", Content: "Hello"},
		},
	}
	infraResp, err := infraProvider.Chat(context.Background(), infraReq)
	require.NoError(t, err)
	assert.Equal(t, "Hello! How can I help you?", infraResp.Content)
	assert.Equal(t, 30, infraResp.TokensUsed)

	// Step 4: Create adapter that converts llm.Provider to ports.LLMProvider
	adapter := llmadapter.NewProviderAdapter(zaiProvider)
	require.NotNil(t, adapter)

	// Step 5: Test adapter implements ports.LLMProvider interface
	var portProvider ports.LLMProvider = adapter

	// Step 6: Test adapter through ports.LLMProvider interface
	portReq := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "Hello"},
		},
		MaxTokens: 100,
	}

	portResp, err := portProvider.Generate(context.Background(), portReq)
	require.NoError(t, err)
	require.NotNil(t, portResp)

	// Verify response
	assert.Equal(t, "assistant", portResp.Message.Role)
	assert.Equal(t, "Hello! How can I help you?", portResp.Message.Content)

	// Verify token calculation (adapter approximates: input = total/2, output = total - input)
	assert.Equal(t, 15, portResp.Tokens.InputTokens)  // 30 / 2
	assert.Equal(t, 15, portResp.Tokens.OutputTokens) // 30 - 15
	assert.Equal(t, 30, portResp.Tokens.TotalTokens)

	// Step 7: Test streaming through adapter
	streamResp, err := portProvider.Stream(context.Background(), portReq)
	require.NoError(t, err)
	require.NotNil(t, streamResp)

	// Read from stream
	messages := []string{}
	for msg := range streamResp {
		messages = append(messages, msg)
	}

	assert.Len(t, messages, 1)
	assert.Equal(t, "Hello! How can I help you?", messages[0])

	// Step 8: Test cost estimation through adapter
	cost, err := portProvider.EstimateCost(portReq)
	require.NoError(t, err)
	assert.Greater(t, cost, 0.0)
}

// TestZaiProviderAdapterWithStreaming tests streaming integration
func TestZaiProviderAdapterWithStreaming(t *testing.T) {
	// Create mock zai server
	// Note: adapter.Stream() calls adapter.Generate() which calls provider.Chat()
	// So we return JSON response, not SSE
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request format
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		if r.Header.Get("Authorization") != "Bearer test-api-key" {
			t.Errorf("Expected correct Authorization header")
		}

		if r.URL.Path != "/chat/completions" {
			t.Errorf("Expected /chat/completions path, got %s", r.URL.Path)
		}

		// Send mock response
		resp := map[string]interface{}{
			"id":      "test-id",
			"created": 1234567890,
			"model":   "glm-4.7",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hello World!",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]int{
				"prompt_tokens":     5,
				"completion_tokens": 12,
				"total_tokens":      17,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()

	// Create zai provider
	zaiProvider, err := zai.NewProvider(&zai.Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	require.NoError(t, err)

	// Create adapter
	adapter := llmadapter.NewProviderAdapter(zaiProvider)

	// Test streaming (adapter.Stream() returns simulated stream with full content)
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "Hello"},
		},
	}

	ch, err := adapter.Stream(context.Background(), req)
	require.NoError(t, err)

	// Collect stream chunks
	var result string
	for chunk := range ch {
		result += chunk
	}

	assert.Equal(t, "Hello World!", result)
}

// TestZaiProviderAdapterIsAvailable tests provider availability check
func TestZaiProviderAdapterIsAvailable(t *testing.T) {
	// Create mock zai server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"id":      "test-id",
			"created": 1234567890,
			"model":   "glm-4.7",
			"choices": []map[string]interface{}{
				{
					"index": 0,
					"message": map[string]interface{}{
						"role":    "assistant",
						"content": "Hi",
					},
					"finish_reason": "stop",
				},
			},
			"usage": map[string]int{
				"total_tokens": 5,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	logger := slog.Default()

	// Create zai provider
	zaiProvider, err := zai.NewProvider(&zai.Config{
		APIKey:  "test-api-key",
		BaseURL: ts.URL,
		Model:   "glm-4.7",
	}, logger)
	require.NoError(t, err)

	// Create adapter
	adapter := llmadapter.NewProviderAdapter(zaiProvider)

	// Test availability (adapter doesn't expose IsAvailable directly, but provider does)
	assert.True(t, zaiProvider.IsAvailable(context.Background()))

	// Test that adapter can be used for generation
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "Test"},
		},
	}

	resp, err := adapter.Generate(context.Background(), req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, "Hi", resp.Message.Content)
}
