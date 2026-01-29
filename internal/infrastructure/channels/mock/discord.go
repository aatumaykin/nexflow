package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
)

// DiscordConnector is a mock implementation of Discord channel connector
type DiscordConnector struct {
	mu        sync.RWMutex
	running   bool
	incoming  chan *channels.Message
	responses []mockResponse
	users     map[string]*entity.User
	name      string
}

// NewDiscordConnector creates a new mock Discord connector
func NewDiscordConnector() *DiscordConnector {
	return &DiscordConnector{
		name:     "discord",
		incoming: make(chan *channels.Message, 100),
		users:    make(map[string]*entity.User),
	}
}

// Name returns the name of the channel
func (c *DiscordConnector) Name() string {
	return c.name
}

// Start initializes and starts the connector
func (c *DiscordConnector) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("discord connector is already running")
	}

	c.running = true
	return nil
}

// Stop gracefully stops the connector
func (c *DiscordConnector) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("discord connector is not running")
	}

	c.running = false
	close(c.incoming)
	c.incoming = make(chan *channels.Message, 100)
	return nil
}

// SendResponse sends a response to a user
func (c *DiscordConnector) SendResponse(ctx context.Context, userID string, response *channels.Response) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("discord connector is not running")
	}

	c.responses = append(c.responses, mockResponse{
		userID:   userID,
		response: response,
	})

	return nil
}

// Incoming returns a channel for incoming messages
func (c *DiscordConnector) Incoming() <-chan *channels.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.incoming
}

// IsRunning returns whether the connector is currently running
func (c *DiscordConnector) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// GetUser retrieves a user by channel-specific ID
func (c *DiscordConnector) GetUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	user, exists := c.users[channelUserID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", channelUserID)
	}

	return user, nil
}

// CreateUser creates a new user in the system
func (c *DiscordConnector) CreateUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.users[channelUserID]; exists {
		return nil, fmt.Errorf("user already exists: %s", channelUserID)
	}

	user := entity.NewUser("discord", channelUserID)
	c.users[channelUserID] = user

	return user, nil
}

// SendTestMessage sends a test message through the connector (for testing purposes)
func (c *DiscordConnector) SendTestMessage(userID, channelID, content string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.running {
		return fmt.Errorf("discord connector is not running")
	}

	msg := &channels.Message{
		UserID:    userID,
		ChannelID: channelID,
		Content:   content,
		Metadata:  make(map[string]interface{}),
	}

	select {
	case c.incoming <- msg:
		return nil
	default:
		return fmt.Errorf("incoming channel is full")
	}
}

// GetResponses returns all sent responses (for testing purposes)
func (c *DiscordConnector) GetResponses() []mockResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	responses := make([]mockResponse, len(c.responses))
	copy(responses, c.responses)
	return responses
}

// ClearResponses clears all sent responses (for testing purposes)
func (c *DiscordConnector) ClearResponses() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.responses = nil
}
