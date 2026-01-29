package mock

import (
	"context"

	"github.com/atumaikin/nexflow/internal/application/ports"
)

// MockConnector is a mock implementation of Connector for testing
type MockConnector struct {
	StartFunc       func(context.Context) error
	StopFunc        func() error
	EventsFunc      func() <-chan ports.Event
	SendMessageFunc func(context.Context, string, string) error
	NameFunc        func() string
}

// NewMockConnector creates a new mock connector
func NewMockConnector() *MockConnector {
	return &MockConnector{
		StartFunc: func(ctx context.Context) error {
			// Mock start - ready to receive events
			return nil
		},
		StopFunc: func() error {
			// Mock stop - gracefully shutdown
			return nil
		},
		EventsFunc: func() <-chan ports.Event {
			// Mock events channel
			ch := make(chan ports.Event, 10)
			go func() {
				// Send some initial mock events
				ch <- ports.Event{
					ID:        "mock-event-1",
					Channel:   "mock",
					UserID:    "mock-user-1",
					Message:   "Hello from mock",
					Timestamp: "2024-01-01T00:00:00Z",
				}
			}()
			return ch
		},
		SendMessageFunc: func(ctx context.Context, userID, message string) error {
			// Mock message send
			return nil
		},
		NameFunc: func() string {
			return "MockConnector"
		},
	}
}

// Start implements Connector interface
func (m *MockConnector) Start(ctx context.Context) error {
	if m.StartFunc != nil {
		return m.StartFunc(ctx)
	}
	return nil
}

// Stop implements Connector interface
func (m *MockConnector) Stop() error {
	if m.StopFunc != nil {
		return m.StopFunc()
	}
	return nil
}

// Events implements Connector interface
func (m *MockConnector) Events() <-chan ports.Event {
	if m.EventsFunc != nil {
		return m.EventsFunc()
	}
	ch := make(chan ports.Event)
	close(ch)
	return ch
}

// SendMessage implements Connector interface
func (m *MockConnector) SendMessage(ctx context.Context, userID, message string) error {
	if m.SendMessageFunc != nil {
		return m.SendMessageFunc(ctx, userID, message)
	}
	return nil
}

// Name implements Connector interface
func (m *MockConnector) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "MockConnector"
}
