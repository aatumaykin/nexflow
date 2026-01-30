package dto

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/atumaikin/nexflow/internal/domain/entity"
	"github.com/atumaikin/nexflow/internal/domain/valueobject"
)

// Helper function to create test timestamps
func getTestTime() time.Time {
	return time.Date(2026, 1, 30, 10, 0, 0, 0, time.UTC)
}

// Helper function to create test timestamps with both created and updated
func getTestTimeWithUpdated() (time.Time, time.Time) {
	created := time.Date(2026, 1, 30, 10, 0, 0, 0, time.UTC)
	updated := time.Date(2026, 1, 30, 12, 0, 0, 0, time.UTC)
	return created, updated
}

// ========== SessionDTO Tests ==========

func TestSessionDTO_ToEntity(t *testing.T) {
	created, updated := getTestTimeWithUpdated()
	dto := &SessionDTO{
		ID:        "session-1",
		UserID:    "user-1",
		CreatedAt: created.Format(time.RFC3339),
		UpdatedAt: updated.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.SessionID("session-1"), entity.ID)
	assert.Equal(t, valueobject.UserID("user-1"), entity.UserID)
	assert.Equal(t, created, entity.CreatedAt)
	assert.Equal(t, updated, entity.UpdatedAt)
}

func TestSessionDTOFromEntity(t *testing.T) {
	created := getTestTime()
	updated := getTestTime()
	entity := &entity.Session{
		ID:        valueobject.MustNewSessionID("session-1"),
		UserID:    valueobject.MustNewUserID("user-1"),
		CreatedAt: created,
		UpdatedAt: updated,
	}

	dto := SessionDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "session-1", dto.ID)
	assert.Equal(t, "user-1", dto.UserID)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
	assert.Equal(t, updated.Format(time.RFC3339), dto.UpdatedAt)
}

// ========== MessageDTO Tests ==========

func TestMessageDTO_ToEntity(t *testing.T) {
	created := getTestTime()
	dto := &MessageDTO{
		ID:        "msg-1",
		SessionID: "session-1",
		Role:      "user",
		Content:   "Hello",
		CreatedAt: created.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.MessageID("msg-1"), entity.ID)
	assert.Equal(t, valueobject.SessionID("session-1"), entity.SessionID)
	assert.Equal(t, valueobject.MessageRole("user"), entity.Role)
	assert.Equal(t, "Hello", entity.Content)
	assert.Equal(t, created, entity.CreatedAt)
}

func TestMessageDTOFromEntity(t *testing.T) {
	created := getTestTime()
	entity := &entity.Message{
		ID:        valueobject.MustNewMessageID("msg-1"),
		SessionID: valueobject.MustNewSessionID("session-1"),
		Role:      valueobject.MustNewMessageRole("user"),
		Content:   "Test content",
		CreatedAt: created,
	}

	dto := MessageDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "msg-1", dto.ID)
	assert.Equal(t, "session-1", dto.SessionID)
	assert.Equal(t, "user", dto.Role)
	assert.Equal(t, "Test content", dto.Content)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
}

// ========== TaskDTO Tests ==========

func TestTaskDTO_ToEntity(t *testing.T) {
	created, updated := getTestTimeWithUpdated()
	dto := &TaskDTO{
		ID:        "task-1",
		SessionID: "session-1",
		Skill:     "skill-1",
		Input:     "Input data",
		Output:    "",
		Status:    "pending",
		Error:     "",
		CreatedAt: created.Format(time.RFC3339),
		UpdatedAt: updated.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.TaskID("task-1"), entity.ID)
	assert.Equal(t, valueobject.SessionID("session-1"), entity.SessionID)
	assert.Equal(t, "skill-1", entity.Skill)
	assert.Equal(t, "Input data", entity.Input)
	assert.Equal(t, "", entity.Output)
	assert.Equal(t, valueobject.TaskStatus("pending"), entity.Status)
	assert.Equal(t, "", entity.Error)
	assert.Equal(t, created, entity.CreatedAt)
	assert.Equal(t, updated, entity.UpdatedAt)
}

