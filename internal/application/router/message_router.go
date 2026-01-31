package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/ports"
	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/repository"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/eventbus"
	"github.com/atumaikin/nexflow/internal/shared/logging"
	"github.com/atumaikin/nexflow/internal/shared/metrics"
)

// RouterMetrics holds all metrics for the MessageRouter
type RouterMetrics struct {
	MessagesReceived          *metrics.Counter
	MessagesProcessed         *metrics.Counter
	MessagesFailed            *metrics.Counter
	MessagesValidated         *metrics.Counter
	MessageValidationFailed   *metrics.Counter
	MessageProcessingDuration *metrics.Histogram
	ResponseSentDuration      *metrics.Histogram
	ConnectorsActive          *metrics.Counter
}

// NewRouterMetrics creates a new RouterMetrics instance
func NewRouterMetrics() *RouterMetrics {
	registry := metrics.NewMetricsRegistry()
	buckets := metrics.DefaultBuckets()

	return &RouterMetrics{
		MessagesReceived:          registry.GetCounter("router_messages_received_total"),
		MessagesProcessed:         registry.GetCounter("router_messages_processed_total"),
		MessagesFailed:            registry.GetCounter("router_messages_failed_total"),
		MessagesValidated:         registry.GetCounter("router_messages_validated_total"),
		MessageValidationFailed:   registry.GetCounter("router_message_validation_failed_total"),
		MessageProcessingDuration: registry.GetHistogram("router_message_processing_duration_seconds", buckets),
		ResponseSentDuration:      registry.GetHistogram("router_response_sent_duration_seconds", buckets),
		ConnectorsActive:          registry.GetCounter("router_connectors_active"),
	}
}

// MessageRouter routes incoming messages from connectors to Orchestrator
// and sends responses back through appropriate connector
type MessageRouter struct {
	connectors    map[string]channels.Connector
	sessionRepo   repository.SessionRepository
	orchestrator  ports.Orchestrator
	eventBus      *eventbus.EventBus
	logger        logging.Logger
	config        *Config
	validator     *MessageValidator
	retryHandler  *RetryHandler
	routerMetrics *RouterMetrics
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
}

// NewMessageRouter creates a new MessageRouter instance
//
// Parameters:
//   - sessionRepo: SessionRepository for managing sessions
//   - orchestrator: Orchestrator for processing messages
//   - eventBus: EventBus for publishing events
//   - logger: Structured logger for logging
//   - config: Router configuration (uses defaults if nil)
//
// Returns:
//   - *MessageRouter: Initialized message router
func NewMessageRouter(sessionRepo repository.SessionRepository, orchestrator ports.Orchestrator, eventBus *eventbus.EventBus, logger logging.Logger, config *Config) *MessageRouter {
	return NewMessageRouterWithMetrics(sessionRepo, orchestrator, eventBus, logger, config, nil)
}

