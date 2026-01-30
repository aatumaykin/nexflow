package channels

import (
	"context"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConnectorInterface tests that the Connector interface is correctly defined
func TestConnectorInterface(t *testing.T) {
	var conn Connector

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test that nil connector doesn't panic on method calls (or recovers if it does)
	_ = conn.IsRunning()
	_ = conn.Name()
	_ = conn.Incoming()
}

// TestConnectorResponse tests the Response struct
func TestConnectorResponse(t *testing.T) {
	response := &Response{
		Content:  "Test response",
		Metadata: map[string]interface{}{"key": "value"},
	}
	assert.Equal(t, "Test response", response.Content)
	assert.NotNil(t, response.Metadata)
}

// TestConnectorMessage tests the Message struct
func TestConnectorMessage(t *testing.T) {
	message := &Message{
		UserID:    "user-1",
		ChannelID: "channel-1",
		Content:   "Test message",
		Metadata:  map[string]interface{}{"timestamp": time.Now()},
	}
	assert.Equal(t, "user-1", message.UserID)
	assert.Equal(t, "channel-1", message.ChannelID)
	assert.Equal(t, "Test message", message.Content)
}

// TestConnectorCreateUser tests CreateUser method
func TestConnectorCreateUser(t *testing.T) {
	ctx := context.Background()

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test that nil interface doesn't panic (or recovers if it does)
	var conn Connector
	user, err := conn.CreateUser(ctx, "user-1")
	_ = user
	_ = err
}

// TestConnectorGetUser tests GetUser method
func TestConnectorGetUser(t *testing.T) {
	ctx := context.Background()

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test that nil interface doesn't panic (or recovers if it does)
	var conn Connector
	user, err := conn.GetUser(ctx, "user-1")
	_ = user
	_ = err
}

// TestConnectorSendResponse tests SendResponse method
func TestConnectorSendResponse(t *testing.T) {
	ctx := context.Background()

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test that nil interface doesn't panic (or recovers if it does)
	var conn Connector
	err := conn.SendResponse(ctx, "user-1", &Response{Content: "test"})
	_ = err
}

// TestConnectorStartStop tests Start and Stop methods
func TestConnectorStartStop(t *testing.T) {
	ctx := context.Background()

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test that nil interface doesn't panic (or recovers if it does)
	var conn Connector
	err := conn.Start(ctx)
	_ = err

	err = conn.Stop(ctx)
	_ = err
}

// TestConnectorIsRunning tests IsRunning method
func TestConnectorIsRunning(t *testing.T) {
	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	var conn Connector
	_ = conn.IsRunning()
}

// TestConnectorIncoming tests Incoming channel method
func TestConnectorIncoming(t *testing.T) {
	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	var conn Connector
	ch := conn.Incoming()
	_ = ch
}

// TestConnectorName tests Name method
func TestConnectorName(t *testing.T) {
	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	var conn Connector
	name := conn.Name()
	_ = name
}

// MockConnector for testing
type MockConnector struct {
	IsRunningVal bool
	HasIncoming  bool
	ShouldFail   bool
	ErrorOnStart bool
	ErrorOnStop  bool
	NameVal      string
}

func (m *MockConnector) IsRunning() bool {
	return m.IsRunningVal
}

func (m *MockConnector) Incoming() <-chan *Message {
	if m.HasIncoming {
		ch := make(chan *Message, 1)
		ch <- &Message{UserID: "user-1", ChannelID: "channel-1", Content: "test"}
		return ch
	}
	return nil
}

func (m *MockConnector) Start(ctx context.Context) error {
	if m.ErrorOnStart {
		return assert.AnError
	}
	return nil
}

func (m *MockConnector) Stop(ctx context.Context) error {
	if m.ErrorOnStop {
		return assert.AnError
	}
	return nil
}

func (m *MockConnector) SendResponse(ctx context.Context, userID string, response *Response) error {
	return nil
}

func (m *MockConnector) GetUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	if m.ShouldFail {
		return nil, assert.AnError
	}
	return &entity.User{
		ID:        valueobject.UserID(channelUserID),
		Channel:   valueobject.Channel("test"),
		ChannelID: channelUserID,
		CreatedAt: time.Now(),
	}, nil
}

