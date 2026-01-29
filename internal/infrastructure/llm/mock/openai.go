package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
)

// OpenAIProvider is a mock implementation of OpenAI LLM provider
type OpenAIProvider struct {
	name      string
	available bool
	responses []mockResponse
}

// NewOpenAIProvider creates a new mock OpenAI provider
func NewOpenAIProvider() *OpenAIProvider {
	return &OpenAIProvider{
		name:      "openai",
		available: true,
		responses: make([]mockResponse, 0),
	}
}

// Name returns the name of the provider
func (p *OpenAIProvider) Name() string {
	return p.name
}

// Completion generates a text completion
func (p *OpenAIProvider) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("openai provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock OpenAI response",
		Model:      req.Model,
		TokensUsed: 10,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// Chat generates a chat completion
func (p *OpenAIProvider) Chat(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("openai provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock OpenAI chat response",
		Model:      req.Model,
		TokensUsed: 15,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// IsAvailable checks if the provider is available
func (p *OpenAIProvider) IsAvailable(ctx context.Context) bool {
	return p.available
}

// SetAvailable sets the availability status
func (p *OpenAIProvider) SetAvailable(available bool) {
	p.available = available
}

// SetResponse sets the next response to return
func (p *OpenAIProvider) SetResponse(response *llm.CompletionResponse) {
	p.responses = append(p.responses, mockResponse{response: response})
}

// ClearResponses clears all queued responses
func (p *OpenAIProvider) ClearResponses() {
	p.responses = make([]mockResponse, 0)
}
