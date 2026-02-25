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
		Log("Read cipher file", err)
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
		Log("Write message", err)
		return
	}

	Log("Decrypt")
}

func getPrivKey() *ecdh.PrivateKey {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		Log("Read private key", err)
		return nil
	}

	block, _ := pem.Decode(privKeyBytes)
	if block == nil {
		Log("Decode private key", emptyErr)
		return nil
	}

	privKeyI, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		Log("Parse private key", err)
		return nil
	}

	privKey, ok := privKeyI.(*ecdh.PrivateKey)
	if !ok {
		ecdsaPrivKey, ok := privKeyI.(*ecdsa.PrivateKey)
		if !ok {
			Log("Private key type assert", emptyErr)
			return nil
		}

		privKey, err = ecdsaPrivKey.ECDH()
		if err != nil {
			Log("ECDSA private key to ecdh", err)
			return nil
		}
	}

	return privKey
}

func doDecrypt(privKey *ecdh.PrivateKey, data []byte) []byte {
	if len(data) < 1 {
		Log("Invalid cipher", emptyErr)
		return nil
	}

	tempPubKeyLength := int(data[0])
	tempPubKeyBytes := data[1 : 1+tempPubKeyLength]
	tempPubKey, err := ecdh.P256().NewPublicKey(tempPubKeyBytes)
	if err != nil {
		Log("New public key", err)
		return nil
	}

	sharedKey, err := privKey.ECDH(tempPubKey)
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

	nonceSize := aesgcm.NonceSize()

	nonce := data[1+tempPubKeyLength : 1+tempPubKeyLength+nonceSize]
	cipherText := data[1+tempPubKeyLength+nonceSize:]

	message, err := aesgcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		Log("Open ciphertext", err)
		return nil
	}

	return message
}
