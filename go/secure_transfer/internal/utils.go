package internal

import (
	"crypto/ecdh"
	"log"
	"os"
	"strings"
)

const (
	privateKeyFilePath         = "./priv.key"
	publicKeyFilePath          = "./PUB.KEY"
	plainTextFileName          = "message"
	cipherFileName             = "CIPHER"
	plainTextDecryptedFileName = "message_decrypted"
	defaultExtension           = ".txt"

	salt            = "mats0319"
	info            = "secure_transfer"
	publicKeyLength = 32
	driveKeyLength  = 32
)

func Success(behavior string) {
	log.Printf("> %s success.\n\n", behavior)
}

func Info(behavior string) {
	log.Printf("> %s\n", behavior)
}

func Error(behavior string, err error) {
	log.Printf("%s failed, error: %v\n\n", behavior, err)
}

// Curve use same curve unified
func Curve() ecdh.Curve {
	return ecdh.X25519()
}

func GetFirstFile(fileName string) (string, []byte) {
	entry, err := os.ReadDir("./")
	if err != nil {
		Error("Read dir", err)
		return "", nil
	}

	var filePath strings.Builder
	filePath.WriteString("./")

	for i := range entry {
		if entry[i].IsDir() {
			continue // ignore folder
		}

		fileInfo, err := entry[i].Info()
		if err != nil {
			Error("Read info", err)
			continue
		}

		if strings.HasPrefix(fileInfo.Name(), fileName) { // match file name without extension
			filePath.WriteString(fileInfo.Name())
			break
		}
	}

	fileBytes, err := os.ReadFile(filePath.String())
	if err != nil {
		Error("Read file", err)
		return "", nil
	}

	return filePath.String(), fileBytes
}

func GetExtension(filePath string, fileName string) string {
	return strings.TrimPrefix(filePath, "./"+fileName)
}
