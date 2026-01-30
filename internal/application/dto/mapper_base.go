package dto

import (
	"fmt"
	"time"
)

// ParseTimeFields parses a single timestamp string to time.Time
// Returns zero time and error if parsing fails
func ParseTimeFields(createdAtStr string) (time.Time, error) {
	if createdAtStr == "" {
		return time.Time{}, fmt.Errorf("created_at timestamp is empty")
	}

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse created_at timestamp '%s': %w", createdAtStr, err)
	}

	return createdAt, nil
}

// MustParseTimeFields parses a single timestamp string to time.Time
// Panics if parsing fails. Use this when the timestamp is guaranteed to be valid.
func MustParseTimeFields(createdAtStr string) time.Time {
	t, err := ParseTimeFields(createdAtStr)
	if err != nil {
		panic(err)
	}
	return t
}

// ParseTimeFieldsWithUpdatedAt parses two timestamp strings to time.Time
// Returns zero times and error if parsing fails
func ParseTimeFieldsWithUpdatedAt(createdAtStr, updatedAtStr string) (time.Time, time.Time, error) {
	createdAt, err := ParseTimeFields(createdAtStr)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	if updatedAtStr == "" {
		return time.Time{}, time.Time{}, fmt.Errorf("updated_at timestamp is empty")
	}

	updatedAt, err := time.Parse(time.RFC3339, updatedAtStr)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("failed to parse updated_at timestamp '%s': %w", updatedAtStr, err)
	}

	return createdAt, updatedAt, nil
}

// MustParseTimeFieldsWithUpdatedAt parses two timestamp strings to time.Time
// Panics if parsing fails. Use this when the timestamps are guaranteed to be valid.
func MustParseTimeFieldsWithUpdatedAt(createdAtStr, updatedAtStr string) (time.Time, time.Time) {
	createdAt, updatedAt, err := ParseTimeFieldsWithUpdatedAt(createdAtStr, updatedAtStr)
	if err != nil {
		panic(err)
	}
	return createdAt, updatedAt
}

// FormatTimeFields formats a single time.Time to RFC3339 string
func FormatTimeFields(createdAt time.Time) string {
	return createdAt.Format(time.RFC3339)
}

// FormatTimeFieldsWithUpdatedAt formats two time.Time to RFC3339 strings
func FormatTimeFieldsWithUpdatedAt(createdAt, updatedAt time.Time) (string, string) {
	return createdAt.Format(time.RFC3339), updatedAt.Format(time.RFC3339)
}