func TestTaskDTOFromEntity(t *testing.T) {
	created, updated := getTestTimeWithUpdated()
	entity := &entity.Task{
		ID:        valueobject.MustNewTaskID("task-1"),
		SessionID: valueobject.MustNewSessionID("session-1"),
		Skill:     "skill-1",
		Input:     "Test input",
		Output:    "Test output",
		Status:    valueobject.MustNewTaskStatus("completed"),
		Error:     "",
		CreatedAt: created,
		UpdatedAt: updated,
	}

	dto := TaskDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "task-1", dto.ID)
	assert.Equal(t, "session-1", dto.SessionID)
	assert.Equal(t, "skill-1", dto.Skill)
	assert.Equal(t, "Test input", dto.Input)
	assert.Equal(t, "Test output", dto.Output)
	assert.Equal(t, "completed", dto.Status)
	assert.Equal(t, "", dto.Error)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
	assert.Equal(t, updated.Format(time.RFC3339), dto.UpdatedAt)
}

// ========== UserDTO Tests ==========

func TestUserDTO_ToEntity(t *testing.T) {
	created := getTestTime()
	dto := &UserDTO{
		ID:        "user-1",
		Channel:   "telegram",
		ChannelID: "12345",
		CreatedAt: created.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.UserID("user-1"), entity.ID)
	assert.Equal(t, valueobject.MustNewChannel("telegram"), entity.Channel)
	assert.Equal(t, "12345", entity.ChannelID)
	assert.Equal(t, created, entity.CreatedAt)
}

func TestUserDTOFromEntity(t *testing.T) {
	created := getTestTime()
	entity := &entity.User{
		ID:        valueobject.UserID("user-1"),
		Channel:   valueobject.MustNewChannel("discord"),
		ChannelID: "67890",
		CreatedAt: created,
	}

	dto := UserDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "user-1", dto.ID)
	assert.Equal(t, "discord", dto.Channel)
	assert.Equal(t, "67890", dto.ChannelID)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
}

// ========== SkillDTO Tests ==========

func TestSkillDTO_ToEntity(t *testing.T) {
	created := getTestTime()
	dto := &SkillDTO{
		ID:          "skill-1",
		Name:        "TestSkill",
		Version:     "1.0.0",
		Location:    "/path/to/skill",
		Permissions: `["read", "write"]`,
		Metadata:    `{"key": "value"}`,
		CreatedAt:   created.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.SkillID("skill-1"), entity.ID)
	assert.Equal(t, "TestSkill", entity.Name)
	assert.Equal(t, valueobject.MustNewVersion("1.0.0"), entity.Version)
	assert.Equal(t, "/path/to/skill", entity.Location)
	assert.Equal(t, `["read", "write"]`, entity.Permissions)
	assert.Equal(t, `{"key": "value"}`, entity.Metadata)
	assert.Equal(t, created, entity.CreatedAt)
}

func TestSkillDTOFromEntity(t *testing.T) {
	created := getTestTime()
	entity := &entity.Skill{
		ID:          valueobject.MustNewSkillID("skill-1"),
		Name:        "TestSkill",
		Version:     valueobject.MustNewVersion("1.0.0"),
		Location:    "/path/to/skill",
		Permissions: `["read", "write"]`,
		Metadata:    `{"key": "value"}`,
		CreatedAt:   created,
	}

	dto := SkillDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "skill-1", dto.ID)
	assert.Equal(t, "TestSkill", dto.Name)
	assert.Equal(t, "1.0.0", dto.Version)
	assert.Equal(t, "/path/to/skill", dto.Location)
	assert.Equal(t, `["read", "write"]`, dto.Permissions)
	assert.Equal(t, `{"key": "value"}`, dto.Metadata)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
}

