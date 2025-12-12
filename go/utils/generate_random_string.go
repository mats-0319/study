package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRandomString(length int) string {
	b := make([]byte, (length+1)/2)
	_, _ = rand.Read(b)

	return hex.EncodeToString(b)[:length]
}
