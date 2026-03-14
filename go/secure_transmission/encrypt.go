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
		Error("Read plain text file", err)
		return
	}

	// do encrypt, ecdh + aes encrypt
	encryptedBytes := doEncrypt(pubKey, fileBytes)
	if encryptedBytes == nil {
		return
	}

	// write cipher file
	extension := strings.ToUpper(getExtension(filePath, plainTextFileName))
	err = os.WriteFile(fmt.Sprintf("./%s%s", cipherFileName, extension), encryptedBytes, 0777)
	if err != nil {
		Error("Write cipher file", err)
		return
	}

	Success("Encrypt")
}

func getPubKey() *ecdh.PublicKey {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		Error("Read public key", err)
		return nil
	}

	block, _ := pem.Decode(pubKeyBytes)
	if block == nil {
		Error("Decode public key", nil)
		return nil
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		Error("Parse public key", err)
		return nil
	}

	pubKey, ok := pubKeyI.(*ecdh.PublicKey)
	if !ok {
		ecdsaPubKey, ok := pubKeyI.(*ecdsa.PublicKey) // x509 parser usually return this type
		if !ok {
			Error("Public key type assert", nil)
			return nil
		}

		pubKey, err = ecdsaPubKey.ECDH()
		if err != nil {
			Error("ECDSA public key to ecdh", err)
			return nil
		}
	}

	return pubKey
}

func doEncrypt(pubKey *ecdh.PublicKey, content []byte) []byte {
	// generate final encrypt key
	tempPrivKey, err := generatePrivateKey()
	if err != nil {
		Error("Generate temp key", err)
		return nil
	}

	sharedKey, err := tempPrivKey.ECDH(pubKey)
	if err != nil {
		Error("ECDH", err)
		return nil
	}

	// aes-gcm encrypt
	var aesgcm cipher.AEAD
	{
		block, err := aes.NewCipher(sharedKey)
		if err != nil {
			Error("Build cipher block", err)
			return nil
		}

		aesgcm, err = cipher.NewGCM(block)
		if err != nil {
			Error("Build GCM", err)
			return nil
		}
	}

	nonce := make([]byte, aesgcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		Error("Generate random nonce", err)
		return nil
	}

	ciphertext := aesgcm.Seal(nil, nonce, content, nil)

	return encodeCipherFile(tempPrivKey.PublicKey().Bytes(), nonce, ciphertext)
}
