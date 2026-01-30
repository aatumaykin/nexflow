package openai

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

// Config represents OpenAI provider configuration
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("openai: api_key is required")
	}
	if c.Model == "" {
		return fmt.Errorf("openai: model is required")
	}
	return nil
}

// Provider is an OpenAI LLM provider
type Provider struct {
	config     *Config
	httpClient *http.Client
	logger     *slog.Logger
}

// NewProvider creates a new OpenAI provider
func NewProvider(config *Config, logger *slog.Logger) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set default base URL if not provided
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	return &Provider{
		config: &Config{
			APIKey:  config.APIKey,
			BaseURL: baseURL,
			Model:   config.Model,
		},
		httpClient: &http.Client{},
		logger:     logger.With("provider", "openai"),
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "openai"
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

	// Convert to OpenAI chat completion request
	openaiReq := chatCompletionRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
	}

	// Add optional parameters
	if req.Temperature > 0 {
		openaiReq.Temperature = req.Temperature
	}
	if req.MaxTokens > 0 {
		openaiReq.MaxTokens = req.MaxTokens
	}

	// Marshal request
	reqBody, err := json.Marshal(openaiReq)
	if err != nil {
		return nil, fmt.Errorf("openai: failed to marshal request: %w", err)
	}

	p.logger.Debug("Sending chat completion request",
		"model", model,
		"messages_count", len(req.Messages))

	// Make HTTP request
	url := fmt.Sprintf("%s/chat/completions", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("openai: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("openai: failed to read response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("openai: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("openai: request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var chatResp chatCompletionResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("openai: failed to parse response: %w", err)
	}

	p.logger.Debug("Chat completion successful",
		"model", chatResp.Model,
		"total_tokens", chatResp.Usage.TotalTokens)

	return &llm.CompletionResponse{
		Content:    chatResp.Choices[0].Message.Content,
		Model:      chatResp.Model,
		TokensUsed: chatResp.Usage.TotalTokens,
		Metadata: map[string]interface{}{
			"prompt_tokens":     chatResp.Usage.PromptTokens,
			"completion_tokens": chatResp.Usage.CompletionTokens,
		},
	}, nil
}

// IsAvailable checks if the provider is available
func (p *Provider) IsAvailable(ctx context.Context) bool {
	url := fmt.Sprintf("%s/models", p.config.BaseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Type definitions for OpenAI API

type chatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	Temperature float64       `json:"temperature,omitempty"`
	MaxTokens   int           `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	ID      string       `json:"id"`
	Object  string       `json:"object"`
	Created int64        `json:"created"`
	Model   string       `json:"model"`
	Choices []chatChoice `json:"choices"`
	Usage   usage        `json:"usage"`
}

type chatChoice struct {
	Index        int         `json:"index"`
	Message      chatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Param   string `json:"param,omitempty"`
	} `json:"error"`
}

// convertMessages converts llm.Messages to chat messages
func convertMessages(messages []*llm.Message) []chatMessage {
	chatMessages := make([]chatMessage, len(messages))
	for i, msg := range messages {
		chatMessages[i] = chatMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
	}
	return chatMessages
}
