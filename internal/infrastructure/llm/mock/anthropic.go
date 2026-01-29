package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
)

// AnthropicProvider is a mock implementation of Anthropic LLM provider
type AnthropicProvider struct {
	name      string
	available bool
	responses []mockResponse
}

type mockResponse struct {
	response *llm.CompletionResponse
}

// NewAnthropicProvider creates a new mock Anthropic provider
func NewAnthropicProvider() *AnthropicProvider {
	return &AnthropicProvider{
		name:      "anthropic",
		available: true,
		responses: make([]mockResponse, 0),
	}
}

// Name returns the name of the provider
func (p *AnthropicProvider) Name() string {
	return p.name
}

// Completion generates a text completion
func (p *AnthropicProvider) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("anthropic provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock Anthropic response",
		Model:      req.Model,
		TokensUsed: 10,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// Chat generates a chat completion
func (p *AnthropicProvider) Chat(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("anthropic provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock Anthropic chat response",
		Model:      req.Model,
		TokensUsed: 15,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// IsAvailable checks if the provider is available
func (p *AnthropicProvider) IsAvailable(ctx context.Context) bool {
	return p.available
}

// SetAvailable sets the availability status
func (p *AnthropicProvider) SetAvailable(available bool) {
	p.available = available
}

// SetResponse sets the next response to return
func (p *AnthropicProvider) SetResponse(response *llm.CompletionResponse) {
	p.responses = append(p.responses, mockResponse{response: response})
}

// ClearResponses clears all queued responses
func (p *AnthropicProvider) ClearResponses() {
	p.responses = make([]mockResponse, 0)
}
