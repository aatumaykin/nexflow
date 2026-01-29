package dto

// MessageDTO represents a message data transfer object
type MessageDTO struct {
	ID        string `json:"id"`
	SessionID string `json:"session_id"`
	Role      string `json:"role"` // "user", "assistant", "system"
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"` // ISO 8601 format
}

// CreateMessageRequest represents a request to create a message
type CreateMessageRequest struct {
	SessionID string `json:"session_id" yaml:"session_id"`
	Role      string `json:"role" yaml:"role"` // "user", "assistant", "system"
	Content   string `json:"content" yaml:"content"`
}

// MessageResponse represents a message response
type MessageResponse struct {
	Success bool        `json:"success"`
	Message *MessageDTO `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// MessagesResponse represents a list of messages response
type MessagesResponse struct {
	Success  bool          `json:"success"`
	Messages []*MessageDTO `json:"messages,omitempty"`
	Error    string        `json:"error,omitempty"`
}

// ChatMessage represents a chat message for LLM interaction
type ChatMessage struct {
	Role    string `json:"role"` // "user", "assistant", "system"
	Content string `json:"content"`
}

// SendMessageRequest represents a request to send a message (for chat flow)
type SendMessageRequest struct {
	UserID  string         `json:"user_id" yaml:"user_id"`
	Message ChatMessage    `json:"message" yaml:"message"`
	Options MessageOptions `json:"options,omitempty" yaml:"options,omitempty"`
}

// MessageOptions represents message options
type MessageOptions struct {
	Model     string `json:"model,omitempty" yaml:"model,omitempty"`
	MaxTokens int    `json:"max_tokens,omitempty" yaml:"max_tokens,omitempty"`
}

// SendMessageResponse represents a response to send message
type SendMessageResponse struct {
	Success  bool          `json:"success"`
	Message  *MessageDTO   `json:"message,omitempty"`
	Messages []*MessageDTO `json:"messages,omitempty"` // Full conversation
	Error    string        `json:"error,omitempty"`
}
