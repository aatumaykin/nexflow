package eventbus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/shared/logging"
)

func TestNewEventBus(t *testing.T) {
	config := &EventBusConfig{
		BatchSize:     50,
		FlushInterval: 50 * time.Millisecond,
		Logger:        logging.NewNoopLogger(),
	}

	eb := NewEventBus(config)

	if eb == nil {
		t.Fatal("Expected non-nil event bus")
	}

	if eb.batchSize != 50 {
		t.Errorf("Expected batch size 50, got %d", eb.batchSize)
	}

	if eb.flushInterval != 50*time.Millisecond {
		t.Errorf("Expected flush interval 50ms, got %v", eb.flushInterval)
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("Expected non-nil config")
	}

	if config.BatchSize != 100 {
		t.Errorf("Expected default batch size 100, got %d", config.BatchSize)
	}

	if config.FlushInterval != 100*time.Millisecond {
		t.Errorf("Expected default flush interval 100ms, got %v", config.FlushInterval)
	}
}

func TestEventBusStartStop(t *testing.T) {
	eb := NewEventBus(DefaultConfig())

	if err := eb.Start(); err != nil {
		t.Fatalf("Failed to start event bus: %v", err)
	}

	if err := eb.Stop(); err != nil {
		t.Fatalf("Failed to stop event bus: %v", err)
	}
}

func TestEventBusPublish(t *testing.T) {
	eb := NewEventBus(DefaultConfig())
	eb.Start()
	defer eb.Stop()

	receivedEvents := make([]Event, 0)
	var mu sync.Mutex

	// Subscribe to events
	eb.Subscribe([]string{"test.event"}, func(ctx context.Context, event Event) error {
		mu.Lock()
		defer mu.Unlock()
		receivedEvents = append(receivedEvents, event)
		return nil
	})

	// Publish events
	event1 := NewBaseEvent("test.event", nil)
	event2 := NewBaseEvent("test.event", nil)

	eb.Publish(event1)
	eb.Publish(event2)

	// Wait for events to be processed
	time.Sleep(200 * time.Millisecond)

	// Check if events were received
	mu.Lock()
	count := len(receivedEvents)
	mu.Unlock()

	if count < 2 {
		t.Errorf("Expected at least 2 events, got %d", count)
	}
}

func TestEventBusSubscribe(t *testing.T) {
	eb := NewEventBus(DefaultConfig())

	receivedCount := 0
	var wg sync.WaitGroup

	handler := func(ctx context.Context, event Event) error {
		receivedCount++
		wg.Done()
		return nil
	}

	sub := eb.Subscribe([]string{"test.event"}, handler)

	if sub == nil {
		t.Fatal("Expected non-nil subscription")
	}

	if sub.ID == "" {
		t.Error("Expected non-empty subscription ID")
	}

	if len(sub.Types) != 1 {
		t.Errorf("Expected 1 event type, got %d", len(sub.Types))
	}

	eb.Unsubscribe(sub)
}

func TestEventBusMultipleSubscriptions(t *testing.T) {
	eb := NewEventBus(DefaultConfig())
	eb.Start()
	defer eb.Stop()

	var mu sync.Mutex
	handlersCalled := make(map[string]int)

	// Create multiple subscriptions
	eb.Subscribe([]string{"event.a"}, func(ctx context.Context, event Event) error {
		mu.Lock()
		handlersCalled["a"]++
		mu.Unlock()
		return nil
	})

	eb.Subscribe([]string{"event.a"}, func(ctx context.Context, event Event) error {
		mu.Lock()
		handlersCalled["a2"]++
		mu.Unlock()
		return nil
	})

	eb.Subscribe([]string{"event.b"}, func(ctx context.Context, event Event) error {
		mu.Lock()
		handlersCalled["b"]++
		mu.Unlock()
		return nil
	})

	// Publish events
	eb.Publish(NewBaseEvent("event.a", nil))
	eb.Publish(NewBaseEvent("event.b", nil))

	// Wait for processing
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	aCount := handlersCalled["a"]
	a2Count := handlersCalled["a2"]
	bCount := handlersCalled["b"]
	mu.Unlock()

	if aCount != 1 {
		t.Errorf("Expected handler 'a' to be called once, got %d", aCount)
	}

	if a2Count != 1 {
		t.Errorf("Expected handler 'a2' to be called once, got %d", a2Count)
	}

	if bCount != 1 {
		t.Errorf("Expected handler 'b' to be called once, got %d", bCount)
	}
}

