package main

import (
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/crypto/ecies"
)

func decrypt() {
	privKey, err := getPrivKey()
	if err != nil {
		ExecLog("Get private key", err)
		return
	}

	filePath, err := getFirstFile(cipherFileName)
	if err != nil {
		ExecLog("Get first file", err)
		return
	}

	// read cipher file
	cipherBytes, err := os.ReadFile(filePath)
	if err != nil {
		ExecLog("Read cipher file", err)
		return
	}

	// decrypt file content
	fileBytes, err := privKey.Decrypt(cipherBytes, nil, nil)
	if err != nil {
		log.Println("> Decrypt failed, error: ", err)
		return
	}

	// write plain text file (decrypted)
	extension := strings.ToLower(getExtension(filePath, cipherFileName))
	err = os.WriteFile(fmt.Sprintf("./%s%s", plainTextDecryptedFileName, extension), fileBytes, 0777)
	if err != nil {
		log.Println("> Write message failed, error: ", err)
		return
	}

	log.Println("> Decrypt success.")
}

func getPrivKey() (*ecies.PrivateKey, error) {
	privKeyBytes, err := os.ReadFile(privateKeyFilePath)
	if err != nil {
		ExecLog("Read private key", err)
		return nil, err
	}

	privKeyECDSA, err := x509.ParseECPrivateKey(privKeyBytes)
	if err != nil {
		ExecLog("Parse private key", err)
		return nil, err
	}

	return ecies.ImportECDSA(privKeyECDSA), nil
}
