package valueobject

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestID_String(t *testing.T) {
	id := ID("test-id")
	assert.Equal(t, "test-id", id.String())
}

func TestID_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		expected bool
	}{
		{"empty ID", ID(""), true},
		{"non-empty ID", ID("test"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsEmpty())
		})
	}
}

func TestID_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		id       ID
		expected bool
	}{
		{"empty ID", ID(""), false},
		{"valid simple ID", ID("test"), true},
		{"valid ID with numbers", ID("test123"), true},
		{"valid ID with underscore", ID("test_id"), true},
		{"valid ID with hyphen", ID("test-id"), true},
		{"valid ID with mixed", ID("test-123_id"), true},
		{"invalid ID with space", ID("test id"), false},
		{"invalid ID with special char", ID("test@id"), false},
		{"invalid ID with dot", ID("test.id"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.id.IsValid())
		})
	}
}

func TestID_Equals(t *testing.T) {
	id1 := ID("test-id")
	id2 := ID("test-id")
	id3 := ID("other-id")

	assert.True(t, id1.Equals(id2))
	assert.False(t, id1.Equals(id3))
}

func TestID_MarshalJSON(t *testing.T) {
	id := ID("test-id")
	data, err := json.Marshal(id)

	require.NoError(t, err)
	assert.JSONEq(t, `"test-id"`, string(data))
}

func TestID_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		data      string
		expected  ID
		expectErr bool
	}{
		{"valid ID", `"test-id"`, ID("test-id"), false},
		{"empty ID", `""`, ID(""), true},
		{"invalid ID", `"test id"`, ID(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id ID
			err := json.Unmarshal([]byte(tt.data), &id)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestNewID(t *testing.T) {
	tests := []struct {
		name      string
		idStr     string
		expected  ID
		expectErr bool
	}{
		{"valid ID", "test-id", ID("test-id"), false},
		{"valid ID with numbers", "test123", ID("test123"), false},
		{"empty ID", "", ID(""), true},
		{"invalid ID", "test id", ID(""), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := NewID(tt.idStr)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, id)
			}
		})
	}
}

func TestMustNewID(t *testing.T) {
	assert.NotPanics(t, func() {
		MustNewID("test-id")
	})

	assert.Panics(t, func() {
		MustNewID("")
	})
}

func TestGenerateID(t *testing.T) {
	tests := []struct {
		name        string
		generator   func() string
		expectEmpty bool
	}{
		{
			name:        "with custom generator",
			generator:   func() string { return "custom-id" },
			expectEmpty: false,
		},
		{
			name:        "with nil generator",
			generator:   nil,
			expectEmpty: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := GenerateID(tt.generator)
			assert.False(t, id.IsEmpty())
		})
	}
}

// Test typed ID types

func TestUserID_String(t *testing.T) {
	id := UserID("user-123")
	assert.Equal(t, "user-123", id.String())
}

func TestNewUserID(t *testing.T) {
	id, err := NewUserID("user-123")

	require.NoError(t, err)
	assert.Equal(t, UserID("user-123"), id)
}

func TestNewUserID_Error(t *testing.T) {
	_, err := NewUserID("")
	assert.Error(t, err)
}

func TestTaskID_String(t *testing.T) {
	id := TaskID("task-456")
	assert.Equal(t, "task-456", id.String())
}

func TestNewTaskID(t *testing.T) {
	id, err := NewTaskID("task-456")

	require.NoError(t, err)
	assert.Equal(t, TaskID("task-456"), id)
}

func TestSessionID_String(t *testing.T) {
	id := SessionID("session-789")
	assert.Equal(t, "session-789", id.String())
}

func TestNewSessionID(t *testing.T) {
	id, err := NewSessionID("session-789")

	require.NoError(t, err)
	assert.Equal(t, SessionID("session-789"), id)
}

func TestMessageID_String(t *testing.T) {
	id := MessageID("msg-000")
	assert.Equal(t, "msg-000", id.String())
}

func TestNewMessageID(t *testing.T) {
	id, err := NewMessageID("msg-000")

	require.NoError(t, err)
	assert.Equal(t, MessageID("msg-000"), id)
}

func TestSkillID_String(t *testing.T) {
	id := SkillID("skill-111")
	assert.Equal(t, "skill-111", id.String())
}

func TestNewSkillID(t *testing.T) {
	id, err := NewSkillID("skill-111")

	require.NoError(t, err)
	assert.Equal(t, SkillID("skill-111"), id)
}

func TestScheduleID_String(t *testing.T) {
	id := ScheduleID("schedule-222")
	assert.Equal(t, "schedule-222", id.String())
}

func TestNewScheduleID(t *testing.T) {
	id, err := NewScheduleID("schedule-222")

	require.NoError(t, err)
	assert.Equal(t, ScheduleID("schedule-222"), id)
}

func TestLogID_String(t *testing.T) {
	id := LogID("log-333")
	assert.Equal(t, "log-333", id.String())
}

func TestNewLogID(t *testing.T) {
	id, err := NewLogID("log-333")

	require.NoError(t, err)
	assert.Equal(t, LogID("log-333"), id)
}

func TestStringToIDType(t *testing.T) {
	tests := []struct {
		name      string
		typeName  string
		idStr     string
		expectErr bool
	}{
		{"UserID", "userid", "user-123", false},
		{"SessionID", "sessionid", "session-456", false},
		{"TaskID", "taskid", "task-789", false},
		{"MessageID", "messageid", "msg-000", false},
		{"SkillID", "skillid", "skill-111", false},
		{"ScheduleID", "scheduleid", "schedule-222", false},
		{"LogID", "logid", "log-333", false},
		{"unknown type", "unknown", "test", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, err := StringToIDType(tt.typeName, tt.idStr)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, id)
			}
		})
	}
}

func TestUserID_MarshalJSON(t *testing.T) {
	id := UserID("user-123")
	data, err := json.Marshal(id)

	require.NoError(t, err)
	assert.JSONEq(t, `"user-123"`, string(data))
}

func TestUserID_UnmarshalJSON(t *testing.T) {
	var id UserID
	err := json.Unmarshal([]byte(`"user-123"`), &id)

	require.NoError(t, err)
	assert.Equal(t, UserID("user-123"), id)
}

func TestTaskID_MarshalJSON(t *testing.T) {
	id := TaskID("task-456")
	data, err := json.Marshal(id)

	require.NoError(t, err)
	assert.JSONEq(t, `"task-456"`, string(data))
}

func TestTaskID_UnmarshalJSON(t *testing.T) {
	var id TaskID
	err := json.Unmarshal([]byte(`"task-456"`), &id)

	require.NoError(t, err)
	assert.Equal(t, TaskID("task-456"), id)
}
