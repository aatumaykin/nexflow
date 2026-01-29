package ports

import (
	"context"
)

// Message represents a chat message in a conversation.
type Message struct {
	Role    string `json:"role"`    // Message role: "user", "assistant", "system"
	Content string `json:"content"` // Message content
}

// CompletionRequest represents a request for LLM completion.
type CompletionRequest struct {
	Messages  []Message `json:"messages"`             // Conversation history messages
	Model     string    `json:"model,omitempty"`      // Model to use for completion
	MaxTokens int       `json:"max_tokens,omitempty"` // Maximum tokens in the response
}

// CompletionResponse represents an LLM completion response.
type CompletionResponse struct {
	Message Message `json:"message"` // Generated message
	Tokens  Tokens  `json:"tokens"`  // Token usage information
}

// Tokens represents token usage information for the completion.
type Tokens struct {
	InputTokens  int `json:"input_tokens"`  // Number of tokens in the input
	OutputTokens int `json:"output_tokens"` // Number of tokens in the output
	TotalTokens  int `json:"total_tokens"`  // Total number of tokens used
}

// ToolCall represents a tool/function call made by the LLM.
type ToolCall struct {
	Name      string                 `json:"name"`      // Name of the tool to call
	Arguments map[string]interface{} `json:"arguments"` // Arguments to pass to the tool
}

// LLMProvider defines the interface for LLM providers.
// Different providers (OpenAI, Anthropic, Ollama, etc.) implement this interface.
type LLMProvider interface {
	// Generate generates a completion for the given request.
	Generate(ctx context.Context, req CompletionRequest) (*CompletionResponse, error)

	// GenerateWithTools generates a completion with tool/function calling support.
	GenerateWithTools(ctx context.Context, req CompletionRequest, tools []ToolDefinition) (*CompletionResponse, error)

	// Stream generates a streaming completion, returning a channel of text chunks.
	Stream(ctx context.Context, req CompletionRequest) (<-chan string, error)

	// EstimateCost estimates the cost of a request in dollars (optional).
	EstimateCost(req CompletionRequest) (float64, error)
}

// ToolDefinition defines a tool/function that the LLM can call.
type ToolDefinition struct {
	Name        string      `json:"name"`        // Unique name of the tool
	Description string      `json:"description"` // Description of what the tool does
	Parameters  interface{} `json:"parameters"`  // JSON Schema for the tool parameters
}
