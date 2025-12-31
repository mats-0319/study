package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"os"
)

// generate P256 curve, keys use 'x509.marshal' to file bytes
func generateKeypair() {
	// generate private key
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		ExecLog("Generate private key", err)
		return
	}

	// save private key
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		ExecLog("Marshal private key", err)
		return
	}

	err = os.WriteFile(privateKeyFilePath, privKeyBytes, 0777)
	if err != nil {
		ExecLog("Save private key", err)
		return
	}

	// save public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		ExecLog("Marshal public key", err)
		return
	}

	err = os.WriteFile(publicKeyFilePath, pubKeyBytes, 0777)
	if err != nil {
		ExecLog("Save public key", err)
		return
	}

	ExecLog("Generate key pair")
}
