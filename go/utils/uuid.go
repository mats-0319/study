package utils

import (
	"crypto/sha256"

	"github.com/google/uuid"
)

// New return uuid v4 string,
// with same 'data', it will return same string,
// without 'data', it will return random string.
func New[T string | []byte](data ...T) string {
	if len(data) < 1 {
		return uuid.NewString()
	}

	return uuid.NewHash(sha256.New(), uuid.Nil, []byte(data[0]), 4).String()
}
