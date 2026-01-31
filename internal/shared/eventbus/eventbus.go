package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// Event represents a generic event that can be published through the event bus
type Event interface {
	// Type returns the event type identifier
	Type() string
	// Timestamp returns when the event was created
	Timestamp() time.Time
}

// EventHandler is a function that handles an event
type EventHandler func(ctx context.Context, event Event) error

// EventSubscription represents a subscription to event types
type EventSubscription struct {
	ID      string
	Types   []string
	Handler EventHandler
	cancel  context.CancelFunc
}

// EventBus implements a publish-subscribe pattern for internal events
type EventBus struct {
	mu            sync.RWMutex
	subscriptions map[string][]*EventSubscription
	handlers      map[string][]EventHandler
	logger        logging.Logger
	eventChannel  chan Event
	subIDCounter  uint64
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	batchSize     int
	flushInterval time.Duration
	eventBuffer   []Event
	bufferMu      sync.Mutex
	started       bool
}

// EventBusConfig contains configuration options for the EventBus
type EventBusConfig struct {
	// BatchSize is the number of events to batch before processing
	BatchSize int
	// FlushInterval is the maximum time to wait before flushing events
	FlushInterval time.Duration
	// Logger is the logger to use
	Logger logging.Logger
}

// DefaultConfig returns the default configuration for EventBus
func DefaultConfig() *EventBusConfig {
	return &EventBusConfig{
		BatchSize:     100,
		FlushInterval: 100 * time.Millisecond,
		Logger:        logging.NewNoopLogger(),
	}
}

// NewEventBus creates a new EventBus with the given configuration
//
// Parameters:
//   - config: Configuration options for the event bus (uses defaults if nil)
//
// Returns:
//   - *EventBus: Initialized event bus
func NewEventBus(config *EventBusConfig) *EventBus {
	if config == nil {
		config = DefaultConfig()
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &EventBus{
		subscriptions: make(map[string][]*EventSubscription),
		handlers:      make(map[string][]EventHandler),
		logger:        config.Logger,
		eventChannel:  make(chan Event, 1000),
		ctx:           ctx,
		cancel:        cancel,
		batchSize:     config.BatchSize,
		flushInterval: config.FlushInterval,
		eventBuffer:   make([]Event, 0, config.BatchSize),
	}
}

// Start starts the event bus and begins processing events
//
// Returns:
//   - error: Error if starting failed
func (eb *EventBus) Start() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if eb.started {
		return nil
	}

	// Start event processor
	eb.wg.Add(1)
	go eb.processEvents()

	// Start batch flusher
	eb.wg.Add(1)
	go eb.flushBatches()

	eb.started = true
	eb.logger.Info("event bus started")
	return nil
}

// Stop stops the event bus gracefully
func (eb *EventBus) Stop() error {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if !eb.started {
		return nil
	}

	// Cancel context to signal shutdown
	eb.cancel()

	// Close event channel
	close(eb.eventChannel)

	// Wait for all goroutines to finish
	eb.wg.Wait()

	// Flush remaining events
	eb.flushBuffer()

	eb.started = false
	eb.logger.Info("event bus stopped")
	return nil
}

// Subscribe subscribes to events of specific types
//
// Parameters:
//   - eventTypes: List of event types to subscribe to
//   - handler: Function to handle events
//
// Returns:
//   - *EventSubscription: Subscription that can be used to unsubscribe
func (eb *EventBus) Subscribe(eventTypes []string, handler EventHandler) *EventSubscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	subID := eb.generateSubID()
	sub := &EventSubscription{
		ID:      subID,
		Types:   eventTypes,
		Handler: handler,
	}

	// Add subscription for each event type
	for _, eventType := range eventTypes {
		eb.subscriptions[eventType] = append(eb.subscriptions[eventType], sub)
	}

	eb.logger.Info("subscription created", "id", subID, "types", eventTypes)
	return sub
}

