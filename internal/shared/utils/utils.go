package utils

import (
	"encoding/json"
	"time"
)

// Now returns the current time in UTC
func Now() time.Time {
	return time.Now().UTC()
}

// FormatTimeRFC3339 formats a time.Time to RFC3339 string
func FormatTimeRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseTimeRFC3339 parses an RFC3339 string to time.Time
// Returns zero time if parsing fails
func ParseTimeRFC3339(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

// MarshalJSON marshals any value to JSON string
// Returns "{}" if marshaling fails
func MarshalJSON(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// UnmarshalJSONToMap unmarshals JSON string to map[string]interface{}
// Returns nil if unmarshaling fails
func UnmarshalJSONToMap(s string) map[string]interface{} {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil
	}
	return result
}

// UnmarshalJSONToSlice unmarshals JSON string to []string
// Returns nil if unmarshaling fails
func UnmarshalJSONToSlice(s string) []string {
	var result []string
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil
	}
	return result
}
