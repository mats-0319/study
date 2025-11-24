package utils

import "math/rand/v2"

const CharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部字符库中的字符

// GenerateRandomBytes generate random 'length' readable Bytes
func GenerateRandomBytes(length int) []byte {
	b := make([]byte, length)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		index := int(randomNum & (1<<useBits - 1)) // 0b0011 1111
		if index < len(CharactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = CharactersLibrary[index]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	return b
}
