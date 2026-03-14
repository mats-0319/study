package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"errors"
	"log"
	"os"
	"strings"
)

func Success(behavior string) {
	log.Printf("> %s success.\n", behavior)
}

func Info(behavior string) {
	log.Printf("> %s .\n", behavior)
}

func Error(behavior string, err error) {
	log.Printf("> %s failed, error: %v\n", behavior, err)
}

func encodeCipherFile(pubKeyBytes []byte, nonce []byte, ciphertext []byte) []byte {
	result := append([]byte{byte(len(pubKeyBytes))}, pubKeyBytes...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result
}

func decodeCipherFile(fileBytes []byte) (pubKeyBytes []byte, nonce []byte, ciphertext []byte, err error) {
	if len(fileBytes) < 1 {
		err = errors.New("invalid cipher")
		Error("Invalid cipher", err)
		return
	}

	pubkLength := int(fileBytes[0])
	pubKeyBytes = fileBytes[1 : 1+pubkLength]

	var nonceSize int
	{
		var block cipher.Block
		block, err = aes.NewCipher(make([]byte, 32))
		if err != nil {
			Error("Build cipher block", err)
			return
		}

		var aesgcm cipher.AEAD
		aesgcm, err = cipher.NewGCM(block)
		if err != nil {
			Error("Build GCM", err)
			return
		}

		nonceSize = aesgcm.NonceSize()
	}

	nonce = fileBytes[1+pubkLength : 1+pubkLength+nonceSize]
	ciphertext = fileBytes[1+pubkLength+nonceSize:]

	return
}

func generatePrivateKey() (*ecdh.PrivateKey, error) {
	// in go 1.26, it will no longer use input io.reader in 'crypto',
	// it will use a global random origin
	return ecdh.P256().GenerateKey(nil)
}

// getFirstFile return full file path matched 'fileName' without extension
func getFirstFile(fileName string) string {
	entry, err := os.ReadDir("./")
	if err != nil {
		Error("Read dir", err)
		return ""
	}

	var filePath strings.Builder
	filePath.WriteString("./")
	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		var fileInfo os.FileInfo
		fileInfo, err = entry[i].Info()
		if err != nil {
			Error("Read info", err)
			continue
		}

		if !strings.HasPrefix(fileInfo.Name()+".", fileName) {
			continue // ignore files with wrong name
		}

		filePath.WriteString(fileInfo.Name())
		break
	}

	return filePath.String()
}

func getExtension(filePath string, fileName string) string {
	return strings.TrimPrefix(filePath, "./"+fileName)
}