// ========== ScheduleDTO Tests ==========

func TestScheduleDTO_ToEntity(t *testing.T) {
	created := getTestTime()
	dto := &ScheduleDTO{
		ID:             "schedule-1",
		Skill:          "skill-1",
		CronExpression: "0 0 * * *",
		Input:          "Test input",
		Enabled:        true,
		CreatedAt:      created.Format(time.RFC3339),
	}

	entity := dto.ToEntity()
	assert.NotNil(t, entity)
	assert.Equal(t, valueobject.ScheduleID("schedule-1"), entity.ID)
	assert.Equal(t, "skill-1", entity.Skill)
	assert.Equal(t, valueobject.MustNewCronExpression("0 0 * * *"), entity.CronExpression)
	assert.Equal(t, "Test input", entity.Input)
	assert.Equal(t, true, entity.Enabled)
	assert.Equal(t, created, entity.CreatedAt)
}

func TestScheduleDTOFromEntity(t *testing.T) {
	created := getTestTime()
	entity := &entity.Schedule{
		ID:             valueobject.MustNewScheduleID("schedule-1"),
		Skill:          "skill-1",
		CronExpression: valueobject.MustNewCronExpression("0 0 * * *"),
		Input:          "Test input",
		Enabled:        true,
		CreatedAt:      created,
	}

	dto := ScheduleDTOFromEntity(entity)
	assert.NotNil(t, dto)
	assert.Equal(t, "schedule-1", dto.ID)
	assert.Equal(t, "skill-1", dto.Skill)
	assert.Equal(t, "0 0 * * *", dto.CronExpression)
	assert.Equal(t, "Test input", dto.Input)
	assert.Equal(t, true, dto.Enabled)
	assert.Equal(t, created.Format(time.RFC3339), dto.CreatedAt)
}

// ========== Utility Functions Tests ==========

