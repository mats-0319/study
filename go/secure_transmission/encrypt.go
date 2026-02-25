package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"
)

func encrypt() {
	pubKey := getPubKey()
	if pubKey == nil {
		return
	}

	// read plain text file
	filePath := getFirstFile(plainTextFileName)
	if len(filePath) < 1 {
		return
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		Log("Read plain text file", err)
		return
	}

	// do encrypt
	encryptedBytes := doEncrypt(pubKey, fileBytes)
	if encryptedBytes == nil {
		return
	}

	// write cipher file
	extension := strings.ToUpper(getExtension(filePath, plainTextFileName))
	err = os.WriteFile(fmt.Sprintf("./%s%s", cipherFileName, extension), encryptedBytes, 0777)
	if err != nil {
		Log("Write cipher file", err)
		return
	}

	Log("Encrypt")
}

func getPubKey() *ecdh.PublicKey {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		Log("Read public key", err)
		return nil
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		Log("Decode public key", emptyErr)
		return nil
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Log("Parse public key", err)
		return nil
	}

	pubKey, ok := pubKeyI.(*ecdh.PublicKey)
	if !ok {
		ecdsaPubKey, ok := pubKeyI.(*ecdsa.PublicKey)
		if !ok {
			Log("Public key type assert", emptyErr)
			return nil
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			Log("ECDSA public key to ecdh", err)
			return nil
		}
	}

	return pubKey
}

func doEncrypt(pubKey *ecdh.PublicKey, content []byte) []byte {
	tempPrivKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		Log("Generate temp key", err)
		return nil
	}

	sharedKey, err := tempPrivKey.ECDH(pubKey)
	if err != nil {
		Log("Generate shared key", err)
		return nil
	}

	block, err := aes.NewCipher(sharedKey)
	if err != nil {
		Log("Build cipher block", err)
		return nil
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		Log("Build GCM", err)
		return nil
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		Log("Read nonce", err)
		return nil
	}

	ciphertext := aesgcm.Seal(nil, nonce, content, nil)

	tempPubKeyBytes := tempPrivKey.PublicKey().Bytes()
	result := append([]byte{byte(len(tempPubKeyBytes))}, tempPubKeyBytes...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result
}
