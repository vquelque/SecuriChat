package gossiper

import "github.com/vquelque/SecuriChat/message"

import "github.com/vquelque/SecuriChat/crypto"

import "log"

import "fmt"

func (gsp *Gossiper) handleRSAEncryptedMessage(message *message.RumorMessage) {
	if !gsp.SubscribedPeers.Contains(message.Origin) {
		return
	}
	plain, err := crypto.RSADecrypt(message.RSAEncryptedMessage, gsp.RSAPrivateKey)
	if err != nil {
		log.Printf("Received encryted rumor but is not for us. \n")
		return
	}
	fmt.Printf("successfully decrypted encrypted rumor. \n Message : %s \n", plain)
}
