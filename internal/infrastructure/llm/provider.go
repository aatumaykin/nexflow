package llm

import (
	"context"
)

// Message represents a message in a conversation with LLM
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// CompletionRequest represents a request to generate completion
type CompletionRequest struct {
	Messages    []*Message
	Model       string
	Temperature float64
	MaxTokens   int
	Metadata    map[string]interface{}
}

// CompletionResponse represents the response from LLM
type CompletionResponse struct {
	Content    string
	Model      string
	TokensUsed int
	Metadata   map[string]interface{}
}

// Provider defines the interface for LLM providers
type Provider interface {
	// Name returns the name of the provider
	Name() string

	// Completion generates a text completion
	Completion(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// Chat generates a chat completion (for providers that support chat API)
	Chat(ctx context.Context, req *CompletionRequest) (*CompletionResponse, error)

	// IsAvailable checks if the provider is available
	IsAvailable(ctx context.Context) bool
}
