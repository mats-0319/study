package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"log"
	"os"
)

// generate P256 curve, keys use 'x509.marshal' to file bytes
func generateKeypair() {
	// generate private key
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Println("> Generate private key failed, error: ", err)
		return
	}

	// save private key
	privKeyBytes, err := x509.MarshalECPrivateKey(privKey)
	if err != nil {
		log.Println("> Save private key failed, error: ", err)
		return
	}

	err = os.WriteFile(privateKeyFilePath, privKeyBytes, 0777)
	if err != nil {
		log.Println("> Save private key failed, error: ", err)
		return
	}

	// save public key
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&privKey.PublicKey)
	if err != nil {
		log.Println("> Marshal public key failed, error: ", err)
		return
	}

	err = os.WriteFile(publicKeyFilePath, pubKeyBytes, 0777)
	if err != nil {
		log.Println("> Save public key failed, error: ", err)
		return
	}

	log.Println("> Generate key pair success.")
}