func TestMapToString(t *testing.T) {
	tests := []struct {
		name    string
		m       map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "non-empty map",
			m: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
			want:    `{"key1":"value1","key2":42,"key3":true}`,
			wantErr: false,
		},
		{
			name:    "nil map",
			m:       nil,
			want:    `{}`,
			wantErr: false,
		},
		{
			name: "empty map",
			m:    map[string]interface{}{},
			want: `{}`,
		},
		{
			name: "nested map",
			m: map[string]interface{}{
				"outer": map[string]interface{}{
					"inner": "value",
				},
			},
			want:    `{"outer":{"inner":"value"}}`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MapToString(tt.m)
			if (err != nil) != tt.wantErr {
				t.Errorf("MapToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringToMap(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "valid JSON",
			s:       `{"key":"value","num":42}`,
			want:    map[string]interface{}{"key": "value", "num": float64(42)},
			wantErr: false,
		},
		{
			name:    "empty string",
			s:       "",
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "valid empty object",
			s:    `{}`,
			want: map[string]interface{}{},
		},
		{
			name:    "complex nested structure",
			s:       `{"a":{"b":{"c":1}}}`,
			want:    map[string]interface{}{"a": map[string]interface{}{"b": map[string]interface{}{"c": float64(1)}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToMap(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSliceToString(t *testing.T) {
	tests := []struct {
		name    string
		s       []string
		want    string
		wantErr bool
	}{
		{
			name:    "non-empty slice",
			s:       []string{"item1", "item2", "item3"},
			want:    `["item1","item2","item3"]`,
			wantErr: false,
		},
		{
			name:    "nil slice",
			s:       nil,
			want:    `[]`,
			wantErr: false,
		},
		{
			name: "empty slice",
			s:    []string{},
			want: `[]`,
		},
		{
			name:    "slice with numbers",
			s:       []string{"1", "2", "3"},
			want:    `["1","2","3"]`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SliceToString(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("SliceToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringToSlice(t *testing.T) {
	tests := []struct {
		name    string
		s       string
		want    []string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			s:       `["item1","item2","item3"]`,
			want:    []string{"item1", "item2", "item3"},
			wantErr: false,
		},
		{
			name:    "empty string",
			s:       "",
			want:    []string{},
			wantErr: false,
		},
		{
			name: "valid empty array",
			s:    `[]`,
			want: []string{},
		},
		{
			name:    "array with numbers",
			s:       `["1","2","3"]`,
			want:    []string{"1", "2", "3"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StringToSlice(tt.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringToSlice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

// ========== Request/Response Structures Tests ==========

func TestCreateSessionRequest(t *testing.T) {
	req := &CreateSessionRequest{
		UserID: "user-1",
	}
	assert.Equal(t, "user-1", req.UserID)
}

func TestUpdateSessionRequest(t *testing.T) {
	tests := []struct {
		name  string
		req   *UpdateSessionRequest
		want  string
		valid bool
	}{
		{
			name:  "with UserID",
			req:   &UpdateSessionRequest{UserID: "user-1"},
			want:  "user-1",
			valid: true,
		},
		{
			name:  "empty UserID",
			req:   &UpdateSessionRequest{},
			want:  "",
			valid: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.req.UserID)
		})
	}
}

func TestSessionResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *SessionResponse
	}{
		{
			name: "success with session",
			resp: &SessionResponse{
				Success: true,
				Session: &SessionDTO{ID: "session-1"},
			},
		},
		{
			name: "success without session",
			resp: &SessionResponse{
				Success: true,
			},
		},
		{
			name: "failure",
			resp: &SessionResponse{
				Success: false,
				Error:   "Something went wrong",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Error != "" {
				assert.False(t, tt.resp.Success)
			} else {
				assert.True(t, tt.resp.Success)
			}
		})
	}
}

func TestSessionsResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *SessionsResponse
	}{
		{
			name: "success with sessions",
			resp: &SessionsResponse{
				Success:  true,
				Sessions: []*SessionDTO{{ID: "session-1"}, {ID: "session-2"}},
			},
		},
		{
			name: "success without sessions",
			resp: &SessionsResponse{
				Success: true,
			},
		},
		{
			name: "failure",
			resp: &SessionsResponse{
				Success: false,
				Error:   "Something went wrong",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Error != "" {
				assert.False(t, tt.resp.Success)
			} else {
				assert.True(t, tt.resp.Success)
			}
		})
	}
}

func TestCreateMessageRequest(t *testing.T) {
	req := &CreateMessageRequest{
		SessionID: "session-1",
		Role:      "user",
		Content:   "Hello",
	}
	assert.Equal(t, "session-1", req.SessionID)
	assert.Equal(t, "user", req.Role)
	assert.Equal(t, "Hello", req.Content)
}

func TestMessageResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *MessageResponse
	}{
		{
			name: "success with message",
			resp: &MessageResponse{
				Success: true,
				Message: &MessageDTO{ID: "msg-1"},
			},
		},
		{
			name: "success without message",
			resp: &MessageResponse{
				Success: true,
			},
		},
		{
			name: "failure",
			resp: &MessageResponse{
				Success: false,
				Error:   "Error message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Error != "" {
				assert.False(t, tt.resp.Success)
			} else {
				assert.True(t, tt.resp.Success)
			}
		})
	}
}

func TestMessagesResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *MessagesResponse
	}{
		{
			name: "success with messages",
			resp: &MessagesResponse{
				Success:  true,
				Messages: []*MessageDTO{{ID: "msg-1"}, {ID: "msg-2"}},
			},
		},
		{
			name: "success without messages",
			resp: &MessagesResponse{
				Success: true,
			},
		},
		{
			name: "failure",
			resp: &MessagesResponse{
				Success: false,
				Error:   "Error message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Error != "" {
				assert.False(t, tt.resp.Success)
			} else {
				assert.True(t, tt.resp.Success)
			}
		})
	}
}

func TestSendMessageRequest(t *testing.T) {
	req := &SendMessageRequest{
		UserID: "user-1",
		Message: ChatMessage{
			Role:    "user",
			Content: "Hello",
		},
		Options: MessageOptions{
			Model:     "gpt-4",
			MaxTokens: 100,
		},
	}
	assert.Equal(t, "user-1", req.UserID)
	assert.Equal(t, "user", req.Message.Role)
	assert.Equal(t, "Hello", req.Message.Content)
}

func TestSendMessageResponse(t *testing.T) {
	tests := []struct {
		name string
		resp *SendMessageResponse
	}{
		{
			name: "success with single message",
			resp: &SendMessageResponse{
				Success: true,
				Message: &MessageDTO{ID: "msg-1"},
			},
		},
		{
			name: "success with multiple messages",
			resp: &SendMessageResponse{
				Success:  true,
				Messages: []*MessageDTO{{ID: "msg-1"}, {ID: "msg-2"}},
			},
		},
		{
			name: "failure",
			resp: &SendMessageResponse{
				Success: false,
				Error:   "Error message",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.resp.Message != nil {
				assert.True(t, tt.resp.Success)
			} else if len(tt.resp.Messages) > 0 {
				assert.True(t, tt.resp.Success)
			} else if tt.resp.Error == "" {
				assert.False(t, tt.resp.Success)
			}
		})
	}
}

// ========== Response Helper Functions Tests ==========

// ========== Skill Response Helper Tests ==========

func TestErrorSkillResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "with error message",
			err:         fmt.Errorf("test error"),
			expectedErr: "operation failed: test error",
		},
		{
			name:        "with nil error",
			err:         nil,
			expectedErr: "operation failed: <nil>",
		},
		{
			name:        "with complex error",
			err:         fmt.Errorf("failed to connect: %w", fmt.Errorf("connection timeout")),
			expectedErr: "operation failed: failed to connect: connection timeout",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorSkillResponse(tt.err)
			assert.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestSuccessSkillResponse(t *testing.T) {
	skill := &SkillDTO{
		ID:          "skill-1",
		Name:        "TestSkill",
		Version:     "1.0.0",
		Location:    "/path/to/skill",
		Permissions: `["read", "write"]`,
		Metadata:    `{"key": "value"}`,
		CreatedAt:   getTestTime().Format(time.RFC3339),
	}

	resp := SuccessSkillResponse(skill)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Skill)
	assert.Equal(t, skill, resp.Skill)
}

// ========== Skills Response Helper Tests ==========

func TestErrorSkillsResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "with error message",
			err:         fmt.Errorf("test error"),
			expectedErr: "operation failed: test error",
		},
		{
			name:        "with nil error",
			err:         nil,
			expectedErr: "operation failed: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorSkillsResponse(tt.err)
			assert.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestSuccessSkillsResponse(t *testing.T) {
	skills := []*SkillDTO{
		{
			ID:        "skill-1",
			Name:      "Skill1",
			Version:   "1.0.0",
			Location:  "/path/to/skill1",
			CreatedAt: getTestTime().Format(time.RFC3339),
		},
		{
			ID:        "skill-2",
			Name:      "Skill2",
			Version:   "2.0.0",
			Location:  "/path/to/skill2",
			CreatedAt: getTestTime().Format(time.RFC3339),
		},
		{
			ID:        "skill-3",
			Name:      "Skill3",
			Version:   "1.5.0",
			Location:  "/path/to/skill3",
			CreatedAt: getTestTime().Format(time.RFC3339),
		},
	}

	resp := SuccessSkillsResponse(skills)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Skills)
	assert.Equal(t, 3, len(resp.Skills))
	assert.Equal(t, skills, resp.Skills)
}

func TestSuccessSkillsResponseEmpty(t *testing.T) {
	resp := SuccessSkillsResponse([]*SkillDTO{})
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Skills)
	assert.Equal(t, 0, len(resp.Skills))
}

