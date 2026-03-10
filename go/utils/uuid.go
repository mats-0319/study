package utils

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

// NewV4 return uuid v4 string,
// with same 'data', it will return same string,
// without 'data', it will return random string.
func NewV4[T string | []byte](data ...T) string {
	if len(data) < 1 {
		return uuid.NewString()
	}

	return uuid.NewHash(sha256.New(), uuid.Nil, []byte(data[0]), 4).String()
}

func NewV7() (uuid.UUID, error) {
	return uuid.NewV7()
}
