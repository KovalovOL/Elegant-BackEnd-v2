package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)


func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generateRandomBytes(lenth_bytes int) ([]byte, error) {
	b := make([]byte, lenth_bytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create byte slice: %w", err)
	}
	return b, nil
}