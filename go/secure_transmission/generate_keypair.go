package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
)

// generate P256 curve, keys use 'x509.marshal' to file bytes
func generateKeypair() {
	// generate private key
	privKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		Log("Generate private key", err)
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

	Log("Generate key pair")
}

func savePrivateKey(privKey *ecdh.PrivateKey) error {
	privKeyBytes, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		Log("Marshal private key", err)
		return err
	}

	block := &pem.Block{Type: "Private Key", Bytes: privKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(privateKeyFilePath, blockBytes, 0600)
	if err != nil {
		Log("Save private key", err)
		return err
	}

	return nil
}

func savePublicKey(pubKey *ecdh.PublicKey) error {
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		Log("Marshal public key", err)
		return err
	}

	block := &pem.Block{Type: "Public Key", Bytes: pubKeyBytes}
	blockBytes := pem.EncodeToMemory(block)

	err = os.WriteFile(publicKeyFilePath, blockBytes, 0644)
	if err != nil {
		Log("Save public key", err)
		return err
	}

	return nil
}