// Unsubscribe removes a subscription from the event bus
//
// Parameters:
//   - sub: Subscription to remove
func (eb *EventBus) Unsubscribe(sub *EventSubscription) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	if sub == nil {
		return
	}

	// Remove subscription from all event types
	for _, eventType := range sub.Types {
		subs := eb.subscriptions[eventType]
		for i, s := range subs {
			if s.ID == sub.ID {
				eb.subscriptions[eventType] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	}

	eb.logger.Info("subscription removed", "id", sub.ID)
}

// Publish publishes an event to the event bus
//
// Parameters:
//   - event: Event to publish
func (eb *EventBus) Publish(event Event) {
	if event == nil {
		eb.logger.Warn("attempted to publish nil event")
		return
	}

	eb.logger.Debug("event published", "type", event.Type())

	// Send event to channel for processing
	select {
	case eb.eventChannel <- event:
		// Event queued successfully
	default:
		eb.logger.Warn("event channel full, dropping event", "type", event.Type())
	}
}

// PublishAsync publishes an event asynchronously
// This is a convenience method that wraps Publish in a goroutine
//
// Parameters:
//   - event: Event to publish
func (eb *EventBus) PublishAsync(event Event) {
	go eb.Publish(event)
}

// SubscribeHandler subscribes a handler to specific event types (convenience method)
//
// Parameters:
//   - eventTypes: List of event types to subscribe to
//   - handler: Handler function
//
// Returns:
//   - *EventSubscription: Subscription that can be used to unsubscribe
func (eb *EventBus) SubscribeHandler(eventTypes []string, handler EventHandler) *EventSubscription {
	return eb.Subscribe(eventTypes, handler)
}

// processEvents processes events from the event channel
func (eb *EventBus) processEvents() {
	defer eb.wg.Done()

	for {
		select {
		case <-eb.ctx.Done():
			return

		case event, ok := <-eb.eventChannel:
			if !ok {
				return
			}

			// Add event to buffer
			eb.bufferMu.Lock()
			eb.eventBuffer = append(eb.eventBuffer, event)

			// Flush if buffer is full
			if len(eb.eventBuffer) >= eb.batchSize {
				eb.flushBuffer()
			}
			eb.bufferMu.Unlock()
		}
	}
}

// flushBatches periodically flushes buffered events
func (eb *EventBus) flushBatches() {
	defer eb.wg.Done()

	ticker := time.NewTicker(eb.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-eb.ctx.Done():
			return
		case <-ticker.C:
			eb.flushBuffer()
		}
	}
}

// flushBuffer processes all buffered events
func (eb *EventBus) flushBuffer() {
	eb.bufferMu.Lock()
	events := eb.eventBuffer
	eb.eventBuffer = make([]Event, 0, eb.batchSize)
	eb.bufferMu.Unlock()

	if len(events) == 0 {
		return
	}

	eb.logger.Debug("flushing events", "count", len(events))

	// Process each event
	for _, event := range events {
		eb.dispatch(event)
	}
}

// dispatch dispatches an event to all subscribed handlers
func (eb *EventBus) dispatch(event Event) {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	// Get all subscriptions for this event type
	subs := eb.subscriptions[event.Type()]

	// Create a copy of handlers to avoid holding the lock during execution
	handlers := make([]EventHandler, len(subs))
	for i, sub := range subs {
		handlers[i] = sub.Handler
	}

	// Execute each handler
	for _, handler := range handlers {
		go func(h EventHandler) {
			ctx, cancel := context.WithTimeout(eb.ctx, 30*time.Second)
			defer cancel()

			if err := h(ctx, event); err != nil {
				eb.logger.Error("event handler failed", "type", event.Type(), "error", err)
			}
		}(handler)
	}
}

// generateSubID generates a unique subscription ID
func (eb *EventBus) generateSubID() string {
	eb.subIDCounter++
	return string(rune(eb.subIDCounter))
}

// AddHandler adds a handler for a specific event type (convenience method)
//
// Parameters:
//   - eventType: Type of event to handle
//   - handler: Handler function
func (eb *EventBus) AddHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.logger.Debug("handler added", "type", eventType)
}

// RemoveHandler removes a handler for a specific event type
//
// Parameters:
//   - eventType: Type of event
//   - handler: Handler function to remove
func (eb *EventBus) RemoveHandler(eventType string, handler EventHandler) {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	handlers := eb.handlers[eventType]
	for i, h := range handlers {
		// Compare function pointers
		if &h == &handler {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			eb.logger.Debug("handler removed", "type", eventType)
			return
		}
	}
}

// GetSubscriptionCount returns the number of subscriptions for a specific event type
//
// Parameters:
//   - eventType: Type of event
//
// Returns:
//   - int: Number of subscriptions
func (eb *EventBus) GetSubscriptionCount(eventType string) int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	return len(eb.subscriptions[eventType])
}

// GetAllSubscriptionCounts returns a map of event types to subscription counts
//
// Returns:
//   - map[string]int: Map of event type to subscription count
func (eb *EventBus) GetAllSubscriptionCounts() map[string]int {
	eb.mu.RLock()
	defer eb.mu.RUnlock()

	counts := make(map[string]int)
	for eventType, subs := range eb.subscriptions {
		counts[eventType] = len(subs)
	}
	return counts
}
