package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
)

// MockLLMProvider is a mock implementation of LLMProvider for testing
type MockLLMProvider struct {
	GenerateFunc          func(context.Context, ports.CompletionRequest) (*ports.CompletionResponse, error)
	GenerateWithToolsFunc func(context.Context, ports.CompletionRequest, []ports.ToolDefinition) (*ports.CompletionResponse, error)
	StreamFunc            func(context.Context, ports.CompletionRequest) (<-chan string, error)
	EstimateCostFunc      func(ports.CompletionRequest) (float64, error)
}

// NewMockLLMProvider creates a new mock LLM provider
func NewMockLLMProvider() *MockLLMProvider {
	return &MockLLMProvider{
		GenerateFunc: func(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
			// Mock response with a simple assistant message
			return &ports.CompletionResponse{
				Message: ports.Message{
					Role:    "assistant",
					Content: fmt.Sprintf("Mock response for %d messages", len(req.Messages)),
				},
				Tokens: ports.Tokens{
					InputTokens:  len(req.Messages) * 10,
					OutputTokens: 20,
					TotalTokens:  len(req.Messages)*10 + 20,
				},
			}, nil
		},
		GenerateWithToolsFunc: func(ctx context.Context, req ports.CompletionRequest, tools []ports.ToolDefinition) (*ports.CompletionResponse, error) {
			// Mock response with tools
			return &ports.CompletionResponse{
				Message: ports.Message{
					Role:    "assistant",
					Content: fmt.Sprintf("Mock response for %d messages with %d tools", len(req.Messages), len(tools)),
				},
				Tokens: ports.Tokens{
					InputTokens:  len(req.Messages) * 10,
					OutputTokens: 25,
					TotalTokens:  len(req.Messages)*10 + 25,
				},
			}, nil
		},
		StreamFunc: func(ctx context.Context, req ports.CompletionRequest) (<-chan string, error) {
			// Mock streaming response
			ch := make(chan string, 1)
			go func() {
				defer close(ch)
				ch <- "Mock"
				ch <- "streaming"
				ch <- "response"
			}()
			return ch, nil
		},
		EstimateCostFunc: func(req ports.CompletionRequest) (float64, error) {
			// Mock cost calculation
			totalTokens := req.MaxTokens
			if totalTokens == 0 {
				totalTokens = 100
			}
			return float64(totalTokens) * 0.00002, nil
		},
	}
}

// Generate implements LLMProvider interface
func (m *MockLLMProvider) Generate(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
	if m.GenerateFunc != nil {
		return m.GenerateFunc(ctx, req)
	}
	return nil, fmt.Errorf("GenerateFunc not set")
}

// GenerateWithTools implements LLMProvider interface
func (m *MockLLMProvider) GenerateWithTools(ctx context.Context, req ports.CompletionRequest, tools []ports.ToolDefinition) (*ports.CompletionResponse, error) {
	if m.GenerateWithToolsFunc != nil {
		return m.GenerateWithToolsFunc(ctx, req, tools)
	}
	return nil, fmt.Errorf("GenerateWithToolsFunc not set")
}

// Stream implements LLMProvider interface
func (m *MockLLMProvider) Stream(ctx context.Context, req ports.CompletionRequest) (<-chan string, error) {
	if m.StreamFunc != nil {
		return m.StreamFunc(ctx, req)
	}
	return nil, fmt.Errorf("StreamFunc not set")
}

// EstimateCost implements LLMProvider interface
func (m *MockLLMProvider) EstimateCost(req ports.CompletionRequest) (float64, error) {
	if m.EstimateCostFunc != nil {
		return m.EstimateCostFunc(req)
	}
	return 0, fmt.Errorf("EstimateCostFunc not set")
}
