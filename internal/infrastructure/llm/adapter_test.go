package llm

import (
	"context"
	"testing"

	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProvider implements Provider interface for testing
type mockProvider struct {
	name      string
	available bool
	responses []*CompletionResponse
	err       error
}

func (m *mockProvider) Name() string {
	return m.name
}

func (m *mockProvider) Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	if m.err != nil {
		return nil, m.err
	}
	if len(m.responses) > 0 {
		resp := m.responses[0]
		m.responses = m.responses[1:]
		return resp, nil
	}
	return &CompletionResponse{
		Content:    "default response",
		Model:      req.Model,
		TokensUsed: 10,
	}, nil
}

func (m *mockProvider) Chat(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error) {
	return m.Completion(ctx, req)
}

func (m *mockProvider) IsAvailable(ctx context.Context) bool {
	return m.available
}

func TestNewProviderAdapter(t *testing.T) {
	provider := &mockProvider{name: "test"}
	adapter := NewProviderAdapter(provider)

	require.NotNil(t, adapter)
	assert.IsType(t, &ProviderAdapter{}, adapter)
}

func TestProviderAdapter_Generate(t *testing.T) {
	provider := &mockProvider{
		name: "test",
		responses: []*CompletionResponse{
			{
				Content:    "test response",
				Model:      "gpt-4",
				TokensUsed: 100,
			},
		},
		available: true,
	}
	adapter := NewProviderAdapter(provider)

	ctx := context.Background()
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
		Model:     "gpt-4",
		MaxTokens: 100,
	}

	resp, err := adapter.Generate(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "assistant", resp.Message.Role)
	assert.Equal(t, "test response", resp.Message.Content)
	assert.Equal(t, 50, resp.Tokens.InputTokens)  // 100 / 2
	assert.Equal(t, 50, resp.Tokens.OutputTokens) // 100 - (100/2)
	assert.Equal(t, 100, resp.Tokens.TotalTokens)
}

func TestProviderAdapter_Generate_Error(t *testing.T) {
	provider := &mockProvider{
		name: "test",
		err:  assert.AnError,
	}
	adapter := NewProviderAdapter(provider)

	ctx := context.Background()
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	resp, err := adapter.Generate(ctx, req)

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "LLMProviderAdapter.Generate")
}

func TestProviderAdapter_GenerateWithTools(t *testing.T) {
	provider := &mockProvider{
		name: "test",
		responses: []*CompletionResponse{
			{
				Content:    "test response with tools",
				Model:      "gpt-4",
				TokensUsed: 100,
			},
		},
		available: true,
	}
	adapter := NewProviderAdapter(provider)

	ctx := context.Background()
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	tools := []ports.ToolDefinition{
		{Name: "test_tool", Description: "A test tool"},
	}

	// Note: Current implementation doesn't use tools, but should still work
	resp, err := adapter.GenerateWithTools(ctx, req, tools)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test response with tools", resp.Message.Content)
}

func TestProviderAdapter_Stream(t *testing.T) {
	provider := &mockProvider{
		name: "test",
		responses: []*CompletionResponse{
			{
				Content:    "streaming response",
				Model:      "gpt-4",
				TokensUsed: 100,
			},
		},
		available: true,
	}
	adapter := NewProviderAdapter(provider)

	ctx := context.Background()
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	ch, err := adapter.Stream(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Read from channel
	messages := []string{}
	for msg := range ch {
		messages = append(messages, msg)
	}

	assert.Len(t, messages, 1)
	assert.Equal(t, "streaming response", messages[0])
}

func TestProviderAdapter_Stream_Error(t *testing.T) {
	provider := &mockProvider{
		name: "test",
		err:  assert.AnError,
	}
	adapter := NewProviderAdapter(provider)

	ctx := context.Background()
	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	ch, err := adapter.Stream(ctx, req)

	require.Error(t, err)
	assert.Nil(t, ch)
}

func TestProviderAdapter_EstimateCost(t *testing.T) {
	provider := &mockProvider{name: "test"}
	adapter := NewProviderAdapter(provider)

	req := ports.CompletionRequest{
		Messages:  []ports.Message{},
		MaxTokens: 100,
	}

	cost, err := adapter.EstimateCost(req)

	require.NoError(t, err)
	assert.Equal(t, 0.002, cost) // 100 * 0.00002
}

func TestProviderAdapter_EstimateCost_DefaultTokens(t *testing.T) {
	provider := &mockProvider{name: "test"}
	adapter := NewProviderAdapter(provider)

	req := ports.CompletionRequest{
		Messages:  []ports.Message{},
		MaxTokens: 0, // No max tokens specified
	}

	cost, err := adapter.EstimateCost(req)

	require.NoError(t, err)
	assert.Equal(t, 0.002, cost) // 100 (default) * 0.00002
}

func TestConvertMessages(t *testing.T) {
	messages := []ports.Message{
		{Role: "system", Content: "system message"},
		{Role: "user", Content: "user message"},
		{Role: "assistant", Content: "assistant message"},
	}

	infraMessages := convertMessages(messages)

	require.Len(t, infraMessages, 3)
	assert.Equal(t, "system", infraMessages[0].Role)
	assert.Equal(t, "system message", infraMessages[0].Content)
	assert.Equal(t, "user", infraMessages[1].Role)
	assert.Equal(t, "user message", infraMessages[1].Content)
	assert.Equal(t, "assistant", infraMessages[2].Role)
	assert.Equal(t, "assistant message", infraMessages[2].Content)
}
