package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
)

// Config represents Ollama provider configuration
type Config struct {
	BaseURL string
	Model   string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return fmt.Errorf("ollama: base_url is required")
	}
	if c.Model == "" {
		return fmt.Errorf("ollama: model is required")
	}
	return nil
}

// Provider is an Ollama LLM provider
type Provider struct {
	config     *Config
	httpClient *http.Client
	logger     *slog.Logger
}

// NewProvider creates a new Ollama provider
func NewProvider(config *Config, logger *slog.Logger) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set default base URL if not provided
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}

	return &Provider{
		config: &Config{
			BaseURL: baseURL,
			Model:   config.Model,
		},
		httpClient: &http.Client{},
		logger:     logger.With("provider", "ollama"),
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "ollama"
}

// Completion generates a text completion
func (p *Provider) Completion(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	return p.Chat(ctx, req)
}

// Chat generates a chat completion
func (p *Provider) Chat(ctx context.Context, req *llm.CompletionRequest) (*llm.CompletionResponse, error) {
	// Determine model to use
	model := req.Model
	if model == "" {
		model = p.config.Model
	}

	// Convert to Ollama chat request
	ollamaReq := chatRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
		Stream:   false,
	}

	// Add optional parameters
	if req.Temperature > 0 {
		ollamaReq.Options = map[string]interface{}{
			"temperature": req.Temperature,
		}
	}
	if req.MaxTokens > 0 {
		if ollamaReq.Options == nil {
			ollamaReq.Options = make(map[string]interface{})
		}
		ollamaReq.Options["num_predict"] = req.MaxTokens
	}

	// Marshal request
	reqBody, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to marshal request: %w", err)
	}

	p.logger.Debug("Sending chat request",
		"model", model,
		"messages_count", len(req.Messages))

	// Make HTTP request
	url := fmt.Sprintf("%s/api/chat", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to read response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error != "" {
			return nil, fmt.Errorf("ollama: %s", errResp.Error)
		}
		return nil, fmt.Errorf("ollama: request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var chatResp chatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("ollama: failed to parse response: %w", err)
	}

	p.logger.Debug("Chat successful",
		"model", chatResp.Model,
		"eval_count", chatResp.EvalCount,
		"prompt_eval_count", chatResp.PromptEvalCount)

	// Estimate tokens (Ollama doesn't always provide accurate token counts)
	totalTokens := chatResp.PromptEvalCount + chatResp.EvalCount
	if totalTokens == 0 {
		// Rough estimation if not provided
		totalTokens = len(chatResp.Message.Content) / 4
	}

	return &llm.CompletionResponse{
		Content:    chatResp.Message.Content,
		Model:      chatResp.Model,
		TokensUsed: totalTokens,
		Metadata: map[string]interface{}{
			"prompt_eval_count": chatResp.PromptEvalCount,
			"eval_count":        chatResp.EvalCount,
			"done":              chatResp.Done,
		},
	}, nil
}

// IsAvailable checks if the provider is available
func (p *Provider) IsAvailable(ctx context.Context) bool {
	url := fmt.Sprintf("%s/api/tags", p.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Type definitions for Ollama API

type chatRequest struct {
	Model    string                 `json:"model"`
	Messages []ollamaMessage        `json:"messages"`
	Stream   bool                   `json:"stream"`
	Options  map[string]interface{} `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Model           string        `json:"model"`
	CreatedAt       string        `json:"created_at"`
	Message         ollamaMessage `json:"message"`
	Done            bool          `json:"done"`
	PromptEvalCount int           `json:"prompt_eval_count"`
	EvalCount       int           `json:"eval_count"`
	TotalDuration   int64         `json:"total_duration"`
	LoadDuration    int64         `json:"load_duration"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// convertMessages converts llm.Messages to Ollama messages
func convertMessages(messages []*llm.Message) []ollamaMessage {
	ollamaMessages := make([]ollamaMessage, len(messages))
	for i, msg := range messages {
		role := msg.Role
		// Convert system role to user (Ollama doesn't support system role)
		if role == "system" {
			role = "user"
		}
		ollamaMessages[i] = ollamaMessage{
			Role:    role,
			Content: msg.Content,
		}
	}
	return ollamaMessages
}
