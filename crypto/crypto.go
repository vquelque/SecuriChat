package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
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
		log.Fatal("Error encrypting the RSA message")
	}
	return encryptedMessage
}

func RSADecrypt(cypher message.RSAEncryptedMessage, privateKey *rsa.PrivateKey) []byte {
	hash := sha256.New()
	label := []byte("")
	data, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		privateKey,
		cypher,
		label,
	)
	if err != nil {
		log.Fatal("Error decrypting the RSA cyphertext")
	}
	return data
}

func GenerateRSAKeypair() (priv crypto.PrivateKey, pub crypto.PublicKey) {
	private, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Error generating the RSA keypair")
	}
	return private, private.Public()
}
