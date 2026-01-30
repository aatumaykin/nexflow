package mock

import (
	"context"
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	"github.com/atumaikin/nexflow/internal/infrastructure/channels"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTelegramConnector_StartStop tests Start and Stop methods
func TestTelegramConnector_StartStop(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name    string
		setup   func(*TelegramConnector)
		wantErr bool
		errMsg  string
	}{
		{
			name:    "successful start",
			setup:   func(c *TelegramConnector) {},
			wantErr: false,
		},
		{
			name: "start when already running",
			setup: func(c *TelegramConnector) {
				_ = c.Start(ctx)
			},
			wantErr: true,
			errMsg:  "already running",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewTelegramConnector()
			if tt.setup != nil {
				tt.setup(conn)
			}

			err := conn.Start(ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.True(t, conn.IsRunning())
			}
		})
	}
}

// TestTelegramConnector_SendResponse tests sending responses
func TestTelegramConnector_SendResponse(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		response  *channels.Response
		setup     func(*TelegramConnector)
		wantErr   bool
		errMsg    string
		checkSent bool
	}{
		{
			name:   "send response when running",
			userID: "user-1",
			response: &channels.Response{
				Content: "Test response",
			},
			checkSent: true,
		},
		{
			name:   "send response with metadata",
			userID: "user-2",
			response: &channels.Response{
				Content:  "Test with metadata",
				Metadata: map[string]interface{}{"source": "test", "time": time.Now()},
			},
			checkSent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewTelegramConnector()
			_ = conn.Start(ctx) // Ensure running

			err := conn.SendResponse(ctx, tt.userID, tt.response)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				assert.True(t, conn.IsRunning())

				if tt.checkSent {
					responses := conn.GetResponses()
					assert.Len(t, responses, 1)
					assert.Equal(t, tt.userID, responses[0].userID)
					assert.Equal(t, tt.response.Content, responses[0].response.Content)
				}
			}
		})
	}
}

// TestTelegramConnector_UserOperations tests GetUser and CreateUser
func TestTelegramConnector_UserOperations(t *testing.T) {
	ctx := context.Background()
	userID := "user-1"

	tests := []struct {
		name          string
		setup         func(*TelegramConnector)
		channelUserID string
		wantErr       bool
		errMsg        string
		checkUser     bool
		action        func(*TelegramConnector, context.Context) (*entity.User, error)
	}{
		{
			name: "create new user",
			setup: func(c *TelegramConnector) {
				_ = c.Start(ctx)
			},
			channelUserID: userID,
			checkUser:     true,
			action: func(c *TelegramConnector, ctx context.Context) (*entity.User, error) {
				return c.CreateUser(ctx, userID)
			},
		},
		{
			name: "get existing user",
			setup: func(c *TelegramConnector) {
				_ = c.Start(ctx)
				_, err := c.CreateUser(ctx, userID)
				require.NoError(t, err)
			},
			channelUserID: userID,
			checkUser:     true,
			action: func(c *TelegramConnector, ctx context.Context) (*entity.User, error) {
				return c.GetUser(ctx, userID)
			},
		},
		{
			name: "get non-existent user",
			setup: func(c *TelegramConnector) {
				_ = c.Start(ctx)
			},
			channelUserID: "non-existent-user",
			wantErr:       true,
			errMsg:        "not found",
			action: func(c *TelegramConnector, ctx context.Context) (*entity.User, error) {
				return c.GetUser(ctx, "non-existent-user")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn := NewTelegramConnector()
			if tt.setup != nil {
				tt.setup(conn)
			}

			user, err := tt.action(conn, ctx)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, user)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.channelUserID, user.ChannelID)
				channelType := valueobject.Channel("telegram")
				assert.True(t, user.Channel.Equals(channelType))
			}
		})
	}
}

