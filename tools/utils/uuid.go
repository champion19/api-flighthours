package utils

import (
	"github.com/google/uuid"
)

// Generate creates a new UUID v4
func Generate() string {
	return uuid.New().String()
}

// IsValid checks if a string is a valid UUID
func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}

// Parse validates and parses a UUID string
func Parse(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

// MustParse parses a UUID string and panics on error
// Use only when you're certain the UUID is valid
func MustParse(id string) uuid.UUID {
	return uuid.MustParse(id)
}
