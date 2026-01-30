package anthropic

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

// Config represents Anthropic provider configuration
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("anthropic: api_key is required")
	}
	if c.Model == "" {
		return fmt.Errorf("anthropic: model is required")
	}
	return nil
}

// Provider is an Anthropic LLM provider
type Provider struct {
	config     *Config
	httpClient *http.Client
	logger     *slog.Logger
}

// NewProvider creates a new Anthropic provider
func NewProvider(config *Config, logger *slog.Logger) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set default base URL if not provided
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.anthropic.com/v1"
	}

	return &Provider{
		config: &Config{
			APIKey:  config.APIKey,
			BaseURL: baseURL,
			Model:   config.Model,
		},
		httpClient: &http.Client{},
		logger:     logger.With("provider", "anthropic"),
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "anthropic"
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

	// Convert to Anthropic message request
	anthropicReq := messageRequest{
		Model:     model,
		MaxTokens: 4096, // Default max tokens for Anthropic
		Messages:  convertMessages(req.Messages),
	}

	// Add optional parameters
	if req.Temperature > 0 {
		anthropicReq.Temperature = req.Temperature
	}
	if req.MaxTokens > 0 {
		anthropicReq.MaxTokens = req.MaxTokens
	}

	// Marshal request
	reqBody, err := json.Marshal(anthropicReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to marshal request: %w", err)
	}

	p.logger.Debug("Sending message request",
		"model", model,
		"messages_count", len(req.Messages))

	// Make HTTP request
	url := fmt.Sprintf("%s/messages", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", p.config.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("anthropic: failed to read response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("anthropic: %s", errResp.Error.Message)
		}
		return nil, fmt.Errorf("anthropic: request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var msgResp messageResponse
	if err := json.Unmarshal(respBody, &msgResp); err != nil {
		return nil, fmt.Errorf("anthropic: failed to parse response: %w", err)
	}

	p.logger.Debug("Message successful",
		"model", msgResp.Model,
		"input_tokens", msgResp.Usage.InputTokens,
		"output_tokens", msgResp.Usage.OutputTokens)

	return &llm.CompletionResponse{
		Content:    msgResp.Content[0].Text,
		Model:      msgResp.Model,
		TokensUsed: msgResp.Usage.InputTokens + msgResp.Usage.OutputTokens,
		Metadata: map[string]interface{}{
			"input_tokens":  msgResp.Usage.InputTokens,
			"output_tokens": msgResp.Usage.OutputTokens,
			"stop_reason":   msgResp.StopReason,
		},
	}, nil
}

// IsAvailable checks if the provider is available
func (p *Provider) IsAvailable(ctx context.Context) bool {
	url := fmt.Sprintf("%s/messages", p.config.BaseURL)
	reqBody := messageRequest{
		Model:     p.config.Model,
		MaxTokens: 10,
		Messages: []anthropicMessage{
			{Role: "user", Content: "Hi"},
		},
	}

	reqBodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return false
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBodyBytes))
	if err != nil {
		return false
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// Type definitions for Anthropic API

type messageRequest struct {
	Model       string             `json:"model"`
	MaxTokens   int                `json:"max_tokens"`
	Messages    []anthropicMessage `json:"messages"`
	Temperature float64            `json:"temperature,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type messageResponse struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Role       string    `json:"role"`
	Content    []content `json:"content"`
	StopReason string    `json:"stop_reason"`
	Model      string    `json:"model"`
	Usage      usage     `json:"usage"`
}

type content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type errorResponse struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

// convertMessages converts llm.Messages to Anthropic messages
func convertMessages(messages []*llm.Message) []anthropicMessage {
	anthropicMessages := make([]anthropicMessage, 0, len(messages))

	// Combine system messages at the start
	var systemContent string
	for _, msg := range messages {
		if msg.Role == "system" {
			if systemContent != "" {
				systemContent += "\n\n"
			}
			systemContent += msg.Content
		}
	}

	// Add system as first user message if present (Anthropic doesn't have separate system role)
	if systemContent != "" {
		anthropicMessages = append(anthropicMessages, anthropicMessage{
			Role:    "user",
			Content: systemContent,
		})
	}

	// Add remaining messages
	for _, msg := range messages {
		if msg.Role != "system" {
			anthropicMessages = append(anthropicMessages, anthropicMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	return anthropicMessages
}