func TestEventBusUnsubscribe(t *testing.T) {
	eb := NewEventBus(DefaultConfig())
	eb.Start()
	defer eb.Stop()

	receivedCount := 0
	var mu sync.Mutex

	handler := func(ctx context.Context, event Event) error {
		mu.Lock()
		receivedCount++
		mu.Unlock()
		return nil
	}

	sub := eb.Subscribe([]string{"test.event"}, handler)

	// Publish first event
	eb.Publish(NewBaseEvent("test.event", nil))
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count1 := receivedCount
	mu.Unlock()

	// Unsubscribe
	eb.Unsubscribe(sub)

	// Publish second event
	eb.Publish(NewBaseEvent("test.event", nil))
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count2 := receivedCount
	mu.Unlock()

	if count1 < 1 {
		t.Errorf("Expected at least 1 event before unsubscribe, got %d", count1)
	}

	if count2 != count1 {
		t.Errorf("Expected no events after unsubscribe, got %d total (started with %d)", count2, count1)
	}
}

func TestEventBusBatchProcessing(t *testing.T) {
	config := &EventBusConfig{
		BatchSize:     10,
		FlushInterval: 500 * time.Millisecond,
		Logger:        logging.NewNoopLogger(),
	}

	eb := NewEventBus(config)
	eb.Start()
	defer eb.Stop()

	receivedCount := 0
	var mu sync.Mutex
	var wg sync.WaitGroup

	handler := func(ctx context.Context, event Event) error {
		mu.Lock()
		receivedCount++
		mu.Unlock()
		wg.Done()
		return nil
	}

	eb.Subscribe([]string{"batch.test"}, handler)

	// Publish batch size events
	wg.Add(10)
	for i := 0; i < 10; i++ {
		eb.Publish(NewBaseEvent("batch.test", nil))
	}

	// Wait for batch to flush
	wg.Wait()

	mu.Lock()
	count := receivedCount
	mu.Unlock()

	if count != 10 {
		t.Errorf("Expected 10 events to be processed, got %d", count)
	}
}

func TestEventBusAsyncPublish(t *testing.T) {
	eb := NewEventBus(DefaultConfig())
	eb.Start()
	defer eb.Stop()

	receivedCount := 0
	var mu sync.Mutex

	handler := func(ctx context.Context, event Event) error {
		mu.Lock()
		receivedCount++
		mu.Unlock()
		return nil
	}

	eb.Subscribe([]string{"async.test"}, handler)

	// Publish asynchronously
	eb.PublishAsync(NewBaseEvent("async.test", nil))
	eb.PublishAsync(NewBaseEvent("async.test", nil))

	// Wait for events to be processed
	time.Sleep(200 * time.Millisecond)

	mu.Lock()
	count := receivedCount
	mu.Unlock()

	if count < 2 {
		t.Errorf("Expected at least 2 events, got %d", count)
	}
}

func TestBaseEvent(t *testing.T) {
	event := NewBaseEvent("test.type", "test data")

	if event.Type() != "test.type" {
		t.Errorf("Expected type 'test.type', got '%s'", event.Type())
	}

	if event.Data() != "test data" {
		t.Errorf("Expected data 'test data', got '%v'", event.Data())
	}

	if event.Timestamp().IsZero() {
		t.Error("Expected non-zero timestamp")
	}

	// Test metadata
	event.SetMetadataValue("key1", "value1")
	val, ok := event.GetMetadataValue("key1")
	if !ok {
		t.Error("Expected metadata value to exist")
	}

	if val != "value1" {
		t.Errorf("Expected metadata value 'value1', got '%v'", val)
	}
}

