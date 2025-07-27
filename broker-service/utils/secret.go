package utils

import (
	"log"
	"os"
)

var Secret []byte

func InitializeSecret() {
	if Secret != nil {
		return
	}

	secretStr := os.Getenv("SECRET")
	if secretStr == "" {
		log.Fatal("SECRET environment variable is required")
	}
	Secret = []byte(secretStr)
}
