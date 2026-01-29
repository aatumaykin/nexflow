package mock

import (
	"context"
	"testing"

	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMockLLMProvider(t *testing.T) {
	provider := NewMockLLMProvider()

	require.NotNil(t, provider)
	require.NotNil(t, provider.GenerateFunc)
	require.NotNil(t, provider.GenerateWithToolsFunc)
	require.NotNil(t, provider.StreamFunc)
	require.NotNil(t, provider.EstimateCostFunc)
}

func TestMockLLMProvider_Generate(t *testing.T) {
	provider := NewMockLLMProvider()
	ctx := context.Background()

	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	resp, err := provider.Generate(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "assistant", resp.Message.Role)
	assert.Contains(t, resp.Message.Content, "Mock response for 1 messages")
	assert.Equal(t, 10, resp.Tokens.InputTokens)
	assert.Equal(t, 20, resp.Tokens.OutputTokens)
	assert.Equal(t, 30, resp.Tokens.TotalTokens)
}

func TestMockLLMProvider_GenerateWithTools(t *testing.T) {
	provider := NewMockLLMProvider()
	ctx := context.Background()

	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
	}

	tools := []ports.ToolDefinition{
		{Name: "test_tool", Description: "A test tool"},
	}

	resp, err := provider.GenerateWithTools(ctx, req, tools)

	require.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "assistant", resp.Message.Role)
	assert.Contains(t, resp.Message.Content, "1 messages with 1 tools")
}

func TestMockLLMProvider_Stream(t *testing.T) {
	provider := NewMockLLMProvider()
	ctx := context.Background()

	req := ports.CompletionRequest{
		Messages: []ports.Message{},
	}

	ch, err := provider.Stream(ctx, req)

	require.NoError(t, err)
	assert.NotNil(t, ch)

	// Read from channel
	messages := []string{}
	for msg := range ch {
		messages = append(messages, msg)
	}

	assert.Len(t, messages, 3)
	assert.Equal(t, "Mock", messages[0])
	assert.Equal(t, "streaming", messages[1])
	assert.Equal(t, "response", messages[2])
}

func TestMockLLMProvider_EstimateCost(t *testing.T) {
	provider := NewMockLLMProvider()

	req := ports.CompletionRequest{
		Messages: []ports.Message{
			{Role: "user", Content: "test message"},
		},
		MaxTokens: 100,
	}

	cost, err := provider.EstimateCost(req)

	require.NoError(t, err)
	assert.Equal(t, 0.002, cost) // 100 tokens * 0.00002
}
