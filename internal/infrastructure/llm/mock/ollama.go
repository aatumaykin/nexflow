package mock

import (
	"context"
	"fmt"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
)

// OllamaProvider is a mock implementation of Ollama LLM provider
type OllamaProvider struct {
	name      string
	available bool
	responses []mockResponse
}

// NewOllamaProvider creates a new mock Ollama provider
func NewOllamaProvider() *OllamaProvider {
	return &OllamaProvider{
		name:      "ollama",
		available: true,
		responses: make([]mockResponse, 0),
	}
}

// Name returns the name of the provider
func (p *OllamaProvider) Name() string {
	return p.name
}

// Completion generates a text completion
func (p *OllamaProvider) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("ollama provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock Ollama response",
		Model:      req.Model,
		TokensUsed: 10,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// Chat generates a chat completion
func (p *OllamaProvider) Chat(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	if !p.available {
		return nil, fmt.Errorf("ollama provider is not available")
	}

	// Return default response or custom response if set
	if len(p.responses) > 0 {
		resp := p.responses[0]
		p.responses = p.responses[1:]
		return resp.response, nil
	}

	return &llm.CompletionResponse{
		Content:    "Mock Ollama chat response",
		Model:      req.Model,
		TokensUsed: 15,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// IsAvailable checks if the provider is available
func (p *OllamaProvider) IsAvailable(ctx context.Context) bool {
	return p.available
}

// SetAvailable sets the availability status
func (p *OllamaProvider) SetAvailable(available bool) {
	p.available = available
}

// SetResponse sets the next response to return
func (p *OllamaProvider) SetResponse(response *llm.CompletionResponse) {
	p.responses = append(p.responses, mockResponse{response: response})
}

// ClearResponses clears all queued responses
func (p *OllamaProvider) ClearResponses() {
	p.responses = make([]mockResponse, 0)
}
