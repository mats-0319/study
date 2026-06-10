package main

import (
	"crypto/ecdh"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func generateKeypair() {
	privKey, err := generatePrivateKey()
	if err != nil {
		Error("Generate private key", err)
		return
	}

	err = savePrivateKey(privKey)
	if err != nil {
		return
	}

	err = savePublicKey(privKey.PublicKey())
	if err != nil {
		return
	}

	Success("Generate key pair")
}

func savePrivateKey(privKey *ecdh.PrivateKey) error {
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		Error("Marshal private key", err)
		return err
	}

	block := &pem.Block{Type: "Private Key", Bytes: privKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(privateKeyFilePath, blockBytes, 0600)
	if err != nil {
		Error("Save private key", err)
		return err
	}

	return nil
}

func savePublicKey(pubKey *ecdh.PublicKey) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		Error("Marshal public key", err)
		return err
	}

	block := &pem.Block{Type: "Public Key", Bytes: pubKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(publicKeyFilePath, blockBytes, 0644)
	if err != nil {
		Error("Save public key", err)
		return err
	}

	return nil
}