func TestSuccessSkillsResponseNil(t *testing.T) {
	resp := SuccessSkillsResponse(nil)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	// Skills can be nil, which is fine
}

// ========== Schedule Response Helper Tests ==========

func TestErrorScheduleResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "with error message",
			err:         fmt.Errorf("test error"),
			expectedErr: "operation failed: test error",
		},
		{
			name:        "with nil error",
			err:         nil,
			expectedErr: "operation failed: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorScheduleResponse(tt.err)
			assert.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestSuccessScheduleResponse(t *testing.T) {
	schedule := &ScheduleDTO{
		ID:             "schedule-1",
		Skill:          "skill-1",
		CronExpression: "0 0 * * *",
		Input:          `{"key": "value"}`,
		Enabled:        true,
		CreatedAt:      getTestTime().Format(time.RFC3339),
	}

	resp := SuccessScheduleResponse(schedule)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Schedule)
	assert.Equal(t, schedule, resp.Schedule)
}

// ========== Schedules Response Helper Tests ==========

func TestErrorSchedulesResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "with error message",
			err:         fmt.Errorf("test error"),
			expectedErr: "operation failed: test error",
		},
		{
			name:        "with nil error",
			err:         nil,
			expectedErr: "operation failed: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorSchedulesResponse(tt.err)
			assert.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestSuccessSchedulesResponse(t *testing.T) {
	schedules := []*ScheduleDTO{
		{
			ID:             "schedule-1",
			Skill:          "skill-1",
			CronExpression: "0 0 * * *",
			Input:          `{"key": "value"}`,
			Enabled:        true,
			CreatedAt:      getTestTime().Format(time.RFC3339),
		},
		{
			ID:             "schedule-2",
			Skill:          "skill-2",
			CronExpression: "0 * * * *",
			Input:          `{"another": "param"}`,
			Enabled:        false,
			CreatedAt:      getTestTime().Format(time.RFC3339),
		},
	}

	resp := SuccessSchedulesResponse(schedules)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Schedules)
	assert.Equal(t, 2, len(resp.Schedules))
	assert.Equal(t, schedules, resp.Schedules)
}