// NewMessageRouterWithMetrics creates a new MessageRouter with custom metrics
func NewMessageRouterWithMetrics(sessionRepo repository.SessionRepository, orchestrator ports.Orchestrator, eventBus *eventbus.EventBus, logger logging.Logger, config *Config, routerMetrics *RouterMetrics) *MessageRouter {
	ctx, cancel := context.WithCancel(context.Background())

	// Use default config if not provided
	if config == nil {
		config = DefaultConfig()
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		logger.Error("invalid router configuration, using defaults", "error", err)
		config = DefaultConfig()
	}

	// Initialize validator and retry handler
	validator := NewMessageValidator(config, logger)
	retryHandler := NewRetryHandler(config.RetryConfig, logger)

	// Create default metrics if not provided
	if routerMetrics == nil {
		routerMetrics = NewRouterMetrics()
	}

	return &MessageRouter{
		connectors:    make(map[string]channels.Connector),
		sessionRepo:   sessionRepo,
		orchestrator:  orchestrator,
		eventBus:      eventBus,
		logger:        logger,
		config:        config,
		validator:     validator,
		retryHandler:  retryHandler,
		routerMetrics: routerMetrics,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// RegisterConnector registers a connector with the router
//
// Parameters:
//   - connector: Connector to register
func (r *MessageRouter) RegisterConnector(connector channels.Connector) {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := connector.Name()
	r.connectors[name] = connector
	r.routerMetrics.ConnectorsActive.Inc()

	r.logger.Info("connector registered", "connector", name, "total_connectors", len(r.connectors))
}

// UnregisterConnector removes a connector from the router
//
// Parameters:
//   - name: Name of the connector to unregister
func (r *MessageRouter) UnregisterConnector(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if conn, exists := r.connectors[name]; exists {
		if err := conn.Stop(r.ctx); err != nil {
			r.logger.Error("failed to stop connector", "connector", name, "error", err)
		}
		delete(r.connectors, name)
		r.routerMetrics.ConnectorsActive.Add(-1) // Decrement counter
		r.logger.Info("connector unregistered",
			"connector", name,
			"total_connectors", len(r.connectors),
		)
	}
}

// Start starts the message router and all registered connectors
//
// Returns:
//   - error: Error if starting failed
func (r *MessageRouter) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.logger.Info("starting message router")

	// Start all connectors
	for name, conn := range r.connectors {
		if err := conn.Start(r.ctx); err != nil {
			return fmt.Errorf("failed to start connector %s: %w", name, err)
		}
		r.logger.Info("connector started", "connector", name)
	}

	// Start message processing for each connector
	for name, conn := range r.connectors {
		r.wg.Add(1)
		go r.processMessages(name, conn)
	}

	r.logger.Info("message router started")

	// Publish router started event
	if r.eventBus != nil {
		event := eventbus.NewRouterEvent(
			eventbus.EventRouterStarted,
			"", "", "", "message router started", "router", nil,
		)
		r.eventBus.Publish(event)
	}

	return nil
}

// Stop stops the message router and all connectors gracefully
func (r *MessageRouter) Stop() error {
	r.logger.Info("stopping message router")

	// Cancel context to signal all goroutines to stop
	r.cancel()

	// Wait for all goroutines to finish
	r.wg.Wait()

	r.mu.Lock()
	defer r.mu.Unlock()

	// Stop all connectors
	for name, conn := range r.connectors {
		if err := conn.Stop(r.ctx); err != nil {
			r.logger.Error("failed to stop connector", "connector", name, "error", err)
		}
	}

	r.logger.Info("message router stopped")

	// Publish router stopped event
	if r.eventBus != nil {
		event := eventbus.NewRouterEvent(
			eventbus.EventRouterStopped,
			"", "", "", "message router stopped", "router", nil,
		)
		r.eventBus.Publish(event)
	}

	return nil
}

// processMessages processes incoming messages from a connector
func (r *MessageRouter) processMessages(connectorName string, conn channels.Connector) {
	defer r.wg.Done()

	r.logger.Info("started processing messages", "connector", connectorName)

	for {
		select {
		case <-r.ctx.Done():
			r.logger.Info("stopped processing messages", "connector", connectorName)
			return

		case msg, ok := <-conn.Incoming():
			if !ok {
				r.logger.Info("connector channel closed", "connector", connectorName)
				return
			}

			// Record message received
			r.routerMetrics.MessagesReceived.Inc()

			// Validate message before processing
			if err := r.validator.Validate(msg); err != nil {
				r.routerMetrics.MessageValidationFailed.Inc()

				r.logger.Error("message validation failed",
					"connector", connectorName,
					"user_id", msg.UserID,
					"message_length", len(msg.Content),
					"error", err,
				)

				// Publish validation error event
				if r.eventBus != nil {
					event := eventbus.NewRouterEvent(
						eventbus.EventRouterError,
						"",
						"",
						msg.UserID,
						msg.Content,
						connectorName,
						err,
					)
					r.eventBus.Publish(event)
				}

				// Send error response to user
				r.sendErrorResponse(r.ctx, conn, msg.UserID,
					"Sorry, your message could not be processed. Please check the format and try again.")

				continue
			}

			r.routerMetrics.MessagesValidated.Inc()

			r.logger.Info("received message",
				"connector", connectorName,
				"user_id", msg.UserID,
				"message_length", len(msg.Content),
			)

			// Publish connector message event
			if r.eventBus != nil {
				event := eventbus.NewConnectorEvent(
					eventbus.EventConnectorMessage,
					connectorName,
					msg.UserID,
					msg.ChannelID,
					msg.Content,
					nil,
				)
				r.eventBus.Publish(event)
			}

			// Process the message in a separate goroutine to avoid blocking
			go func() {
				// Wrap handling with metrics
				err := metrics.RecordDurationWithError(r.routerMetrics.MessageProcessingDuration, func() error {
					r.handleMessage(connectorName, conn, msg)
					return nil
				})

				// Record processing metrics
				if err != nil {
					r.routerMetrics.MessagesFailed.Inc()
					r.logger.Error("message processing failed",
						"connector", connectorName,
						"user_id", msg.UserID,
						"error", err,
					)
				} else {
					r.routerMetrics.MessagesProcessed.Inc()
				}
			}()
		}
	}
}

// handleMessage handles a single message from a connector
func (r *MessageRouter) handleMessage(connectorName string, conn channels.Connector, msg *channels.Message) {
	ctx := r.ctx
	var err error

	// Check if orchestrator is available (nil check for testing)
	if r.orchestrator == nil {
		r.logger.Warn("orchestrator not available, skipping message processing", "connector", connectorName, "user_id", msg.UserID)
		return
	}

	// Get or create user with retry
	var user *entity.User
	err = r.retryHandler.Do(ctx, "get_or_create_user", func() error {
		var err error
		user, err = conn.GetUser(ctx, msg.UserID)
		if err != nil {
			r.logger.Info("user not found, creating new user", "connector", connectorName, "user_id", msg.UserID, "error", err)
			user, err = conn.CreateUser(ctx, msg.UserID)
		}
		return err
	})

	if err != nil {
		r.logger.Error("failed to get or create user",
			"connector", connectorName,
			"user_id", msg.UserID,
			"error", err,
		)
		r.sendErrorResponse(ctx, conn, msg.UserID, "Sorry, I encountered an error processing your request.")
		return
	}

	// Get or create session with retry
	var session *entity.Session
	err = r.retryHandler.Do(ctx, "get_or_create_session", func() error {
		sessions, err := r.sessionRepo.FindByUserID(ctx, string(user.ID))
		if err != nil || len(sessions) == 0 {
			// No session exists, create a new one
			newSession := entity.NewSession(string(user.ID))
			if err := r.sessionRepo.Create(ctx, newSession); err != nil {
				return fmt.Errorf("failed to create session: %w", err)
			}
			session = newSession
			r.logger.Info("session created",
				"connector", connectorName,
				"session_id", session.ID,
				"user_id", user.ID,
			)
			return nil
		}

		// Use the most recent session (last in the list)
		session = sessions[len(sessions)-1]
		return nil
	})

	if err != nil {
		r.logger.Error("failed to get or create session",
			"connector", connectorName,
			"user_id", msg.UserID,
			"error", err,
		)
		r.sendErrorResponse(ctx, conn, msg.UserID, "Sorry, I encountered an error creating a session.")
		return
	}

	// Prepare message options with session ID
	options := dto.MessageOptions{
		MaxTokens: 1000,
	}

	// Process message through Orchestrator
	resp, err := r.orchestrator.ProcessMessage(ctx, string(user.ID), msg.Content, options)
	if err != nil {
		r.logger.Error("failed to process message",
			"connector", connectorName,
			"user_id", msg.UserID,
			"session_id", session.ID,
			"error", err,
		)
		r.sendErrorResponse(ctx, conn, msg.UserID, "Sorry, I encountered an error generating a response.")
		return
	}

	// Send response back through connector
	if resp != nil && resp.Message != nil {
		// Check if the response contains an assistant message
		if resp.Message.Role == "assistant" {
			response := &channels.Response{
				Content: resp.Message.Content,
				Metadata: map[string]interface{}{
					"message_id": resp.Message.ID,
				},
			}

			// Add session ID to response metadata
			if session.ID != "" {
				response.Metadata["session_id"] = session.ID.String()
			}

			if err := conn.SendResponse(ctx, msg.UserID, response); err != nil {
				r.logger.Error("failed to send response",
					"connector", connectorName,
					"user_id", msg.UserID,
					"session_id", session.ID,
					"message_id", resp.Message.ID,
					"error", err,
				)

				// Publish router error event
				if r.eventBus != nil {
					event := eventbus.NewRouterEvent(
						eventbus.EventRouterError,
						resp.Message.ID,
						session.ID.String(),
						msg.UserID,
						resp.Message.Content,
						connectorName,
						err,
					)
					r.eventBus.Publish(event)
				}
			} else {
				r.logger.Info("response sent",
					"connector", connectorName,
					"user_id", msg.UserID,
					"session_id", session.ID,
					"message_id", resp.Message.ID,
				)

				// Publish router message event
				if r.eventBus != nil {
					event := eventbus.NewRouterEvent(
						eventbus.EventRouterMessage,
						resp.Message.ID,
						session.ID.String(),
						msg.UserID,
						resp.Message.Content,
						connectorName,
						nil,
					)
					r.eventBus.Publish(event)
				}
			}
		}
	}
}

// sendErrorResponse sends an error message to user through connector
func (r *MessageRouter) sendErrorResponse(ctx context.Context, conn channels.Connector, userID string, message string) {
	response := &channels.Response{
		Content: message,
		Metadata: map[string]interface{}{
			"error": true,
		},
	}

	err := conn.SendResponse(ctx, userID, response)
	if err != nil {
		r.logger.Error("failed to send error response",
			"user_id", userID,
			"message_length", len(message),
			"error", err,
		)
		r.routerMetrics.MessagesFailed.Inc()
	} else {
		r.routerMetrics.MessagesProcessed.Inc()
	}
}

// GetConnector returns a connector by name
//
// Parameters:
//   - name: Name of the connector
//
// Returns:
//   - channels.Connector: Connector instance
//   - bool: True if connector exists
func (r *MessageRouter) GetConnector(name string) (channels.Connector, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	conn, exists := r.connectors[name]
	return conn, exists
}

// ListConnectors returns a list of registered connector names
//
// Returns:
//   - []string: List of connector names
func (r *MessageRouter) ListConnectors() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.connectors))
	for name := range r.connectors {
		names = append(names, name)
	}
	return names
}
