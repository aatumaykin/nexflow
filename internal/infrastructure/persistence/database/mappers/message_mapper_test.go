package mappers

import (
	"testing"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
	dbmodel "github.com/atumaikin/nexflow/internal/infrastructure/persistence/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessageToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbMessage   *dbmodel.Message
		expected    *entity.Message
		expectedNil bool
	}{
		{
			name: "Valid message",
			dbMessage: &dbmodel.Message{
				ID:        "message-id",
				SessionID: "session-id",
				Role:      "user",
				Content:   "Hello, world!",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Message{
				ID:        valueobject.MessageID("message-id"),
				SessionID: valueobject.MustNewSessionID("session-id"),
				Role:      valueobject.MustNewMessageRole("user"),
				Content:   "Hello, world!",
				CreatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name: "Valid assistant message",
			dbMessage: &dbmodel.Message{
				ID:        "message-id-2",
				SessionID: "session-id-2",
				Role:      "assistant",
				Content:   "This is a response",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expected: &entity.Message{
				ID:        valueobject.MessageID("message-id-2"),
				SessionID: valueobject.MustNewSessionID("session-id-2"),
				Role:      valueobject.MustNewMessageRole("assistant"),
				Content:   "This is a response",
				CreatedAt: time.Now(),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			dbMessage:   nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessageToDomain(tt.dbMessage)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SessionID, result.SessionID)
			assert.Equal(t, tt.expected.Role, result.Role)
			assert.Equal(t, tt.expected.Content, result.Content)
			assert.WithinDuration(t, tt.expected.CreatedAt, result.CreatedAt, time.Second)
		})
	}
}

func TestMessageToDB(t *testing.T) {
	tests := []struct {
		name        string
		message     *entity.Message
		expected    *dbmodel.Message
		expectedNil bool
	}{
		{
			name: "Valid message",
			message: &entity.Message{
				ID:        valueobject.MessageID("message-id"),
				SessionID: valueobject.MustNewSessionID("session-id"),
				Role:      valueobject.MustNewMessageRole("user"),
				Content:   "Hello, world!",
				CreatedAt: time.Now(),
			},
			expected: &dbmodel.Message{
				ID:        "message-id",
				SessionID: "session-id",
				Role:      "user",
				Content:   "Hello, world!",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name: "Valid assistant message",
			message: &entity.Message{
				ID:        valueobject.MessageID("message-id-2"),
				SessionID: valueobject.MustNewSessionID("session-id-2"),
				Role:      valueobject.MustNewMessageRole("assistant"),
				Content:   "This is a response",
				CreatedAt: time.Now(),
			},
			expected: &dbmodel.Message{
				ID:        "message-id-2",
				SessionID: "session-id-2",
				Role:      "assistant",
				Content:   "This is a response",
				CreatedAt: time.Now().Format(time.RFC3339),
			},
			expectedNil: false,
		},
		{
			name:        "Nil input",
			message:     nil,
			expected:    nil,
			expectedNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessageToDB(tt.message)

			if tt.expectedNil {
				assert.Nil(t, result)
				return
			}

			require.NotNil(t, result)
			assert.Equal(t, tt.expected.ID, result.ID)
			assert.Equal(t, tt.expected.SessionID, result.SessionID)
			assert.Equal(t, tt.expected.Role, result.Role)
			assert.Equal(t, tt.expected.Content, result.Content)
			assert.Equal(t, tt.expected.CreatedAt, result.CreatedAt)
		})
	}
}

func TestMessagesToDomain(t *testing.T) {
	tests := []struct {
		name        string
		dbMessages  []dbmodel.Message
		expectedLen int
	}{
		{
			name: "Multiple messages",
			dbMessages: []dbmodel.Message{
				{
					ID:        "msg-1",
					SessionID: "session-1",
					Role:      "user",
					Content:   "Hello",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "msg-2",
					SessionID: "session-1",
					Role:      "assistant",
					Content:   "Hi there",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
				{
					ID:        "msg-3",
					SessionID: "session-1",
					Role:      "user",
					Content:   "How are you?",
					CreatedAt: time.Now().Format(time.RFC3339),
				},
			},
			expectedLen: 3,
		},
		{
			name:        "Empty slice",
			dbMessages:  []dbmodel.Message{},
			expectedLen: 0,
		},
		{
			name:        "Nil input",
			dbMessages:  nil,
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MessagesToDomain(tt.dbMessages)

			if tt.dbMessages == nil {
				assert.Empty(t, result)
				return
			}

			assert.Len(t, result, tt.expectedLen)
			for i, msg := range result {
				assert.Equal(t, tt.dbMessages[i].ID, string(msg.ID))
				assert.Equal(t, tt.dbMessages[i].Role, string(msg.Role))
				assert.Equal(t, tt.dbMessages[i].Content, msg.Content)
			}
		})
	}
}