// TestTelegramConnector_IncomingChannel tests incoming messages channel
func TestTelegramConnector_IncomingChannel(t *testing.T) {
	conn := NewTelegramConnector()
	_ = conn.Start(context.Background())

	ch := conn.Incoming()
	assert.NotNil(t, ch)
	// Send a test message to verify the channel is working
	err := conn.SendTestMessage("user-1", "channel-1", "test message")
	require.NoError(t, err)

	// Verify message was received
	select {
	case msg := <-ch:
		assert.Equal(t, "user-1", msg.UserID)
		assert.Equal(t, "channel-1", msg.ChannelID)
		assert.Equal(t, "test message", msg.Content)
	default:
		t.Fatal("expected to receive message from incoming channel")
	}
}

// TestTelegramConnector_IsRunning tests IsRunning method
func TestTelegramConnector_IsRunning(t *testing.T) {
	conn := NewTelegramConnector()
	assert.False(t, conn.IsRunning())

	_ = conn.Start(context.Background())
	assert.True(t, conn.IsRunning())

	_ = conn.Stop(context.Background())
	assert.False(t, conn.IsRunning())
}

// TestTelegramConnector_SendTestMessage tests sending test messages
func TestTelegramConnector_SendTestMessage(t *testing.T) {
	ctx := context.Background()
	conn := NewTelegramConnector()
	_ = conn.Start(ctx)

	tests := []struct {
		name      string
		userID    string
		channelID string
		content   string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "send valid test message",
			userID:    "user-1",
			channelID: "channel-1",
			content:   "Hello, Telegram!",
		},
		{
			name:      "send empty content",
			userID:    "user-2",
			channelID: "channel-2",
			content:   "",
		},
		{
			name:      "send long content",
			userID:    "user-3",
			channelID: "channel-3",
			content:   "This is a very long message that should still be handled correctly for Telegram.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := conn.SendTestMessage(tt.userID, tt.channelID, tt.content)
			require.NoError(t, err)

			select {
			case msg := <-conn.Incoming():
				assert.Equal(t, tt.userID, msg.UserID)
				assert.Equal(t, tt.channelID, msg.ChannelID)
				assert.Equal(t, tt.content, msg.Content)
			default:
				t.Fatal("expected message to be received")
			}
		})
	}
}

// TestTelegramConnector_GetResponses tests GetResponses method
func TestTelegramConnector_GetResponses(t *testing.T) {
	ctx := context.Background()
	conn := NewTelegramConnector()
	_ = conn.Start(ctx)

	// Send multiple responses
	for i := 0; i < 5; i++ {
		err := conn.SendResponse(ctx, "user-1", &channels.Response{
			Content: "Response " + string(rune('0'+i)),
		})
		require.NoError(t, err)
	}

	responses := conn.GetResponses()
	assert.Len(t, responses, 5)

	// Clear and verify
	conn.ClearResponses()
	responses = conn.GetResponses()
	assert.Len(t, responses, 0)
}

// TestTelegramConnector_MultipleSenders tests concurrent sends
func TestTelegramConnector_MultipleSenders(t *testing.T) {
	ctx := context.Background()
	conn := NewTelegramConnector()
	_ = conn.Start(ctx)

	done := make(chan bool, 5)

	for i := 0; i < 5; i++ {
		go func(i int) {
			err := conn.SendResponse(ctx, "user-"+string(rune('0'+i)), &channels.Response{
				Content: "Message " + string(rune('0'+i)),
			})
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	for i := 0; i < 5; i++ {
		<-done
	}

	responses := conn.GetResponses()
	assert.Len(t, responses, 5)
}

// TestTelegramConnector_NilPointer tests nil pointer handling
func TestTelegramConnector_NilPointer(t *testing.T) {
	var conn *TelegramConnector

	ctx := context.Background()

	// Defer recover to handle potential panics from nil pointer calls
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic (expected for nil pointer): %v", r)
		}
	}()

	// These calls will panic with nil pointer, so we test with recover
	_ = conn.IsRunning()
	_ = conn.Incoming()
	_ = conn.Start(ctx)
	_ = conn.Stop(ctx)
	_, _ = conn.GetUser(ctx, "user-1")
	_, _ = conn.CreateUser(ctx, "user-1")
}
