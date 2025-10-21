package util

import "github.com/google/uuid"

// RandomUUID generate string of UUIDv4
func RandomUUID() string {
	return uuid.New().String()
}
