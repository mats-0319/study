package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func encrypt() {
	pubKey, err := getPubKey()
	if err != nil {
		ExecLog("Get public key", err)
		return
	}

	filePath, err := getFirstFile(plainTextFileName)
	if err != nil {
		ExecLog("Get first file", err)
		return
	}

	// read plain text file
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		ExecLog("Read plain text file", err)
		return
	}

	// encrypt file content
	encryptedBytes, err := ecies.Encrypt(rand.Reader, pubKey, fileBytes, nil, nil)
	if err != nil {
		ExecLog("Encrypt", err)
		return
	}

	// write cipher file
	extension := strings.ToUpper(getExtension(filePath, plainTextFileName))
	err = os.WriteFile(fmt.Sprintf("./%s%s", cipherFileName, extension), encryptedBytes, 0777)
	if err != nil {
		ExecLog("Write cipher file", err)
		return
	}

	ExecLog("Encrypt")
}

func getPubKey() (*ecies.PublicKey, error) {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		ExecLog("Read public key", err)
		return nil, err
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		ExecLog("Parse public key", err)
		return nil, err
	}

	pubKeyIns, ok := pubKeyI.(*ecdsa.PublicKey)
	if !ok {
		ExecLog("Public key type assert", err)
		return nil, err
	}

	return ecies.ImportECDSAPublic(pubKeyIns), nil
}
