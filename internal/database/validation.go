package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidID       = errors.New("invalid id: must be a valid UUID")
	ErrInvalidDateTime = errors.New("invalid datetime: must be in RFC3339 format")
	ErrRequiredField   = errors.New("required field is empty")
)

// ValidateUser checks if User has valid data
func ValidateUser(u *User) error {
	if u.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(u.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, u.ID)
	}
	if u.Channel == "" {
		return fmt.Errorf("%w: channel", ErrRequiredField)
	}
	if u.ChannelUserID == "" {
		return fmt.Errorf("%w: channel_user_id", ErrRequiredField)
	}
	if u.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, u.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateSession checks if Session has valid data
func ValidateSession(s *Session) error {
	if s.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(s.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, s.ID)
	}
	if s.UserID == "" {
		return fmt.Errorf("%w: user_id", ErrRequiredField)
	}
	if _, err := uuid.Parse(s.UserID); err != nil {
		return fmt.Errorf("%w: user_id %s", ErrInvalidID, s.UserID)
	}
	if s.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, s.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	if s.UpdatedAt == "" {
		return fmt.Errorf("%w: updated_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, s.UpdatedAt); err != nil {
		return fmt.Errorf("%w: updated_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateMessage checks if Message has valid data
func ValidateMessage(m *Message) error {
	if m.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(m.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, m.ID)
	}
	if m.SessionID == "" {
		return fmt.Errorf("%w: session_id", ErrRequiredField)
	}
	if _, err := uuid.Parse(m.SessionID); err != nil {
		return fmt.Errorf("%w: session_id %s", ErrInvalidID, m.SessionID)
	}
	if m.Role == "" {
		return fmt.Errorf("%w: role", ErrRequiredField)
	}
	if m.Content == "" {
		return fmt.Errorf("%w: content", ErrRequiredField)
	}
	if m.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, m.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateTask checks if Task has valid data
func ValidateTask(t *Task) error {
	if t.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(t.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, t.ID)
	}
	if t.SessionID == "" {
		return fmt.Errorf("%w: session_id", ErrRequiredField)
	}
	if _, err := uuid.Parse(t.SessionID); err != nil {
		return fmt.Errorf("%w: session_id %s", ErrInvalidID, t.SessionID)
	}
	if t.Skill == "" {
		return fmt.Errorf("%w: skill", ErrRequiredField)
	}
	if t.Input == "" {
		return fmt.Errorf("%w: input", ErrRequiredField)
	}
	if t.Status == "" {
		return fmt.Errorf("%w: status", ErrRequiredField)
	}
	if t.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, t.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	if t.UpdatedAt == "" {
		return fmt.Errorf("%w: updated_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, t.UpdatedAt); err != nil {
		return fmt.Errorf("%w: updated_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateSkill checks if Skill has valid data
func ValidateSkill(s *Skill) error {
	if s.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(s.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, s.ID)
	}
	if s.Name == "" {
		return fmt.Errorf("%w: name", ErrRequiredField)
	}
	if s.Version == "" {
		return fmt.Errorf("%w: version", ErrRequiredField)
	}
	if s.Location == "" {
		return fmt.Errorf("%w: location", ErrRequiredField)
	}
	if s.Permissions == "" {
		return fmt.Errorf("%w: permissions", ErrRequiredField)
	}
	if s.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, s.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateSchedule checks if Schedule has valid data
func ValidateSchedule(s *Schedule) error {
	if s.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(s.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, s.ID)
	}
	if s.Skill == "" {
		return fmt.Errorf("%w: skill", ErrRequiredField)
	}
	if s.CronExpression == "" {
		return fmt.Errorf("%w: cron_expression", ErrRequiredField)
	}
	if s.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, s.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	return nil
}

// ValidateLog checks if Log has valid data
func ValidateLog(l *Log) error {
	if l.ID == "" {
		return fmt.Errorf("%w: id", ErrRequiredField)
	}
	if _, err := uuid.Parse(l.ID); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidID, l.ID)
	}
	if l.Level == "" {
		return fmt.Errorf("%w: level", ErrRequiredField)
	}
	if l.Source == "" {
		return fmt.Errorf("%w: source", ErrRequiredField)
	}
	if l.Message == "" {
		return fmt.Errorf("%w: message", ErrRequiredField)
	}
	if l.CreatedAt == "" {
		return fmt.Errorf("%w: created_at", ErrRequiredField)
	}
	if _, err := time.Parse(time.RFC3339, l.CreatedAt); err != nil {
		return fmt.Errorf("%w: created_at", ErrInvalidDateTime)
	}
	return nil
}
