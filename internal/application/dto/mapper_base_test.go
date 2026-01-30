package dto

import (
	"testing"
	"time"
)

func TestParseTimeFields_Success(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantTime time.Time
		wantErr  bool
	}{
		{
			name:     "valid RFC3339 timestamp",
			input:    "2024-01-30T10:30:00Z",
			wantTime: time.Date(2024, 1, 30, 10, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "valid RFC3339 timestamp with offset",
			input:    "2024-01-30T12:30:00+02:00",
			wantTime: time.Date(2024, 1, 30, 10, 30, 0, 0, time.UTC),
			wantErr:  false,
		},
		{
			name:     "empty timestamp",
			input:    "",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "invalid timestamp",
			input:    "2024-01-30",
			wantTime: time.Time{},
			wantErr:  true,
		},
		{
			name:     "malformed timestamp",
			input:    "not-a-timestamp",
			wantTime: time.Time{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTime, err := ParseTimeFields(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTimeFields() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !gotTime.Equal(tt.wantTime) {
				t.Errorf("ParseTimeFields() = %v, want %v", gotTime, tt.wantTime)
			}
		})
	}
}

func TestMustParseTimeFields_Success(t *testing.T) {
	input := "2024-01-30T10:30:00Z"
	wantTime := time.Date(2024, 1, 30, 10, 30, 0, 0, time.UTC)

	gotTime := MustParseTimeFields(input)

	if !gotTime.Equal(wantTime) {
		t.Errorf("MustParseTimeFields() = %v, want %v", gotTime, wantTime)
	}
}

func TestMustParseTimeFields_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParseTimeFields() did not panic with invalid input")
		}
	}()

	MustParseTimeFields("") // Should panic
}

func TestMustParseTimeFields_PanicOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParseTimeFields() did not panic with invalid timestamp")
		}
	}()

	MustParseTimeFields("not-a-timestamp") // Should panic
}

func TestParseTimeFieldsWithUpdatedAt_Success(t *testing.T) {
	createdAtStr := "2024-01-30T10:00:00Z"
	updatedAtStr := "2024-01-30T11:00:00Z"

	wantCreatedAt := time.Date(2024, 1, 30, 10, 0, 0, 0, time.UTC)
	wantUpdatedAt := time.Date(2024, 1, 30, 11, 0, 0, 0, time.UTC)

	gotCreatedAt, gotUpdatedAt, err := ParseTimeFieldsWithUpdatedAt(createdAtStr, updatedAtStr)

	if err != nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() error = %v", err)
		return
	}

	if !gotCreatedAt.Equal(wantCreatedAt) {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() createdAt = %v, want %v", gotCreatedAt, wantCreatedAt)
	}

	if !gotUpdatedAt.Equal(wantUpdatedAt) {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() updatedAt = %v, want %v", gotUpdatedAt, wantUpdatedAt)
	}
}

func TestParseTimeFieldsWithUpdatedAt_EmptyCreatedAt(t *testing.T) {
	_, _, err := ParseTimeFieldsWithUpdatedAt("", "2024-01-30T11:00:00Z")

	if err == nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() expected error for empty created_at")
	}
}

func TestParseTimeFieldsWithUpdatedAt_EmptyUpdatedAt(t *testing.T) {
	_, _, err := ParseTimeFieldsWithUpdatedAt("2024-01-30T10:00:00Z", "")

	if err == nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() expected error for empty updated_at")
	}
}

func TestParseTimeFieldsWithUpdatedAt_InvalidCreatedAt(t *testing.T) {
	_, _, err := ParseTimeFieldsWithUpdatedAt("not-a-timestamp", "2024-01-30T11:00:00Z")

	if err == nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() expected error for invalid created_at")
	}
}

func TestParseTimeFieldsWithUpdatedAt_InvalidUpdatedAt(t *testing.T) {
	_, _, err := ParseTimeFieldsWithUpdatedAt("2024-01-30T10:00:00Z", "not-a-timestamp")

	if err == nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() expected error for invalid updated_at")
	}
}

