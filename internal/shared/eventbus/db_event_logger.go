package eventbus

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// DatabaseEventLogger persists events to the database
type DatabaseEventLogger struct {
	logRepo repository.LogRepository
	logger  logging.Logger
}

// NewDatabaseEventLogger creates a new database event logger
//
// Parameters:
//   - logRepo: Repository for log persistence
//   - logger: Structured logger for logging
//
// Returns:
//   - *DatabaseEventLogger: Initialized database event logger
func NewDatabaseEventLogger(logRepo repository.LogRepository, logger logging.Logger) *DatabaseEventLogger {
	return &DatabaseEventLogger{
		logRepo: logRepo,
		logger:  logger,
	}
}

// Handle logs an event to the database
//
// Parameters:
//   - ctx: Context for the operation
//   - event: Event to log
//
// Returns:
//   - error: Error if persistence failed
func (del *DatabaseEventLogger) Handle(ctx context.Context, event Event) error {
	// Build metadata map
	metadata := map[string]interface{}{
		"timestamp": event.Timestamp().Format("2006-01-02T15:04:05Z07:00"),
	}

	// Add event-specific metadata based on event type
	switch e := event.(type) {
	case *ConnectorEvent:
		metadata["connector"] = e.ConnectorName
		metadata["user_id"] = e.UserID
		metadata["channel_id"] = e.ChannelID
		if e.Message != "" {
			metadata["message"] = e.Message
		}
		if e.Error != nil {
			metadata["error"] = e.Error.Error()
		}

	case *RouterEvent:
		metadata["message_id"] = e.MessageID
		metadata["session_id"] = e.SessionID
		metadata["user_id"] = e.UserID
		metadata["source"] = e.Source
		if e.Content != "" {
			metadata["content"] = e.Content
		}
		if e.Error != nil {
			metadata["error"] = e.Error.Error()
		}

	case *LLMPublishedEvent:
		metadata["provider"] = e.ProviderName
		metadata["model"] = e.Model
		metadata["tokens"] = e.Tokens
		metadata["cost"] = e.Cost
		metadata["duration_ms"] = e.Duration.Milliseconds()
		if e.Error != nil {
			metadata["error"] = e.Error.Error()
		}

	case *UserEvent:
		metadata["user_id"] = e.UserID
		metadata["email"] = e.Email
		metadata["channel"] = e.Channel

	case *SessionEvent:
		metadata["session_id"] = e.SessionID
		metadata["user_id"] = e.UserID
		metadata["message_count"] = e.MessageCount

	case *SkillEvent:
		metadata["skill"] = e.SkillName
		if e.Input != "" {
			metadata["input"] = e.Input
		}
		if e.Output != "" {
			metadata["output"] = e.Output
		}
		if e.Duration > 0 {
			metadata["duration_ms"] = e.Duration.Milliseconds()
		}
		if e.Error != nil {
			metadata["error"] = e.Error.Error()
		}

	case *TaskEvent:
		metadata["task_id"] = e.TaskID
		metadata["session_id"] = e.SessionID
		metadata["skill"] = e.SkillName
		metadata["status"] = e.Status
		if e.Input != "" {
			metadata["input"] = e.Input
		}
		if e.Output != "" {
			metadata["output"] = e.Output
		}
		if e.Error != "" {
			metadata["error"] = e.Error
		}
	}

	// Determine log level
	level := valueobject.LogLevelInfo
	switch event.Type() {
	case EventConnectorError, EventRouterError, EventOrchestratorError, EventLLMError, EventSkillFailed, EventTaskFailed:
		level = valueobject.LogLevelError
	}

	// Create log entity
	logEntry := entity.NewLog(level, "eventbus", event.Type(), metadata)

	// Persist log to database
	if err := del.logRepo.Create(ctx, logEntry); err != nil {
		del.logger.Error("failed to persist event log",
			"type", event.Type(),
			"error", err,
		)
		return err
	}

	return nil
}
