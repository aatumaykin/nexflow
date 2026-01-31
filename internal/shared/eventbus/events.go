package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// BaseEvent provides a base implementation for events
type BaseEvent struct {
	eventType string
	timestamp time.Time
	metadata  map[string]interface{}
	data      interface{}
	mu        sync.RWMutex
}

// NewBaseEvent creates a new base event
//
// Parameters:
//   - eventType: Type identifier for the event
//   - data: Event data (optional)
//
// Returns:
//   - *BaseEvent: Initialized base event
func NewBaseEvent(eventType string, data interface{}) *BaseEvent {
	return &BaseEvent{
		eventType: eventType,
		timestamp: time.Now(),
		metadata:  make(map[string]interface{}),
		data:      data,
	}
}

// Type returns the event type
func (e *BaseEvent) Type() string {
	return e.eventType
}

// Timestamp returns the event timestamp
func (e *BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// Data returns the event data
func (e *BaseEvent) Data() interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.data
}

// SetData sets the event data
func (e *BaseEvent) SetData(data interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data = data
}

// Metadata returns the event metadata
func (e *BaseEvent) Metadata() map[string]interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.metadata
}

// SetMetadata sets the event metadata
func (e *BaseEvent) SetMetadata(metadata map[string]interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metadata = metadata
}

// GetMetadataValue returns a value from metadata
func (e *BaseEvent) GetMetadataValue(key string) (interface{}, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	val, ok := e.metadata[key]
	return val, ok
}

// SetMetadataValue sets a value in metadata
func (e *BaseEvent) SetMetadataValue(key string, value interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.metadata[key] = value
}

// Event types
const (
	// Connector events
	EventConnectorStarted = "connector.started"
	EventConnectorStopped = "connector.stopped"
	EventConnectorError   = "connector.error"
	EventConnectorMessage = "connector.message"

	// Router events
	EventRouterStarted = "router.started"
	EventRouterStopped = "router.stopped"
	EventRouterError   = "router.error"
	EventRouterMessage = "router.message"

	// Orchestrator events
	EventOrchestratorStarted = "orchestrator.started"
	EventOrchestratorStopped = "orchestrator.stopped"
	EventOrchestratorError   = "orchestrator.error"
	EventOrchestratorTask    = "orchestrator.task"

	// LLM events
	EventLLMRequest  = "llm.request"
	EventLLMResponse = "llm.response"
	EventLLMError    = "llm.error"

	// User events
	EventUserCreated = "user.created"
	EventUserUpdated = "user.updated"
	EventUserDeleted = "user.deleted"

	// Session events
	EventSessionCreated = "session.created"
	EventSessionUpdated = "session.updated"
	EventSessionEnded   = "session.ended"

	// Skill events
	EventSkillStarted   = "skill.started"
	EventSkillCompleted = "skill.completed"
	EventSkillFailed    = "skill.failed"

	// Task events
	EventTaskCreated   = "task.created"
	EventTaskStarted   = "task.started"
	EventTaskCompleted = "task.completed"
	EventTaskFailed    = "task.failed"
)

// ConnectorEvent represents an event from a connector
type ConnectorEvent struct {
	*BaseEvent
	ConnectorName string
	UserID        string
	ChannelID     string
	Message       string
	Error         error
}

// NewConnectorEvent creates a new connector event
func NewConnectorEvent(eventType, connectorName, userID, channelID, message string, err error) *ConnectorEvent {
	return &ConnectorEvent{
		BaseEvent:     NewBaseEvent(eventType, nil),
		ConnectorName: connectorName,
		UserID:        userID,
		ChannelID:     channelID,
		Message:       message,
		Error:         err,
	}
}

// RouterEvent represents an event from the router
type RouterEvent struct {
	*BaseEvent
	MessageID string
	SessionID string
	UserID    string
	Content   string
	Source    string
	Error     error
}

