package dto

import (
	"encoding/json"
	"time"

	"github.com/atumaikin/nexflow/internal/domain/entity"
)

// ToEntity converts UserDTO to entity.User
func (dto *UserDTO) ToEntity() *entity.User {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	return &entity.User{
		ID:        dto.ID,
		Channel:   dto.Channel,
		ChannelID: dto.ChannelID,
		CreatedAt: createdAt,
	}
}

// FromEntity converts entity.User to UserDTO
func UserDTOFromEntity(user *entity.User) *UserDTO {
	return &UserDTO{
		ID:        user.ID,
		Channel:   user.Channel,
		ChannelID: user.ChannelID,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}
}

// ToEntity converts SessionDTO to entity.Session
func (dto *SessionDTO) ToEntity() *entity.Session {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dto.UpdatedAt)
	return &entity.Session{
		ID:        dto.ID,
		UserID:    dto.UserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromEntity converts entity.Session to SessionDTO
func SessionDTOFromEntity(session *entity.Session) *SessionDTO {
	return &SessionDTO{
		ID:        session.ID,
		UserID:    session.UserID,
		CreatedAt: session.CreatedAt.Format(time.RFC3339),
		UpdatedAt: session.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEntity converts MessageDTO to entity.Message
func (dto *MessageDTO) ToEntity() *entity.Message {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	return &entity.Message{
		ID:        dto.ID,
		SessionID: dto.SessionID,
		Role:      dto.Role,
		Content:   dto.Content,
		CreatedAt: createdAt,
	}
}

// FromEntity converts entity.Message to MessageDTO
func MessageDTOFromEntity(message *entity.Message) *MessageDTO {
	return &MessageDTO{
		ID:        message.ID,
		SessionID: message.SessionID,
		Role:      message.Role,
		Content:   message.Content,
		CreatedAt: message.CreatedAt.Format(time.RFC3339),
	}
}

// ToEntity converts TaskDTO to entity.Task
func (dto *TaskDTO) ToEntity() *entity.Task {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, dto.UpdatedAt)
	return &entity.Task{
		ID:        dto.ID,
		SessionID: dto.SessionID,
		Skill:     dto.Skill,
		Input:     dto.Input,
		Output:    dto.Output,
		Status:    dto.Status,
		Error:     dto.Error,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

// FromEntity converts entity.Task to TaskDTO
func TaskDTOFromEntity(task *entity.Task) *TaskDTO {
	return &TaskDTO{
		ID:        task.ID,
		SessionID: task.SessionID,
		Skill:     task.Skill,
		Input:     task.Input,
		Output:    task.Output,
		Status:    task.Status,
		Error:     task.Error,
		CreatedAt: task.CreatedAt.Format(time.RFC3339),
		UpdatedAt: task.UpdatedAt.Format(time.RFC3339),
	}
}

// ToEntity converts SkillDTO to entity.Skill
func (dto *SkillDTO) ToEntity() *entity.Skill {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	return &entity.Skill{
		ID:          dto.ID,
		Name:        dto.Name,
		Version:     dto.Version,
		Location:    dto.Location,
		Permissions: dto.Permissions,
		Metadata:    dto.Metadata,
		CreatedAt:   createdAt,
	}
}

// FromEntity converts entity.Skill to SkillDTO
func SkillDTOFromEntity(skill *entity.Skill) *SkillDTO {
	return &SkillDTO{
		ID:          skill.ID,
		Name:        skill.Name,
		Version:     skill.Version,
		Location:    skill.Location,
		Permissions: skill.Permissions,
		Metadata:    skill.Metadata,
		CreatedAt:   skill.CreatedAt.Format(time.RFC3339),
	}
}

// ToEntity converts ScheduleDTO to entity.Schedule
func (dto *ScheduleDTO) ToEntity() *entity.Schedule {
	createdAt, _ := time.Parse(time.RFC3339, dto.CreatedAt)
	return &entity.Schedule{
		ID:             dto.ID,
		Skill:          dto.Skill,
		CronExpression: dto.CronExpression,
		Input:          dto.Input,
		Enabled:        dto.Enabled,
		CreatedAt:      createdAt,
	}
}

// FromEntity converts entity.Schedule to ScheduleDTO
func ScheduleDTOFromEntity(schedule *entity.Schedule) *ScheduleDTO {
	return &ScheduleDTO{
		ID:             schedule.ID,
		Skill:          schedule.Skill,
		CronExpression: schedule.CronExpression,
		Input:          schedule.Input,
		Enabled:        schedule.Enabled,
		CreatedAt:      schedule.CreatedAt.Format(time.RFC3339),
	}
}

// MapToString converts map to JSON string
func MapToString(m map[string]interface{}) (string, error) {
	if m == nil {
		return "{}", nil
	}
	b, err := json.Marshal(m)
	if err != nil {
		return "{}", err
	}
	return string(b), nil
}

// StringToMap converts JSON string to map
func StringToMap(s string) (map[string]interface{}, error) {
	if s == "" {
		return make(map[string]interface{}), nil
	}
	var m map[string]interface{}
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		return make(map[string]interface{}), err
	}
	return m, nil
}

// SliceToString converts slice to JSON string
func SliceToString(s []string) (string, error) {
	if s == nil {
		return "[]", nil
	}
	b, err := json.Marshal(s)
	if err != nil {
		return "[]", err
	}
	return string(b), nil
}

// StringToSlice converts JSON string to slice
func StringToSlice(s string) ([]string, error) {
	if s == "" {
		return []string{}, nil
	}
	var slice []string
	err := json.Unmarshal([]byte(s), &slice)
	if err != nil {
		return []string{}, err
	}
	return slice, nil
}
