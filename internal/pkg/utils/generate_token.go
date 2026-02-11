package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateRandomToken creates a secure random string of a specific byte length
func GenerateRandomToken(length int) (string, error) {
	// 1. Create a byte slice of the desired length
	b := make([]byte, length)

	// 2. Fill the slice with random bytes from the crypto/rand reader
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// 3. Encode the bytes into a string
	// URLEncoding avoids characters like '+' and '/' which can cause issues in APIs
	return base64.URLEncoding.EncodeToString(b), nil
}
