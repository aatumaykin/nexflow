package zai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/atumaikin/nexflow/internal/infrastructure/llm"
)

// Config represents z.ai provider configuration
type Config struct {
	APIKey  string
	BaseURL string
	Model   string
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("zai: api_key is required")
	}
	if c.Model == "" {
		return fmt.Errorf("zai: model is required")
	}
	return nil
}

// Provider is a z.ai LLM provider
type Provider struct {
	config     *Config
	httpClient *http.Client
	logger     *slog.Logger
}

// NewProvider creates a new z.ai provider
func NewProvider(config *Config, logger *slog.Logger) (*Provider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Set default base URL if not provided
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.z.ai/api/paas/v4"
	}

	return &Provider{
		config: &Config{
			APIKey:  config.APIKey,
			BaseURL: baseURL,
			Model:   config.Model,
		},
		httpClient: &http.Client{},
		logger:     logger.With("provider", "zai"),
	}, nil
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "zai"
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

	// Convert to z.ai chat completion request
	zaiReq := &zaiChatRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
	}

	// Add optional parameters
	if req.Temperature > 0 {
		zaiReq.Temperature = req.Temperature
	}
	if req.MaxTokens > 0 {
		zaiReq.MaxTokens = req.MaxTokens
	}

	// Marshal request
	reqBody, err := json.Marshal(zaiReq)
	if err != nil {
		return nil, fmt.Errorf("zai: failed to marshal request: %w", err)
	}

	p.logger.Debug("Sending chat completion request",
		"model", model,
		"messages_count", len(req.Messages),
		"max_tokens", zaiReq.MaxTokens,
		"temperature", zaiReq.Temperature)

	// Make HTTP request
	url := fmt.Sprintf("%s/chat/completions", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("zai: failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("zai: failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("zai: failed to read response: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		var errResp zaiErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err == nil && errResp.Message != "" {
			return nil, fmt.Errorf("zai: %s (code: %d)", errResp.Message, errResp.Code)
		}
		return nil, fmt.Errorf("zai: request failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse successful response
	var zaiResp zaiChatResponse
	if err := json.Unmarshal(respBody, &zaiResp); err != nil {
		return nil, fmt.Errorf("zai: failed to parse response: %w", err)
	}

	if len(zaiResp.Choices) == 0 {
		return nil, fmt.Errorf("zai: no choices in response")
	}

	p.logger.Debug("Chat completion successful",
		"model", zaiResp.Model,
		"prompt_tokens", zaiResp.Usage.PromptTokens,
		"completion_tokens", zaiResp.Usage.CompletionTokens,
		"total_tokens", zaiResp.Usage.TotalTokens,
		"finish_reason", zaiResp.Choices[0].FinishReason)

	// Extract content from response
	content := zaiResp.Choices[0].Message.Content

	// Build metadata
	metadata := map[string]interface{}{
		"model":             zaiResp.Model,
		"prompt_tokens":     zaiResp.Usage.PromptTokens,
		"completion_tokens": zaiResp.Usage.CompletionTokens,
		"total_tokens":      zaiResp.Usage.TotalTokens,
		"finish_reason":     zaiResp.Choices[0].FinishReason,
	}

	// Add thinking content if present
	if zaiResp.Choices[0].Message.ReasoningContent != "" {
		metadata["reasoning_content"] = zaiResp.Choices[0].Message.ReasoningContent
	}

	return &llm.CompletionResponse{
		Content:    content,
		Model:      zaiResp.Model,
		TokensUsed: zaiResp.Usage.TotalTokens,
		Metadata:   metadata,
	}, nil
}

// Stream generates a streaming completion
func (p *Provider) Stream(ctx context.Context, req *llm.CompletionRequest) (<-chan string, error) {
	// Determine model to use
	model := req.Model
	if model == "" {
		model = p.config.Model
	}

	// Convert to z.ai chat completion request with streaming enabled
	zaiReq := &zaiChatRequest{
		Model:    model,
		Messages: convertMessages(req.Messages),
		Stream:   true,
	}

	// Add optional parameters
	if req.Temperature > 0 {
		zaiReq.Temperature = req.Temperature
	}
	if req.MaxTokens > 0 {
		zaiReq.MaxTokens = req.MaxTokens
	}

	// Marshal request
	reqBody, err := json.Marshal(zaiReq)
	if err != nil {
		return nil, fmt.Errorf("zai: failed to marshal stream request: %w", err)
	}

	p.logger.Debug("Sending streaming chat completion request",
		"model", model,
		"messages_count", len(req.Messages))

	// Make HTTP request
	url := fmt.Sprintf("%s/chat/completions", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("zai: failed to create stream request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("zai: failed to send stream request: %w", err)
	}

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("zai: stream request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create channel for streaming chunks
	ch := make(chan string, 10)

	// Start goroutine to read stream
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()

			// Skip empty lines and keep-alive
			if line == "" || line == ": keep-alive" {
				continue
			}

			// SSE format: "data: {...}"
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")

			// Check for end of stream
			if data == "[DONE]" {
				p.logger.Debug("Stream complete")
				break
			}

			// Parse JSON chunk
			var chunk zaiStreamChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				p.logger.Warn("Failed to parse stream chunk", "error", err, "data", data)
				continue
			}

			// Extract content from delta
			if len(chunk.Choices) > 0 {
				delta := chunk.Choices[0].Delta
				if delta.Content != "" {
					ch <- delta.Content
				}
			}
		}

		if err := scanner.Err(); err != nil {
			p.logger.Error("Stream scanner error", "error", err)
		}
	}()

	return ch, nil
}

// IsAvailable checks if the provider is available
func (p *Provider) IsAvailable(ctx context.Context) bool {
	// Send a minimal request to check availability
	zaiReq := &zaiChatRequest{
		Model: p.config.Model,
		Messages: []zaiMessage{
			{Role: "user", Content: "Hi"},
		},
		MaxTokens: 10,
	}

	reqBody, err := json.Marshal(zaiReq)
	if err != nil {
		return false
	}

	url := fmt.Sprintf("%s/chat/completions", p.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(reqBody))
	if err != nil {
		return false
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// convertMessages converts llm.Messages to z.ai messages
func convertMessages(messages []*llm.Message) []zaiMessage {
	zaiMessages := make([]zaiMessage, 0, len(messages))

	for _, msg := range messages {
		zaiMsg := zaiMessage{
			Role:    msg.Role,
			Content: msg.Content,
		}
		zaiMessages = append(zaiMessages, zaiMsg)
	}

	return zaiMessages
}