func TestEventLogger(t *testing.T) {
	logger := logging.NewNoopLogger()
	eventLogger := NewEventLogger(logger)

	if eventLogger == nil {
		t.Fatal("Expected non-nil event logger")
	}

	event := NewBaseEvent("test.event", nil)
	err := eventLogger.Handle(context.Background(), event)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestEventBusGetSubscriptionCount(t *testing.T) {
	eb := NewEventBus(DefaultConfig())

	eb.Subscribe([]string{"event.a"}, func(ctx context.Context, event Event) error {
		return nil
	})

	eb.Subscribe([]string{"event.a"}, func(ctx context.Context, event Event) error {
		return nil
	})

	eb.Subscribe([]string{"event.b"}, func(ctx context.Context, event Event) error {
		return nil
	})

	countA := eb.GetSubscriptionCount("event.a")
	countB := eb.GetSubscriptionCount("event.b")
	countC := eb.GetSubscriptionCount("event.c")

	if countA != 2 {
		t.Errorf("Expected 2 subscriptions for event.a, got %d", countA)
	}

	if countB != 1 {
		t.Errorf("Expected 1 subscription for event.b, got %d", countB)
	}

	if countC != 0 {
		t.Errorf("Expected 0 subscriptions for event.c, got %d", countC)
	}
}

func TestEventBusGetAllSubscriptionCounts(t *testing.T) {
	eb := NewEventBus(DefaultConfig())

	eb.Subscribe([]string{"event.a"}, func(ctx context.Context, event Event) error {
		return nil
	})

	eb.Subscribe([]string{"event.a", "event.b"}, func(ctx context.Context, event Event) error {
		return nil
	})

	counts := eb.GetAllSubscriptionCounts()

	if len(counts) != 2 {
		t.Errorf("Expected 2 event types, got %d", len(counts))
	}

	if counts["event.a"] != 2 {
		t.Errorf("Expected 2 subscriptions for event.a, got %d", counts["event.a"])
	}

	if counts["event.b"] != 1 {
		t.Errorf("Expected 1 subscription for event.b, got %d", counts["event.b"])
	}
}

func TestSpecificEventTypes(t *testing.T) {
	// Test connector event
	connEvent := NewConnectorEvent(EventConnectorStarted, "telegram", "user123", "chat456", "", nil)
	if connEvent.Type() != EventConnectorStarted {
		t.Errorf("Expected type %s, got %s", EventConnectorStarted, connEvent.Type())
	}

	// Test router event
	routerEvent := NewRouterEvent(EventRouterMessage, "msg123", "sess456", "user789", "hello", "telegram", nil)
	if routerEvent.Type() != EventRouterMessage {
		t.Errorf("Expected type %s, got %s", EventRouterMessage, routerEvent.Type())
	}

	// Test LLM event
	llmEvent := NewLLMEvent(EventLLMRequest, "openai", "gpt-4", 100, 0.05, 0, nil)
	if llmEvent.Type() != EventLLMRequest {
		t.Errorf("Expected type %s, got %s", EventLLMRequest, llmEvent.Type())
	}

	// Test user event
	userEvent := NewUserEvent(EventUserCreated, "user123", "test@example.com", "telegram")
	if userEvent.Type() != EventUserCreated {
		t.Errorf("Expected type %s, got %s", EventUserCreated, userEvent.Type())
	}

	// Test session event
	sessionEvent := NewSessionEvent(EventSessionCreated, "sess456", "user789", 0)
	if sessionEvent.Type() != EventSessionCreated {
		t.Errorf("Expected type %s, got %s", EventSessionCreated, sessionEvent.Type())
	}

	// Test skill event
	skillEvent := NewSkillEvent(EventSkillStarted, "weather", "current", "", nil, 0)
	if skillEvent.Type() != EventSkillStarted {
		t.Errorf("Expected type %s, got %s", EventSkillStarted, skillEvent.Type())
	}

	// Test task event
	taskEvent := NewTaskEvent(EventTaskCreated, "task789", "sess456", "weather", "pending", "get weather", "", "")
	if taskEvent.Type() != EventTaskCreated {
		t.Errorf("Expected type %s, got %s", EventTaskCreated, taskEvent.Type())
	}
}
