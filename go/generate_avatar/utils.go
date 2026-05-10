package generate_avatar

import (
	"crypto/sha256"
	"encoding/hex"
)

// calcSHA256 returns a 64-length hex string
func calcSHA256(text string) string {
	hasher := sha256.New()
	hasher.Write([]byte(text))
	hashBytes := hasher.Sum(nil)

	return hex.EncodeToString(hashBytes)
}
