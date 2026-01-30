package dto

import (
	"encoding/json"
)

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
