package auth

import (
	"crypto/rand"
	"encoding/hex"
)


func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}