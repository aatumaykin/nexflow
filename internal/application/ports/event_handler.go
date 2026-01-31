package ports

import (
	"context"

	"github.com/atumaikin/nexflow/internal/shared/eventbus"
)

// EventHandler defines the interface for handling events in system components.
// Components that need to react to events (connectors, routers, skills, etc.) implement this interface.
type EventHandler interface {
	// Handle processes an event and returns an error if handling fails.
	//
	// Parameters:
	//   - ctx: Context for the operation (can be used for cancellation or timeouts)
	//   - event: Event to handle (implements eventbus.Event interface)
	//
	// Returns:
	//   - error: Error if event handling failed (nil on success)
	Handle(ctx context.Context, event eventbus.Event) error
}
