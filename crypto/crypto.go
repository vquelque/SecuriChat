package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"github.com/vquelque/SecuriChat/message"
)

func RSAEncrypt(data []byte, destPublicKey *rsa.PublicKey) message.RSAEncryptedMessage {
	hash := sha256.New()
	label := []byte("")
	encryptedMessage, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		destPublicKey,
		data,
		label,
	)
	if err != nil {
		log.Println(len(data))
		log.Println(err.Error())
		log.Fatal("Error encrypting the RSA message")

	}
	return encryptedMessage
}

func RSADecrypt(cypher message.RSAEncryptedMessage, privateKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	label := []byte("")
	data, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		privateKey,
		cypher,
		label,
	)
	return data, err
}

func GenerateRSAKeypair() (priv *rsa.PrivateKey, pub *rsa.PublicKey) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Error generating the RSA keypair")
	}
	return private, &private.PublicKey
}

func PrintPublicKey(pub *rsa.PublicKey) []byte {
	asn1Key, _ := x509.MarshalPKIXPublicKey(pub)
	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: asn1Key,
	})
	fmt.Println(string(pubBytes))
	return pubBytes
}
