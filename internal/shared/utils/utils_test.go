package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	now := Now()

	// Should return UTC time
	assert.True(t, now.UTC().Equal(now))
	assert.False(t, time.Time{}.Equal(now))
}

func TestFormatTimeRFC3339(t *testing.T) {
	testTime := time.Date(2024, 1, 15, 12, 30, 45, 0, time.UTC)

	formatted := FormatTimeRFC3339(testTime)

	assert.Equal(t, "2024-01-15T12:30:45Z", formatted)
}

func TestParseTimeRFC3339_Valid(t *testing.T) {
	timeStr := "2024-01-15T12:30:45Z"
	parsed := ParseTimeRFC3339(timeStr)

	assert.False(t, parsed.IsZero())
	assert.Equal(t, 2024, parsed.Year())
	assert.Equal(t, time.January, parsed.Month())
	assert.Equal(t, 15, parsed.Day())
}

func TestParseTimeRFC3339_Invalid(t *testing.T) {
	timeStr := "invalid-time"
	parsed := ParseTimeRFC3339(timeStr)

	assert.True(t, parsed.IsZero())
}

func TestParseTimeRFC3339_Empty(t *testing.T) {
	parsed := ParseTimeRFC3339("")

	assert.True(t, parsed.IsZero())
}

func TestMarshalJSON_Valid(t *testing.T) {
	data := map[string]string{"key": "value"}

	result := MarshalJSON(data)

	assert.Contains(t, result, `"key"`)
	assert.Contains(t, result, `"value"`)
	assert.Contains(t, result, "{")
	assert.Contains(t, result, "}")
}

func TestMarshalJSON_Invalid(t *testing.T) {
	// Pass a value that can't be marshaled (e.g., a channel)
	ch := make(chan int)

	result := MarshalJSON(ch)

	// Should return empty object string on error
	assert.Equal(t, "{}", result)
}

func TestUnmarshalJSONToMap_Valid(t *testing.T) {
	jsonStr := `{"key":"value","number":42}`

	result := UnmarshalJSONToMap(jsonStr)

	assert.NotNil(t, result)
	assert.Equal(t, "value", result["key"])
	assert.Equal(t, float64(42), result["number"])
}

func TestUnmarshalJSONToMap_Invalid(t *testing.T) {
	jsonStr := `invalid json`

	result := UnmarshalJSONToMap(jsonStr)

	assert.Nil(t, result)
}

func TestUnmarshalJSONToMap_Empty(t *testing.T) {
	result := UnmarshalJSONToMap("{}")

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestUnmarshalJSONToSlice_Valid(t *testing.T) {
	jsonStr := `["item1","item2","item3"]`

	result := UnmarshalJSONToSlice(jsonStr)

	assert.NotNil(t, result)
	assert.Len(t, result, 3)
	assert.Equal(t, "item1", result[0])
	assert.Equal(t, "item2", result[1])
	assert.Equal(t, "item3", result[2])
}

func TestUnmarshalJSONToSlice_Invalid(t *testing.T) {
	jsonStr := `invalid json`

	result := UnmarshalJSONToSlice(jsonStr)

	assert.Nil(t, result)
}

func TestUnmarshalJSONToSlice_Empty(t *testing.T) {
	result := UnmarshalJSONToSlice("[]")

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestUnmarshalJSONToSlice_NonArray(t *testing.T) {
	// Valid JSON but not an array
	jsonStr := `{"key":"value"}`

	result := UnmarshalJSONToSlice(jsonStr)

	assert.Nil(t, result)
}
