package common

import (
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
)

func GetEnv(key string, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}

	return value
}

func GenerateToken() string {
	return uuid.New().String()
}

func GetIsoString() string {
	return time.Now().Format(time.RFC3339)
}

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrorRecordNotFound = errors.New("record not found")
	ErrorEditConflict   = errors.New("edit conflict")
)
