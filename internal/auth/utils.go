package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)


func GenerateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateRandomBytes(lenth_bytes int) ([]byte, error) {
	b := make([]byte, lenth_bytes)
	_, err := rand.Read(b)
	if err != nil {
		return nil, fmt.Errorf("failed to create byte slice: %w", err)
	}
	return b, nil
}

func HashBytes(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}