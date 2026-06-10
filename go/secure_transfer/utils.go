package main

import (
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

func encodeCipherFile(pubKeyBytes []byte, ciphertext []byte) []byte {
	result := append([]byte{byte(len(pubKeyBytes))}, pubKeyBytes...)
	result = append(result, ciphertext...)

	return result
}

func decodeCipherFile(fileBytes []byte) (pubKeyBytes []byte, ciphertext []byte, err error) {
	if len(fileBytes) < 1 {
		err = errors.New("invalid cipher")
		Error("Invalid cipher", err)
		return
	}

	pubkLength := int(fileBytes[0])
	pubKeyBytes = fileBytes[1 : 1+pubkLength]
	ciphertext = fileBytes[1+pubkLength:]

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
