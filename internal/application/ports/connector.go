package ports

import (
	"context"
)

// Event represents an incoming event from a communication channel.
// Events are messages received from users through various channels.
type Event struct {
	ID        string            `json:"id"`        // Unique identifier for the event
	Channel   string            `json:"channel"`   // Channel type: "telegram", "discord", "web", etc.
	UserID    string            `json:"user_id"`   // ID of the user who sent the event
	Message   string            `json:"message"`   // Event message content
	Metadata  map[string]string `json:"metadata"`  // Additional event metadata
	Timestamp string            `json:"timestamp"` // ISO 8601 format timestamp
}

// Connector defines the interface for communication channels.
// Connectors handle bidirectional communication with users through various channels.
type Connector interface {
	// Start begins listening for events from the channel.
	// The context can be used to cancel the operation.
	Start(ctx context.Context) error

	// Stop gracefully stops the connector.
	Stop() error

	// Events returns a read-only channel that receives incoming events.
	Events() <-chan Event

	// SendMessage sends a response message to the specified user.
	SendMessage(ctx context.Context, userID, message string) error

	// Name returns the connector name.
	Name() string
}
