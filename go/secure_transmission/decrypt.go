package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
)

func decrypt() {
	privKey := getPrivKey()
	if privKey == nil {
		return
	}

	// read cipher file
	filePath := getFirstFile(cipherFileName)
	if len(filePath) < 1 {
		return
	}

	cipherBytes, err := os.ReadFile(filePath)
	if err != nil {
		Error("Read cipher file", err)
		return
	}

	// do decrypt
	fileBytes := doDecrypt(privKey, cipherBytes)
	if fileBytes == nil {
		return
	}

	// write plain text file (decrypted)
	extension := strings.ToLower(getExtension(filePath, cipherFileName))
	err = os.WriteFile(fmt.Sprintf("./%s%s", plainTextDecryptedFileName, extension), fileBytes, 0777)
	if err != nil {
		Error("Write message", err)
		return
	}

	Success("Decrypt")
}

func getPrivKey() *ecdh.PrivateKey {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		Error("Read private key", err)
		return nil
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		Error("Decode private key", nil)
		return nil
	}

	privKeyI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Error("Parse private key", err)
		return nil
	}

	privKey, ok := privKeyI.(*ecdh.PrivateKey)
	if !ok {
		ecdsaPrivKey, ok := privKeyI.(*ecdsa.PrivateKey)
		if !ok {
			Error("Private key type assert", nil)
			return nil
		}

		privKey, err = ecdsaPrivKey.ECDH()
		if err != nil {
			Error("ECDSA private key to ecdh", err)
			return nil
		}
	}

	return privKey
}

func doDecrypt(privKey *ecdh.PrivateKey, fileBytes []byte) []byte {
	if len(fileBytes) < 1 {
		Error("Invalid cipher", nil)
		return nil
	}

	pubKeyBytes, nonce, ciphertext, err := decodeCipherFile(fileBytes)
	if err != nil {
		return nil
	}

	// generate final decrypt key
	tempPubKey, err := ecdh.P256().NewPublicKey(pubKeyBytes)
	if err != nil {
		Error("Load public key", err)
		return nil
	}

	sharedKey, err := privKey.ECDH(tempPubKey)
	if err != nil {
		Error("ECDH", err)
		return nil
	}

	// aes-gcm decrypt
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

	message, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		Error("Decrypt", err)
		return nil
	}

	return message
}