func TestSuccessSchedulesResponseEmpty(t *testing.T) {
	resp := SuccessSchedulesResponse([]*ScheduleDTO{})
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.NotNil(t, resp.Schedules)
	assert.Equal(t, 0, len(resp.Schedules))
}

func TestSuccessSchedulesResponseNil(t *testing.T) {
	resp := SuccessSchedulesResponse(nil)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	// Schedules can be nil, which is fine
}

// ========== Skill Execution Response Helper Tests ==========

func TestErrorSkillExecutionResponse(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		expectedErr string
	}{
		{
			name:        "with error message",
			err:         fmt.Errorf("test error"),
			expectedErr: "operation failed: test error",
		},
		{
			name:        "with nil error",
			err:         nil,
			expectedErr: "operation failed: <nil>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorSkillExecutionResponse(tt.err)
			assert.NotNil(t, resp)
			assert.False(t, resp.Success)
			assert.Equal(t, tt.expectedErr, resp.Error)
		})
	}
}

func TestSuccessSkillExecutionResponse(t *testing.T) {
	tests := []struct {
		name   string
		output string
	}{
		{
			name:   "with output",
			output: "Success output from skill execution",
		},
		{
			name:   "with empty output",
			output: "",
		},
		{
			name:   "with multi-line output",
			output: "Line 1\nLine 2\nLine 3",
		},
		{
			name:   "with JSON output",
			output: `{"result": "success", "data": [1, 2, 3]}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := SuccessSkillExecutionResponse(tt.output)
			assert.NotNil(t, resp)
			assert.True(t, resp.Success)
			assert.Equal(t, tt.output, resp.Output)
		})
	}
}
