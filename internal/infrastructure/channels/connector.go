package channels

import (
	"context"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// Message represents a message from a channel
type Message struct {
	UserID    string // Channel-specific user ID
	ChannelID string // Channel-specific channel or chat ID
	Content   string // Message content
	Metadata  map[string]interface{}
}

// Response represents a response to send back to a channel
type Response struct {
	Content  string
	Metadata map[string]interface{}
}

// Connector defines the interface for all channel connectors
type Connector interface {
	// Name returns the name of the channel (telegram, discord, web, etc.)
	Name() string

	// Start initializes and starts the connector
	Start(ctx context.Context) error

	// Stop gracefully stops the connector
	Stop(ctx context.Context) error

	// SendResponse sends a response to a user
	SendResponse(ctx context.Context, userID string, response *Response) error

	// Incoming returns a channel for incoming messages
	Incoming() <-chan *Message

	// IsRunning returns whether the connector is currently running
	IsRunning() bool

	// GetUser retrieves a user by channel-specific ID
	GetUser(ctx context.Context, channelUserID string) (*entity.User, error)

	// CreateUser creates a new user in the system
	CreateUser(ctx context.Context, channelUserID string) (*entity.User, error)
}