// NewRouterEvent creates a new router event
func NewRouterEvent(eventType, messageID, sessionID, userID, content, source string, err error) *RouterEvent {
	return &RouterEvent{
		BaseEvent: NewBaseEvent(eventType, nil),
		MessageID: messageID,
		SessionID: sessionID,
		UserID:    userID,
		Content:   content,
		Source:    source,
		Error:     err,
	}
}

// LLMPublishedEvent represents an LLM event
type LLMPublishedEvent struct {
	*BaseEvent
	ProviderName string
	Model        string
	Tokens       int
	Cost         float64
	Duration     time.Duration
	Error        error
}

// NewLLMEvent creates a new LLM event
func NewLLMEvent(eventType, providerName, model string, tokens int, cost float64, duration time.Duration, err error) *LLMPublishedEvent {
	return &LLMPublishedEvent{
		BaseEvent:    NewBaseEvent(eventType, nil),
		ProviderName: providerName,
		Model:        model,
		Tokens:       tokens,
		Cost:         cost,
		Duration:     duration,
		Error:        err,
	}
}

// UserEvent represents a user-related event
type UserEvent struct {
	*BaseEvent
	UserID  string
	Email   string
	Channel string
}

// NewUserEvent creates a new user event
func NewUserEvent(eventType, userID, email, channel string) *UserEvent {
	return &UserEvent{
		BaseEvent: NewBaseEvent(eventType, nil),
		UserID:    userID,
		Email:     email,
		Channel:   channel,
	}
}

// SessionEvent represents a session-related event
type SessionEvent struct {
	*BaseEvent
	SessionID    string
	UserID       string
	MessageCount int
}

// NewSessionEvent creates a new session event
func NewSessionEvent(eventType, sessionID, userID string, messageCount int) *SessionEvent {
	return &SessionEvent{
		BaseEvent:    NewBaseEvent(eventType, nil),
		SessionID:    sessionID,
		UserID:       userID,
		MessageCount: messageCount,
	}
}

// SkillEvent represents a skill-related event
type SkillEvent struct {
	*BaseEvent
	SkillName string
	Input     string
	Output    string
	Error     error
	Duration  time.Duration
}

// NewSkillEvent creates a new skill event
func NewSkillEvent(eventType, skillName, input, output string, err error, duration time.Duration) *SkillEvent {
	return &SkillEvent{
		BaseEvent: NewBaseEvent(eventType, nil),
		SkillName: skillName,
		Input:     input,
		Output:    output,
		Error:     err,
		Duration:  duration,
	}
}

// TaskEvent represents a task-related event
type TaskEvent struct {
	*BaseEvent
	TaskID    string
	SessionID string
	SkillName string
	Status    string
	Input     string
	Output    string
	Error     string
}

// NewTaskEvent creates a new task event
func NewTaskEvent(eventType, taskID, sessionID, skillName, status, input, output, errMsg string) *TaskEvent {
	return &TaskEvent{
		BaseEvent: NewBaseEvent(eventType, nil),
		TaskID:    taskID,
		SessionID: sessionID,
		SkillName: skillName,
		Status:    status,
		Input:     input,
		Output:    output,
		Error:     errMsg,
	}
}

// EventLogger is a built-in event handler that logs events
type EventLogger struct {
	logger logging.Logger
}

// NewEventLogger creates a new event logger
//
// Parameters:
//   - logger: Logger to use
//
// Returns:
//   - *EventLogger: Event logger instance
func NewEventLogger(logger logging.Logger) *EventLogger {
	return &EventLogger{
		logger: logger,
	}
}

// Handle logs an event
//
// Parameters:
//   - ctx: Context for the operation
//   - event: Event to log
//
// Returns:
//   - error: Always nil
func (el *EventLogger) Handle(ctx context.Context, event Event) error {
	el.logger.Info("event received",
		"type", event.Type(),
		"timestamp", event.Timestamp(),
	)
	return nil
}