func TestMustParseTimeFieldsWithUpdatedAt_Success(t *testing.T) {
	createdAtStr := "2024-01-30T10:00:00Z"
	updatedAtStr := "2024-01-30T11:00:00Z"

	wantCreatedAt := time.Date(2024, 1, 30, 10, 0, 0, 0, time.UTC)
	wantUpdatedAt := time.Date(2024, 1, 30, 11, 0, 0, 0, time.UTC)

	gotCreatedAt, gotUpdatedAt := MustParseTimeFieldsWithUpdatedAt(createdAtStr, updatedAtStr)

	if !gotCreatedAt.Equal(wantCreatedAt) {
		t.Errorf("MustParseTimeFieldsWithUpdatedAt() createdAt = %v, want %v", gotCreatedAt, wantCreatedAt)
	}

	if !gotUpdatedAt.Equal(wantUpdatedAt) {
		t.Errorf("MustParseTimeFieldsWithUpdatedAt() updatedAt = %v, want %v", gotUpdatedAt, wantUpdatedAt)
	}
}

func TestMustParseTimeFieldsWithUpdatedAt_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParseTimeFieldsWithUpdatedAt() did not panic with invalid input")
		}
	}()

	MustParseTimeFieldsWithUpdatedAt("", "") // Should panic
}

func TestFormatTimeFields(t *testing.T) {
	inputTime := time.Date(2024, 1, 30, 10, 30, 45, 123456789, time.UTC)
	wantOutput := "2024-01-30T10:30:45Z" // RFC3339Nano truncates to seconds in Go's Format

	gotOutput := FormatTimeFields(inputTime)

	if gotOutput != wantOutput {
		t.Errorf("FormatTimeFields() = %v, want %v", gotOutput, wantOutput)
	}
}

func TestFormatTimeFieldsWithUpdatedAt(t *testing.T) {
	createdAt := time.Date(2024, 1, 30, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2024, 1, 30, 11, 30, 45, 0, time.UTC)

	wantCreatedAt := "2024-01-30T10:00:00Z"
	wantUpdatedAt := "2024-01-30T11:30:45Z"

	gotCreatedAt, gotUpdatedAt := FormatTimeFieldsWithUpdatedAt(createdAt, updatedAt)

	if gotCreatedAt != wantCreatedAt {
		t.Errorf("FormatTimeFieldsWithUpdatedAt() createdAt = %v, want %v", gotCreatedAt, wantCreatedAt)
	}

	if gotUpdatedAt != wantUpdatedAt {
		t.Errorf("FormatTimeFieldsWithUpdatedAt() updatedAt = %v, want %v", gotUpdatedAt, wantUpdatedAt)
	}
}

// Test that helper functions maintain RFC3339 consistency
func TestRFC3339Consistency(t *testing.T) {
	originalTime := time.Date(2024, 1, 30, 10, 30, 0, 0, time.UTC)

	// Format and then parse back
	formatted := FormatTimeFields(originalTime)
	parsed, err := ParseTimeFields(formatted)

	if err != nil {
		t.Errorf("ParseTimeFields() error = %v", err)
		return
	}

	if !parsed.Equal(originalTime) {
		t.Errorf("Format then Parse: got %v, want %v", parsed, originalTime)
	}
}

func TestRFC3339ConsistencyWithUpdatedAt(t *testing.T) {
	originalCreatedAt := time.Date(2024, 1, 30, 10, 0, 0, 0, time.UTC)
	originalUpdatedAt := time.Date(2024, 1, 30, 11, 30, 0, 0, time.UTC)

	// Format and then parse back
	formattedCreatedAt, formattedUpdatedAt := FormatTimeFieldsWithUpdatedAt(originalCreatedAt, originalUpdatedAt)
	parsedCreatedAt, parsedUpdatedAt, err := ParseTimeFieldsWithUpdatedAt(formattedCreatedAt, formattedUpdatedAt)

	if err != nil {
		t.Errorf("ParseTimeFieldsWithUpdatedAt() error = %v", err)
		return
	}

	if !parsedCreatedAt.Equal(originalCreatedAt) {
		t.Errorf("Format then Parse createdAt: got %v, want %v", parsedCreatedAt, originalCreatedAt)
	}

	if !parsedUpdatedAt.Equal(originalUpdatedAt) {
		t.Errorf("Format then Parse updatedAt: got %v, want %v", parsedUpdatedAt, originalUpdatedAt)
	}
}
