package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"math/rand/v2"
)

func main() {
	//otp()
	//base32Demo()
	//generateRandomBytes()
}

const randomCharactersLibrary = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const useBits = 6 // 6个bit位可以表示全部随机字符库中的字符

// generate random 20 Bytes
func generateRandomBytes() {
	b := make([]byte, 20)

	randomNum, remainBits := rand.Int64(), 64
	for i := 0; i < len(b); {
		if remainBits < useBits {
			randomNum, remainBits = rand.Int64(), 64
		}

		idx := int(randomNum & 0x3f) // 0b0011 1111
		if idx < len(randomCharactersLibrary) {
			randomNum >>= useBits
			remainBits -= useBits

			b[i] = randomCharactersLibrary[idx]
			i++
		} else {
			randomNum >>= 1
			remainBits -= 1
		}
	}

	fmt.Println(string(b))
}

func base32Demo() {
	randomBytes := []byte("a 20 bytes str      ")

	maxLen := base32.StdEncoding.EncodedLen(len(randomBytes))
	encodeBase32 := make([]byte, maxLen)
	base32.StdEncoding.Encode(encodeBase32, randomBytes)

	fmt.Println(string(encodeBase32), len(encodeBase32))

	maxLen = base32.StdEncoding.DecodedLen(len(encodeBase32))
	decodeBase32 := make([]byte, maxLen)
	n, err := base32.StdEncoding.Decode(decodeBase32, encodeBase32)
	if err != nil {
		fmt.Println("decode base32 failed, error: ", n, err)
		return
	}

	fmt.Println(string(decodeBase32), len(decodeBase32), n)

	fmt.Println(string(randomBytes) == string(decodeBase32[:n]))
}

func otp() {
	key := []byte("secret key")
	content := []byte("test text")

	// sha-1(key, content)
	hasher := hmac.New(sha1.New, key)
	hasher.Write(content)
	hmacHash := hasher.Sum(nil)

	// get a int32 from hash
	offset := int(hmacHash[len(hmacHash)-1] & 0x0f)
	// 算法要求屏蔽最高有效位
	longPassword := int(hmacHash[offset]&0x7f)<<24 |
		int(hmacHash[offset+1])<<16 |
		int(hmacHash[offset+2])<<8 |
		int(hmacHash[offset+3])

	// get 6 digits
	password := longPassword % int(math.Pow10(6))

	fmt.Println(fmt.Sprintf("%06d", password))
}
