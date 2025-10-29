package main

import (
	"crypto/x509"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func decrypt() {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		log.Println("> Read private key failed, error: ", err)
		return
	}

	privKeyECDSA, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		log.Println("> Parse private key failed, error: ", err)
		return
	}

	privKey := ecies.ImportECDSA(privKeyECDSA)

	cipherBytes, err := os.ReadFile(cipherTextFilePath)
	if err != nil {
		log.Println("> Read cipher text failed, error: ", err)
		return
	}

	messageBytes, err := privKey.Decrypt(cipherBytes, nil, nil)
	if err != nil {
		log.Println("> Decrypt failed, error: ", err)
		return
	}

	err = os.WriteFile(plainTextDecryptedFilePath, messageBytes, 0777)
	if err != nil {
		log.Println("> Write message failed, error: ", err)
		return
	}

	log.Println("> Decrypt success.")
}
