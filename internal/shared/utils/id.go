package utils

import "github.com/google/uuid"

// GenerateID generates a new UUID-based ID
func GenerateID() string {
	return uuid.New().String()
}
