package llm

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/application/ports"
)

// ProviderAdapter adapts infrastructure.Provider to ports.LLMProvider
type ProviderAdapter struct {
	provider Provider
}

// NewProviderAdapter creates a new adapter that implements ports.LLMProvider
func NewProviderAdapter(provider Provider) ports.LLMProvider {
	return &ProviderAdapter{
		provider: provider,
	}
}

// Generate implements ports.LLMProvider.Generate
func (a *ProviderAdapter) Generate(ctx context.Context, req ports.CompletionRequest) (*ports.CompletionResponse, error) {
	// Convert ports.CompletionRequest to llm.CompletionRequest
	infraReq := &CompletionRequest{
		Messages:    convertMessages(req.Messages),
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: 0.7, // Default temperature
		Metadata:    make(map[string]interface{}),
	}

	// Use Chat method for better chat support
	resp, err := a.provider.Chat(ctx, infraReq)
	if err != nil {
		return nil, fmt.Errorf("LLMProviderAdapter.Generate: %w", err)
	}

	// Convert llm.CompletionResponse to ports.CompletionResponse
	return &ports.CompletionResponse{
		Message: ports.Message{
			Role:    "assistant",
			Content: resp.Content,
		},
		Tokens: ports.Tokens{
			InputTokens:  resp.TokensUsed / 2, // Approximation
			OutputTokens: resp.TokensUsed - (resp.TokensUsed / 2),
			TotalTokens:  resp.TokensUsed,
		},
	}, nil
}

// GenerateWithTools implements ports.LLMProvider.GenerateWithTools
// Note: Current infrastructure.Provider doesn't support tools, so we'll call Generate
func (a *ProviderAdapter) GenerateWithTools(ctx context.Context, req ports.CompletionRequest, tools []ports.ToolDefinition) (*ports.CompletionResponse, error) {
	// For now, just call Generate without tool support
	// This can be extended when infrastructure.Provider supports tools
	return a.Generate(ctx, req)
}

// Stream implements ports.LLMProvider.Stream
// Note: Current infrastructure.Provider doesn't support streaming, so we'll return a simulated stream
func (a *ProviderAdapter) Stream(ctx context.Context, req ports.CompletionRequest) (<-chan string, error) {
	// Get the full response first
	resp, err := a.Generate(ctx, req)
	if err != nil {
		return nil, err
	}

	// Create a channel and send the full content
	ch := make(chan string, 1)
	go func() {
		defer close(ch)
		ch <- resp.Message.Content
	}()

	return ch, nil
}

// EstimateCost implements ports.LLMProvider.EstimateCost
// Note: Current infrastructure.Provider doesn't provide cost estimation
func (a *ProviderAdapter) EstimateCost(req ports.CompletionRequest) (float64, error) {
	// Simple estimation: $0.00002 per token (approximate)
	totalTokens := req.MaxTokens
	if totalTokens == 0 {
		totalTokens = 100 // Default estimation
	}
	return float64(totalTokens) * 0.00002, nil
}

// convertMessages converts ports.Message slice to llm.Message slice
func convertMessages(messages []ports.Message) []*Message {
	infraMessages := make([]*Message, len(messages))
	for i, msg := range messages {
		infraMessages[i] = &Message{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return infraMessages
}
