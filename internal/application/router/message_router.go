package router

import (
	"context"
	"fmt"
	"sync"

	"github.com/atumaikin/nexflow/internal/application/dto"
	"github.com/atumaikin/nexflow/internal/application/usecase"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/atumaikin/nexflow/internal/shared/eventbus"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// MessageRouter routes incoming messages from connectors to ChatUseCase
// and sends responses back through appropriate connector
type MessageRouter struct {
	connectors  map[string]channels.Connector
	chatUseCase *usecase.ChatUseCase
	eventBus    *eventbus.EventBus
	logger      logging.Logger
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
}

// NewMessageRouter creates a new MessageRouter instance
//
// Parameters:
//   - chatUseCase: ChatUseCase for processing messages
//   - eventBus: EventBus for publishing events
//   - logger: Structured logger for logging
//
// Returns:
//   - *MessageRouter: Initialized message router
func NewMessageRouter(chatUseCase *usecase.ChatUseCase, eventBus *eventbus.EventBus, logger logging.Logger) *MessageRouter {
	ctx, cancel := context.WithCancel(context.Background())

	return &MessageRouter{
		connectors:  make(map[string]channels.Connector),
		chatUseCase: chatUseCase,
		eventBus:    eventBus,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
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

	r.logger.Info("connector registered", "connector", name)
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
		r.logger.Info("connector unregistered", "connector", name)
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

			r.logger.Info("received message", "connector", connectorName, "user_id", msg.UserID)

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
			go r.handleMessage(connectorName, conn, msg)
		}
	}
}

// handleMessage handles a single message from a connector
func (r *MessageRouter) handleMessage(connectorName string, conn channels.Connector, msg *channels.Message) {
	ctx := r.ctx

	// Check if chatUseCase is available (nil check for testing)
	if r.chatUseCase == nil {
		r.logger.Warn("chat use case not available, skipping message processing", "connector", connectorName, "user_id", msg.UserID)
		return
	}

	// Get or create user for the channel user ID
	user, err := conn.GetUser(ctx, msg.UserID)
	if err != nil {
		r.logger.Error("failed to get user, creating new user", "connector", connectorName, "user_id", msg.UserID, "error", err)
		user, err = conn.CreateUser(ctx, msg.UserID)
		if err != nil {
			r.logger.Error("failed to create user", "connector", connectorName, "user_id", msg.UserID, "error", err)
			// Send error message back to user
			r.sendErrorResponse(ctx, conn, msg.UserID, "Sorry, I encountered an error processing your request.")
			return
		}
	}

	// Create SendMessageRequest
	req := dto.SendMessageRequest{
		UserID: string(user.ID),
		Message: dto.ChatMessage{
			Role:    "user",
			Content: msg.Content,
		},
		Options: dto.MessageOptions{
			MaxTokens: 1000,
		},
	}

	// Process message through ChatUseCase
	resp, err := r.chatUseCase.SendMessage(ctx, req)
	if err != nil {
		r.logger.Error("failed to process message", "connector", connectorName, "user_id", msg.UserID, "error", err)
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

			if resp.Message.SessionID != "" {
				response.Metadata["session_id"] = resp.Message.SessionID
			}

			if err := conn.SendResponse(ctx, msg.UserID, response); err != nil {
				r.logger.Error("failed to send response", "connector", connectorName, "user_id", msg.UserID, "error", err)

				// Publish router error event
				if r.eventBus != nil {
					event := eventbus.NewRouterEvent(
						eventbus.EventRouterError,
						resp.Message.ID,
						resp.Message.SessionID,
						msg.UserID,
						resp.Message.Content,
						connectorName,
						err,
					)
					r.eventBus.Publish(event)
				}
			} else {
				r.logger.Info("response sent", "connector", connectorName, "user_id", msg.UserID)

				// Publish router message event
				if r.eventBus != nil {
					event := eventbus.NewRouterEvent(
						eventbus.EventRouterMessage,
						resp.Message.ID,
						resp.Message.SessionID,
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

// sendErrorResponse sends an error message to the user through the connector
func (r *MessageRouter) sendErrorResponse(ctx context.Context, conn channels.Connector, userID string, message string) {
	response := &channels.Response{
		Content: message,
		Metadata: map[string]interface{}{
			"error": true,
		},
	}

	if err := conn.SendResponse(ctx, userID, response); err != nil {
		r.logger.Error("failed to send error response", "user_id", userID, "error", err)
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
