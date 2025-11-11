package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"math"
	"time"
)

func main() {
	base32Demo()
	totp()
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

func totp() {
	keyBase32 := []byte("HFLEOZBUOVKXMVRY")
	key := make([]byte, 32)
	n, err := base32.StdEncoding.Decode(key, keyBase32)
	if err != nil {
		fmt.Println("decode base32 failed, error: ", n, err)
		return
	}
	key = key[:n]

	timestampSecond := time.Now().Unix()

	// sha-1(key, content)
	hasher := hmac.New(sha1.New, key)
	hasher.Write(itob(timestampSecond / 30))
	hmacHash := hasher.Sum(nil)

	// get an int32 from hash
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

func itob(integer int64) []byte {
	byteArr := make([]byte, 8)
	for i := 7; i >= 0; i-- {
		byteArr[i] = byte(integer & 0xff)
		integer = integer >> 8
	}
	return byteArr
}
