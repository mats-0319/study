package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/x509"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func encrypt() {
	pubKeyBytes, err := os.ReadFile(publicKeyFilePath)
	if err != nil {
		log.Println("> Read public key failed, err: ", err)
		return
	}

	pubKeyI, err := x509.ParsePKIXPublicKey(pubKeyBytes)
	if err != nil {
		log.Println("> Parse public key failed, err: ", err)
		return
	}

	pubKeyIns, ok := pubKeyI.(*ecdsa.PublicKey)
	if !ok {
		log.Println("> Parse public key failed, err: ", err)
		return
	}
	pubKey := ecies.ImportECDSAPublic(pubKeyIns)

	messageBytes, err := os.ReadFile(plainTextFilePath)
	if err != nil {
		log.Println("> Read plain text failed, err: ", err)
		return
	}

	cipherBytes, err := ecies.Encrypt(rand.Reader, pubKey, messageBytes, nil, nil)
	if err != nil {
		log.Println("> Encrypt failed, err: ", err)
		return
	}

	err = os.WriteFile(cipherTextFilePath, cipherBytes, 0777)
	if err != nil {
		log.Println("> Write cipher text failed, err: ", err)
		return
	}

	log.Println("> Encrypt success.")
}
