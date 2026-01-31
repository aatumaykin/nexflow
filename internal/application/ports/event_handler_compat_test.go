package ports

import (
	"context"
	"testing"

	"github.com/atumaikin/nexflow/internal/shared/eventbus"
	"github.com/atumaikin/nexflow/internal/shared/logging"
)

// TestEventLoggerImplementsEventHandler verifies that EventLogger implements EventHandler interface
func TestEventLoggerImplementsEventHandler(t *testing.T) {
	// EventLogger is defined in eventbus package and has a Handle method
	// This test verifies it implements our EventHandler interface

	logger := logging.NewNoopLogger()
	eventLogger := eventbus.NewEventLogger(logger)

	var handler EventHandler = eventLogger

	if handler == nil {
		t.Error("expected EventLogger to implement EventHandler interface")
	}

	// Test calling Handle method
	event := eventbus.NewBaseEvent("test.event", nil)
	err := handler.Handle(context.Background(), event)
	if err != nil {
		t.Errorf("expected no error from EventLogger.Handle, got: %v", err)
	}
}
