package mock

import (
	"context"
	"fmt"
	"sync"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
)

// TelegramConnector is a mock implementation of the Telegram channel connector
type TelegramConnector struct {
	mu        sync.RWMutex
	running   bool
	incoming  chan *channels.Message
	responses []mockResponse
	users     map[string]*entity.User
	name      string
}

type mockResponse struct {
	userID   string
	response *channels.Response
}

// NewTelegramConnector creates a new mock Telegram connector
func NewTelegramConnector() *TelegramConnector {
	return &TelegramConnector{
		name:     "telegram",
		incoming: make(chan *channels.Message, 100),
		users:    make(map[string]*entity.User),
	}
}

// Name returns the name of the channel
func (c *TelegramConnector) Name() string {
	return c.name
}

// Start initializes and starts the connector
func (c *TelegramConnector) Start(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.running {
		return fmt.Errorf("telegram connector is already running")
	}

	c.running = true
	return nil
}

// Stop gracefully stops the connector
func (c *TelegramConnector) Stop(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("telegram connector is not running")
	}

	c.running = false
	close(c.incoming)
	c.incoming = make(chan *channels.Message, 100)
	return nil
}

// SendResponse sends a response to a user
func (c *TelegramConnector) SendResponse(ctx context.Context, userID string, response *channels.Response) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.running {
		return fmt.Errorf("telegram connector is not running")
	}

	c.responses = append(c.responses, mockResponse{
		userID:   userID,
		response: response,
	})

	return nil
}

// Incoming returns a channel for incoming messages
func (c *TelegramConnector) Incoming() <-chan *channels.Message {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.incoming
}

// IsRunning returns whether the connector is currently running
func (c *TelegramConnector) IsRunning() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.running
}

// GetUser retrieves a user by channel-specific ID
func (c *TelegramConnector) GetUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	user, exists := c.users[channelUserID]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", channelUserID)
	}

	return user, nil
}

// CreateUser creates a new user in the system
func (c *TelegramConnector) CreateUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.users[channelUserID]; exists {
		return nil, fmt.Errorf("user already exists: %s", channelUserID)
	}

	user := entity.NewUser("telegram", channelUserID)
	c.users[channelUserID] = user

	return user, nil
}

// SendTestMessage sends a test message through the connector (for testing purposes)
func (c *TelegramConnector) SendTestMessage(userID, channelID, content string) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if !c.running {
		return fmt.Errorf("telegram connector is not running")
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
func (c *TelegramConnector) GetResponses() []mockResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	responses := make([]mockResponse, len(c.responses))
	copy(responses, c.responses)
	return responses
}

// ClearResponses clears all sent responses (for testing purposes)
func (c *TelegramConnector) ClearResponses() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.responses = nil
}
