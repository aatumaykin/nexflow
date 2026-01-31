package ports

import (
	"context"
	"errors"
	"testing"

	"github.com/atumaikin/nexflow/internal/shared/eventbus"
)

// mockEventHandler is a mock implementation of EventHandler for testing
type mockEventHandler struct {
	handleFunc func(ctx context.Context, event eventbus.Event) error
}

// Handle implements EventHandler interface
func (m *mockEventHandler) Handle(ctx context.Context, event eventbus.Event) error {
	if m.handleFunc != nil {
		return m.handleFunc(ctx, event)
	}
	return nil
}

// TestEventHandlerInterface verifies that EventHandler interface is correctly defined
func TestEventHandlerInterface(t *testing.T) {
	// Create a mock event handler
	handler := &mockEventHandler{
		handleFunc: func(ctx context.Context, event eventbus.Event) error {
			return nil
		},
	}

	// Verify that the handler implements the interface
	var _ EventHandler = handler
}

// TestEventHandlerHandleSuccess tests successful event handling
func TestEventHandlerHandleSuccess(t *testing.T) {
	ctx := context.Background()
	event := eventbus.NewBaseEvent("test.event", nil)

	handler := &mockEventHandler{
		handleFunc: func(ctx context.Context, event eventbus.Event) error {
			// Verify event is not nil
			if event == nil {
				t.Error("expected event to be non-nil")
			}
			return nil
		},
	}

	err := handler.Handle(ctx, event)
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

// TestEventHandlerHandleError tests event handling with error
func TestEventHandlerHandleError(t *testing.T) {
	ctx := context.Background()
	event := eventbus.NewBaseEvent("test.event", nil)

	expectedError := errors.New("test error")

	handler := &mockEventHandler{
		handleFunc: func(ctx context.Context, event eventbus.Event) error {
			return expectedError
		},
	}

	err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("expected error, got nil")
	}
	if err != expectedError {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

// TestEventHandlerHandleWithContextCancellation tests event handling with cancelled context
func TestEventHandlerHandleWithContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	event := eventbus.NewBaseEvent("test.event", nil)

	handler := &mockEventHandler{
		handleFunc: func(ctx context.Context, event eventbus.Event) error {
			// Check if context is cancelled
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return nil
		},
	}

	err := handler.Handle(ctx, event)
	if err == nil {
		t.Error("expected context cancellation error, got nil")
	}
}
