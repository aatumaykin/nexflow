package zai

// zaiChatRequest represents a request to z.ai chat completion API
type zaiChatRequest struct {
	Model       string       `json:"model"`                 // glm-4.7, glm-4.6, etc.
	Messages    []zaiMessage `json:"messages"`              // Conversation messages
	Temperature float64      `json:"temperature,omitempty"` // 0.0-1.0
	TopP        float64      `json:"top_p,omitempty"`       // 0.01-1.0
	MaxTokens   int          `json:"max_tokens,omitempty"`  // Max output tokens
	Stream      bool         `json:"stream,omitempty"`      // Enable streaming
	Tools       []zaiTool    `json:"tools,omitempty"`       // Function calling tools
	ToolChoice  string       `json:"tool_choice,omitempty"` // "auto"
	Thinking    *zaiThinking `json:"thinking,omitempty"`    // Thinking mode
	Stop        []string     `json:"stop,omitempty"`        // Stop sequences
	RequestID   string       `json:"request_id,omitempty"`  // Unique request ID
}

// zaiMessage represents a message in the conversation
type zaiMessage struct {
	Role       string        `json:"role"`                   // "system", "user", "assistant", "tool"
	Content    interface{}   `json:"content"`                // string or []zaiMultimodalContent
	ToolCalls  []zaiToolCall `json:"tool_calls,omitempty"`   // Tool calls from assistant
	ToolCallID string        `json:"tool_call_id,omitempty"` // For tool responses
}

// zaiMultimodalContent represents multimodal content (text, image, video, file)
type zaiMultimodalContent struct {
	Type     string  `json:"type"`                // "text", "image_url", "video_url", "file_url"
	Text     string  `json:"text,omitempty"`      // Text content
	ImageURL *zaiURL `json:"image_url,omitempty"` // Image URL
	VideoURL *zaiURL `json:"video_url,omitempty"` // Video URL
	FileURL  *zaiURL `json:"file_url,omitempty"`  // File URL
}

// zaiURL represents a URL for multimodal content
type zaiURL struct {
	URL string `json:"url"`
}

// zaiTool represents a tool that the model can call
type zaiTool struct {
	Type     string       `json:"type"`               // "function", "web_search", "retrieval"
	Function *zaiFunction `json:"function,omitempty"` // Function definition
}

// zaiFunction represents a function definition
type zaiFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// zaiToolCall represents a tool call made by the model
type zaiToolCall struct {
	ID       string       `json:"id"`       // Tool call ID
	Type     string       `json:"type"`     // "function"
	Function *zaiFunction `json:"function"` // Function name and arguments
}

// zaiThinking represents thinking mode configuration
type zaiThinking struct {
	Type          string `json:"type"`                     // "enabled" or "disabled"
	ClearThinking *bool  `json:"clear_thinking,omitempty"` // Clear thinking from previous turns
}

// zaiChatResponse represents a response from z.ai chat completion API
type zaiChatResponse struct {
	ID        string         `json:"id"`
	RequestID string         `json:"request_id"`
	Created   int64          `json:"created"`
	Model     string         `json:"model"`
	Choices   []zaiChoice    `json:"choices"`
	Usage     zaiUsage       `json:"usage"`
	WebSearch []zaiWebSearch `json:"web_search,omitempty"`
}

// zaiChoice represents a choice in the response
type zaiChoice struct {
	Index        int                `json:"index"`
	Message      zaiResponseMessage `json:"message"`
	FinishReason string             `json:"finish_reason"` // "stop", "tool_calls", "length", "sensitive", "network_error"
}

// zaiResponseMessage represents a message in the response
type zaiResponseMessage struct {
	Role             string        `json:"role"`
	Content          string        `json:"content"`
	ReasoningContent string        `json:"reasoning_content,omitempty"` // For thinking mode
	ToolCalls        []zaiToolCall `json:"tool_calls,omitempty"`
}

// zaiUsage represents token usage statistics
type zaiUsage struct {
	PromptTokens        int                     `json:"prompt_tokens"`
	CompletionTokens    int                     `json:"completion_tokens"`
	TotalTokens         int                     `json:"total_tokens"`
	PromptTokensDetails *zaiPromptTokensDetails `json:"prompt_tokens_details,omitempty"`
}

// zaiPromptTokensDetails represents details about prompt tokens
type zaiPromptTokensDetails struct {
	CachedTokens int `json:"cached_tokens"`
}

// zaiWebSearch represents a web search result
type zaiWebSearch struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Link        string `json:"link"`
	Media       string `json:"media"`
	Icon        string `json:"icon"`
	Refer       string `json:"refer"`
	PublishDate string `json:"publish_date,omitempty"`
}

// zaiStreamChunk represents a chunk in streaming response
type zaiStreamChunk struct {
	ID      string            `json:"id"`
	Created int64             `json:"created"`
	Model   string            `json:"model"`
	Choices []zaiStreamChoice `json:"choices"`
}

// zaiStreamChoice represents a choice in streaming response
type zaiStreamChoice struct {
	Index        int                   `json:"index"`
	Delta        zaiStreamMessageDelta `json:"delta"`
	FinishReason *string               `json:"finish_reason"`
}

// zaiStreamMessageDelta represents a delta in streaming message
type zaiStreamMessageDelta struct {
	Role             string        `json:"role,omitempty"`
	Content          string        `json:"content,omitempty"`
	ReasoningContent string        `json:"reasoning_content,omitempty"`
	ToolCalls        []zaiToolCall `json:"tool_calls,omitempty"`
}

// zaiErrorResponse represents an error response from z.ai API
type zaiErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
