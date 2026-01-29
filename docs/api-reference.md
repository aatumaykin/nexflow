# Справочник по API Nexflow

Это документация по всем интерфейсам, типам и функциям проекта Nexflow.

## Содержание

- [Application Ports](#application-ports)
- [Domain Entities](#domain-entities)
- [Domain Repositories](#domain-repositories)
- [Application DTOs](#application-dtos)
- [Application Use Cases](#application-use-cases)
- [Infrastructure Components](#infrastructure-components)
- [Shared Utilities](#shared-utilities)

## Application Ports

### Connector

Интерфейс для communication channels (Telegram, Discord, Web, etc.).

```go
package ports

// Event represents an incoming event from a communication channel.
type Event struct {
    ID        string            `json:"id"`         // Unique identifier for the event
    Channel   string            `json:"channel"`    // Channel type: "telegram", "discord", "web", etc.
    UserID    string            `json:"user_id"`    // ID of the user who sent the event
    Message   string            `json:"message"`    // Event message content
    Metadata  map[string]string `json:"metadata"`   // Additional event metadata
    Timestamp string            `json:"timestamp"`  // ISO 8601 format timestamp
}

// Connector defines the interface for communication channels.
type Connector interface {
    // Start begins listening for events from the channel.
    Start(ctx context.Context) error

    // Stop gracefully stops the connector.
    Stop() error

    // Events returns a read-only channel that receives incoming events.
    Events() <-chan Event

    // SendMessage sends a response message to the specified user.
    SendMessage(ctx context.Context, userID, message string) error

    // Name returns the connector name.
    Name() string
}
```

### LLMProvider

Интерфейс для LLM providers (Anthropic, OpenAI, Ollama, etc.).

```go
package ports

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
```

### SkillRuntime

Интерфейс для выполнения навыков (skills).

```go
package ports

// SkillExecution represents the result of a skill execution.
type SkillExecution struct {
    Success bool   `json:"success"`       // Whether the execution was successful
    Output  string `json:"output"`        // Output result in JSON format
    Error   string `json:"error,omitempty"` // Error message if execution failed
}

// SkillRuntime defines the interface for skill execution.
type SkillRuntime interface {
    // Execute runs a skill with the given input parameters.
    Execute(ctx context.Context, skillName string, input map[string]interface{}) (*SkillExecution, error)

    // Validate checks if a skill is valid (permissions, configuration, etc.).
    Validate(skillName string) error

    // List returns all available skill names.
    List() ([]string, error)

    // GetSkill returns skill details by name.
    GetSkill(skillName string) (map[string]interface{}, error)
}
```

## Domain Entities

### User

```go
package entity

// User represents a user in the system.
// A user can interact through different channels (Telegram, Discord, Web, etc.).
type User struct {
    ID        string    `json:"id"`          // Unique identifier for the user
    Channel   string    `json:"channel"`     // Channel type: "telegram", "discord", "web", etc.
    ChannelID string    `json:"channel_id"`  // Channel-specific user identifier
    CreatedAt time.Time `json:"created_at"`  // Timestamp when the user was created
}

// NewUser creates a new user with the specified channel and channel ID.
func NewUser(channel, channelID string) *User

// CanAccessSession checks if the user can access the specified session.
func (u *User) CanAccessSession(sessionID string) bool

// GetChannelUserID returns the channel-specific user identifier.
func (u *User) GetChannelUserID() string

// IsSameChannel returns true if the user is from the same channel as the other user.
func (u *User) IsSameChannel(other *User) bool
```

### Session

```go
package entity

// Session represents a conversation session between a user and the AI.
// A session contains all messages exchanged during a conversation.
type Session struct {
    ID        string    `json:"id"`         // Unique identifier for the session
    UserID    string    `json:"user_id"`    // ID of the user who owns this session
    CreatedAt time.Time `json:"created_at"` // Timestamp when the session was created
    UpdatedAt time.Time `json:"updated_at"` // Timestamp when the session was last updated
}

// NewSession creates a new session for the specified user.
func NewSession(userID string) *Session

// UpdateTimestamp updates the last modified timestamp to the current time.
func (s *Session) UpdateTimestamp()

// IsOwnedBy returns true if the session belongs to the specified user.
func (s *Session) IsOwnedBy(userID string) bool
```

### Message

```go
package entity

// Message represents a message in a conversation session.
// Messages can be from user, assistant (AI), or system.
type Message struct {
    ID        string    `json:"id"`          // Unique identifier for the message
    SessionID string    `json:"session_id"`  // ID of the session this message belongs to
    Role      string    `json:"role"`        // Message role: "user", "assistant", "system"
    Content   string    `json:"content"`     // Message content
    CreatedAt time.Time `json:"created_at"`  // Timestamp when the message was created
}

// MessageRole represents the role of the message sender.
type MessageRole string

const (
    RoleUser      MessageRole = "user"      // Message from a human user
    RoleAssistant MessageRole = "assistant" // Message from the AI assistant
    RoleSystem    MessageRole = "system"    // System-level message
)

// NewUserMessage creates a new user message in the specified session.
func NewUserMessage(sessionID, content string) *Message

// NewAssistantMessage creates a new assistant (AI) message in the specified session.
func NewAssistantMessage(sessionID, content string) *Message

// NewSystemMessage creates a new system message in the specified session.
func NewSystemMessage(sessionID, content string) *Message

// IsFromUser returns true if the message is from a user.
func (m *Message) IsFromUser() bool

// IsFromAssistant returns true if the message is from the AI assistant.
func (m *Message) IsFromAssistant() bool

// IsSystem returns true if the message is a system message.
func (m *Message) IsSystem() bool

// IsPartOfSession returns true if the message belongs to the specified session.
func (m *Message) IsPartOfSession(sessionID string) bool
```

### Task

```go
package entity

// Task represents a skill execution task.
// Tasks track skill execution, status, and results.
type Task struct {
    ID        string    `json:"id"`         // Unique identifier for the task
    SessionID string    `json:"session_id"` // ID of the session this task belongs to
    Skill     string    `json:"skill"`      // Name of the skill to execute
    Input     string    `json:"input"`      // Input parameters in JSON format
    Output    string    `json:"output"`     // Output result in JSON format
    Status    string    `json:"status"`     // Task status: "pending", "running", "completed", "failed"
    Error     string    `json:"error"`      // Error message if the task failed
    CreatedAt time.Time `json:"created_at"` // Timestamp when the task was created
    UpdatedAt time.Time `json:"updated_at"` // Timestamp when the task was last updated
}

// TaskStatus represents the status of a task.
type TaskStatus string

const (
    TaskStatusPending   TaskStatus = "pending"   // Task is waiting to be executed
    TaskStatusRunning   TaskStatus = "running"   // Task is currently running
    TaskStatusCompleted TaskStatus = "completed" // Task completed successfully
    TaskStatusFailed    TaskStatus = "failed"    // Task failed with an error
)

// NewTask creates a new pending task for the specified session and skill with input parameters.
func NewTask(sessionID, skill, input string) *Task

// SetRunning sets the task status to running and updates the timestamp.
func (t *Task) SetRunning()

// SetCompleted sets the task status to completed with the output and updates the timestamp.
func (t *Task) SetCompleted(output string)

// SetFailed sets the task status to failed with an error message and updates the timestamp.
func (t *Task) SetFailed(err error)

// IsPending returns true if the task is pending.
func (t *Task) IsPending() bool

// IsRunning returns true if the task is currently running.
func (t *Task) IsRunning() bool

// IsCompleted returns true if the task completed successfully.
func (t *Task) IsCompleted() bool

// IsFailed returns true if the task failed.
func (t *Task) IsFailed() bool

// BelongsToSession returns true if the task belongs to the specified session.
func (t *Task) BelongsToSession(sessionID string) bool

// GetInput parses and returns the input parameters as a map.
func (t *Task) GetInput() map[string]interface{}

// GetOutput parses and returns the output result as a map.
func (t *Task) GetOutput() map[string]interface{}
```

### Skill

```go
package entity

// Skill represents a registered skill that can be executed by the AI.
// Skills are tools with specific permissions and metadata.
type Skill struct {
    ID          string                 `json:"id"`          // Unique identifier for the skill
    Name        string                 `json:"name"`        // Unique skill name
    Version     string                 `json:"version"`     // Skill version (e.g., "1.0.0")
    Location    string                 `json:"location"`    // Path to skill directory
    Permissions string                 `json:"permissions"` // JSON array of required permissions
    Metadata    string                 `json:"metadata"`    // JSON metadata (timeout, description, etc.)
    CreatedAt   time.Time              `json:"created_at"`  // Timestamp when the skill was registered
    MetadataMap map[string]interface{} `json:"-"`          // Parsed metadata (not persisted)
}

// NewSkill creates a new skill with the specified name, version, location, permissions, and metadata.
func NewSkill(name, version, location string, permissions []string, metadata map[string]interface{}) *Skill

// GetPermissions parses and returns the list of permissions.
func (s *Skill) GetPermissions() []string

// RequiresPermission checks if the skill requires a specific permission.
func (s *Skill) RequiresPermission(permission string) bool

// RequiresSandbox checks if the skill needs sandbox execution.
func (s *Skill) RequiresSandbox() bool

// GetTimeout returns the execution timeout in seconds.
func (s *Skill) GetTimeout() int

// HasPermission checks if the skill has a specific permission.
func (s *Skill) HasPermission(perm string) bool

// GetMetadata parses and returns the metadata as a map.
func (s *Skill) GetMetadata() map[string]interface{}
```

### Schedule

```go
package entity

// Schedule represents a cron-based scheduled task.
// Schedules allow automatic skill execution at specific times defined by cron expressions.
type Schedule struct {
    ID             string    `json:"id"`              // Unique identifier for the schedule
    Skill          string    `json:"skill"`            // Name of the skill to execute
    CronExpression string    `json:"cron_expression"`  // Cron syntax (e.g., "0 * * * *")
    Input          string    `json:"input"`            // Input parameters in JSON format
    Enabled        bool      `json:"enabled"`          // Whether the schedule is active
    CreatedAt      time.Time `json:"created_at"`       // Timestamp when the schedule was created
}

// NewSchedule creates a new enabled schedule for the specified skill with a cron expression and input.
func NewSchedule(skill, cronExpression, input string) *Schedule

// Enable sets the schedule as enabled.
func (s *Schedule) Enable()

// Disable sets the schedule as disabled.
func (s *Schedule) Disable()

// IsEnabled returns true if the schedule is enabled.
func (s *Schedule) IsEnabled() bool

// BelongsToSkill returns true if the schedule belongs to the specified skill.
func (s *Schedule) BelongsToSkill(skill string) bool

// GetInput parses and returns the input parameters as a map.
func (s *Schedule) GetInput() map[string]interface{}
```

### Log

```go
package entity

// Log represents an application log entry.
// Logs are stored in the database for observability and debugging.
type Log struct {
    ID        string    `json:"id"`         // Unique identifier for the log entry
    Level     string    `json:"level"`      // Log level: "debug", "info", "warn", "error"
    Source    string    `json:"source"`     // Source component/module that generated the log
    Message   string    `json:"message"`    // Log message content
    Metadata  string    `json:"metadata"`   // Additional metadata in JSON format
    CreatedAt time.Time `json:"created_at"` // Timestamp when the log was created
}

// LogLevel represents the severity level of a log message.
type LogLevel string

const (
    LogLevelDebug LogLevel = "debug" // Debug level for detailed information
    LogLevelInfo  LogLevel = "info"  // Info level for general information
    LogLevelWarn  LogLevel = "warn"  // Warning level for potential issues
    LogLevelError LogLevel = "error" // Error level for errors and failures
)

// NewLog creates a new log entry with the specified level, source, message, and metadata.
func NewLog(level LogLevel, source, message string, metadata map[string]interface{}) *Log

// IsDebug returns true if the log is at debug level.
func (l *Log) IsDebug() bool

// IsInfo returns true if the log is at info level.
func (l *Log) IsInfo() bool

// IsWarn returns true if the log is at warn level.
func (l *Log) IsWarn() bool

// IsError returns true if the log is at error level.
func (l *Log) IsError() bool

// IsFromSource returns true if the log originated from the specified source.
func (l *Log) IsFromSource(source string) bool

// GetMetadata parses and returns the metadata as a map.
func (l *Log) GetMetadata() map[string]interface{}
```

## Domain Repositories

```go
package repository

import (
    "context"
    "github.com/atumaikin/nexflow/internal/domain/entity"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    GetByID(ctx context.Context, id string) (*entity.User, error)
    GetByChannel(ctx context.Context, channel, channelID string) (*entity.User, error)
    List(ctx context.Context) ([]*entity.User, error)
    Delete(ctx context.Context, id string) error
}

// SessionRepository defines the interface for session data access.
type SessionRepository interface {
    Create(ctx context.Context, session *entity.Session) error
    GetByID(ctx context.Context, id string) (*entity.Session, error)
    GetByUserID(ctx context.Context, userID string) ([]*entity.Session, error)
    Update(ctx context.Context, session *entity.Session) error
    Delete(ctx context.Context, id string) error
}

// MessageRepository defines the interface for message data access.
type MessageRepository interface {
    Create(ctx context.Context, message *entity.Message) error
    GetByID(ctx context.Context, id string) (*entity.Message, error)
    GetBySessionID(ctx context.Context, sessionID string) ([]*entity.Message, error)
    Delete(ctx context.Context, id string) error
}

// TaskRepository defines the interface for task data access.
type TaskRepository interface {
    Create(ctx context.Context, task *entity.Task) error
    GetByID(ctx context.Context, id string) (*entity.Task, error)
    GetBySessionID(ctx context.Context, sessionID string) ([]*entity.Task, error)
    Update(ctx context.Context, task *entity.Task) error
    Delete(ctx context.Context, id string) error
}

// SkillRepository defines the interface for skill data access.
type SkillRepository interface {
    Create(ctx context.Context, skill *entity.Skill) error
    GetByID(ctx context.Context, id string) (*entity.Skill, error)
    GetByName(ctx context.Context, name string) (*entity.Skill, error)
    List(ctx context.Context) ([]*entity.Skill, error)
    Delete(ctx context.Context, id string) error
}

// ScheduleRepository defines the interface for schedule data access.
type ScheduleRepository interface {
    Create(ctx context.Context, schedule *entity.Schedule) error
    GetByID(ctx context.Context, id string) (*entity.Schedule, error)
    GetBySkill(ctx context.Context, skill string) ([]*entity.Schedule, error)
    Update(ctx context.Context, schedule *entity.Schedule) error
    Delete(ctx context.Context, id string) error
}

// LogRepository defines the interface for log data access.
type LogRepository interface {
    Create(ctx context.Context, log *entity.Log) error
    GetByID(ctx context.Context, id string) (*entity.Log, error)
    List(ctx context.Context, source string, limit int) ([]*entity.Log, error)
    Delete(ctx context.Context, id string) error
}
```

## Application DTOs

### User DTOs

```go
package dto

// UserDTO represents a user data transfer object.
type UserDTO struct {
    ID        string `json:"id"`          // Unique identifier for the user
    Channel   string `json:"channel"`     // Channel type: "telegram", "discord", "web", etc.
    ChannelID string `json:"channel_id"`  // Channel-specific user identifier
    CreatedAt string `json:"created_at"`  // ISO 8601 format timestamp when the user was created
}

// CreateUserRequest represents a request to create a new user.
type CreateUserRequest struct {
    Channel   string `json:"channel" yaml:"channel"`     // Channel type
    ChannelID string `json:"channel_id" yaml:"channel_id"` // Channel-specific user identifier
}

// UpdateUserRequest represents a request to update an existing user.
type UpdateUserRequest struct {
    ChannelID string `json:"channel_id,omitempty" yaml:"channel_id,omitempty"` // New channel ID (optional)
}

// UserResponse represents a response containing a single user.
type UserResponse struct {
    Success bool     `json:"success"`     // Whether the operation was successful
    User    *UserDTO `json:"user,omitempty"` // User data (if successful)
    Error   string   `json:"error,omitempty"` // Error message (if failed)
}

// UsersResponse represents a response containing multiple users.
type UsersResponse struct {
    Success bool       `json:"success"`      // Whether the operation was successful
    Users   []*UserDTO `json:"users,omitempty"` // List of users (if successful)
    Error   string     `json:"error,omitempty"` // Error message (if failed)
}
```

### Session DTOs

```go
package dto

// SessionDTO represents a session data transfer object.
type SessionDTO struct {
    ID        string `json:"id"`         // Unique identifier for the session
    UserID    string `json:"user_id"`    // ID of the user who owns the session
    CreatedAt string `json:"created_at"` // ISO 8601 format timestamp when the session was created
    UpdatedAt string `json:"updated_at"` // ISO 8601 format timestamp when the session was last updated
}

// CreateSessionRequest represents a request to create a new session.
type CreateSessionRequest struct {
    UserID string `json:"user_id" yaml:"user_id"` // ID of the user who will own the session
}

// UpdateSessionRequest represents a request to update an existing session.
type UpdateSessionRequest struct {
    UserID string `json:"user_id,omitempty" yaml:"user_id,omitempty"` // New user ID (optional)
}

// SessionResponse represents a response containing a single session.
type SessionResponse struct {
    Success bool        `json:"success"`       // Whether the operation was successful
    Session *SessionDTO `json:"session,omitempty"` // Session data (if successful)
    Error   string      `json:"error,omitempty"`   // Error message (if failed)
}

// SessionsResponse represents a response containing multiple sessions.
type SessionsResponse struct {
    Success  bool          `json:"success"`       // Whether the operation was successful
    Sessions []*SessionDTO `json:"sessions,omitempty"` // List of sessions (if successful)
    Error    string        `json:"error,omitempty"`   // Error message (if failed)
}
```

## Shared Utilities

### Time Utilities

```go
package utils

// Now returns the current time in UTC.
func Now() time.Time

// FormatTimeRFC3339 formats a time.Time to RFC3339 string.
func FormatTimeRFC3339(t time.Time) string

// ParseTimeRFC3339 parses an RFC3339 string to time.Time.
// Returns zero time if parsing fails.
func ParseTimeRFC3339(s string) time.Time
```

### ID Generation

```go
package utils

// GenerateID generates a new UUID-based ID.
func GenerateID() string
```

### JSON Utilities

```go
package utils

// MarshalJSON marshals any value to JSON string.
// Returns "{}" if marshaling fails.
func MarshalJSON(v interface{}) string

// UnmarshalJSONToMap unmarshals JSON string to map[string]interface{}.
// Returns nil if unmarshaling fails.
func UnmarshalJSONToMap(s string) map[string]interface{}

// UnmarshalJSONToSlice unmarshals JSON string to []string.
// Returns nil if unmarshaling fails.
func UnmarshalJSONToSlice(s string) []string
```

## Интеграция с внешними системами

### Конфигурация через ENV

Конфигурация поддерживает подстановку переменных окружения:

```yaml
database:
  password: "${DB_PASSWORD}"  # Будет заменено значением из ENV

llm:
  providers:
    anthropic:
      api_key: "${ANTHROPIC_API_KEY}"
```

### Секреты в логах

Logger автоматически маскирует поля с ключами: `token`, `key`, `password`, `secret`.

```go
logger.Info("Connecting to database",
    "host", "localhost",
    "password", "secret123", // Будет замаскировано
)
// Output: {"password":"***"}
```

## Ресурсы

- [Godoc](https://pkg.go.dev/github.com/atumaikin/nexflow)
- [Go по API](https://golang.org/pkg/)
- [Эффективный Go](https://golang.org/doc/effective_go)