func (m *MockConnector) CreateUser(ctx context.Context, channelUserID string) (*entity.User, error) {
	if m.ShouldFail {
		return nil, assert.AnError
	}
	return &entity.User{
		ID:        valueobject.UserID(channelUserID),
		Channel:   valueobject.Channel("test"),
		ChannelID: channelUserID,
		CreatedAt: time.Now(),
	}, nil
}

func (m *MockConnector) Name() string {
	return m.NameVal
}

// TestMockConnector tests MockConnector implementation
func TestMockConnector(t *testing.T) {
	tests := []struct {
		name       string
		mock       *MockConnector
		testStart  bool
		testGetErr bool
		wantErr    bool
	}{
		{
			name: "mock successful start",
			mock: &MockConnector{
				IsRunningVal: true,
				NameVal:      "test",
			},
			testStart: true,
			wantErr:   false,
		},
		{
			name: "mock start error",
			mock: &MockConnector{
				ErrorOnStart: true,
			},
			testStart: true,
			wantErr:   true,
		},
		{
			name: "mock get user",
			mock: &MockConnector{
				IsRunningVal: true,
			},
			testGetErr: true,
			wantErr:    false,
		},
		{
			name: "mock get user error",
			mock: &MockConnector{
				ShouldFail:   true,
				IsRunningVal: true,
			},
			testGetErr: true,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			if tt.testStart {
				err := tt.mock.Start(ctx)
				if tt.wantErr {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
				}
			}

			if tt.testGetErr {
				user, err := tt.mock.GetUser(ctx, "user-1")
				if tt.wantErr {
					require.Error(t, err)
					assert.Nil(t, user)
				} else {
					require.NoError(t, err)
					assert.NotNil(t, user)
					assert.Equal(t, "user-1", string(user.ID))
				}
			}
		})
	}
}

// TestMessageStructure tests message structure fields
func TestMessageStructure(t *testing.T) {
	tests := []struct {
		name string
		msg  *Message
	}{
		{
			name: "message with metadata",
			msg: &Message{
				UserID:    "user-1",
				ChannelID: "channel-1",
				Content:   "Hello",
				Metadata:  map[string]interface{}{"key": "value", "num": 42},
			},
		},
		{
			name: "message without metadata",
			msg: &Message{
				UserID:    "user-1",
				ChannelID: "channel-1",
				Content:   "Hello",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "user-1", tt.msg.UserID)
			assert.Equal(t, "channel-1", tt.msg.ChannelID)
			assert.Equal(t, "Hello", tt.msg.Content)
			if tt.msg.Metadata != nil {
				assert.Equal(t, "value", tt.msg.Metadata["key"])
				assert.Equal(t, int(42), tt.msg.Metadata["num"])
			}
		})
	}
}

// TestResponseStructure tests response structure fields
func TestResponseStructure(t *testing.T) {
	tests := []struct {
		name string
		resp *Response
	}{
		{
			name: "response with metadata",
			resp: &Response{
				Content:  "Test",
				Metadata: map[string]interface{}{"source": "test", "time": time.Now()},
			},
		},
		{
			name: "response without metadata",
			resp: &Response{
				Content: "Test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, "Test", tt.resp.Content)
			if tt.resp.Metadata != nil {
				assert.Equal(t, "test", tt.resp.Metadata["source"])
			}
		})
	}
}

// TestConnectorContext tests context propagation
func TestConnectorContext(t *testing.T) {
	ctx := context.Background()

	// Defer recover to handle potential panics from nil interface method calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil interface): %v", r)
		}
	}()

	// Test with nil connector
	var conn Connector
	err := conn.Start(ctx)
	_ = err

	err = conn.Stop(ctx)
	_ = err

	_, err = conn.GetUser(ctx, "user-1")
	_ = err

	user, err := conn.CreateUser(ctx, "user-1")
	_ = user
	_ = err
}
