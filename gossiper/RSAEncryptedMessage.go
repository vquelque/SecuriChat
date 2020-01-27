package gossiper

import (
	"log"

	"github.com/dedis/protobuf"
	"github.com/vquelque/SecuriChat/crypto"
	"github.com/vquelque/SecuriChat/message"
)

func (gsp *Gossiper) handleRSAEncryptedMessage(rumor *message.RumorMessage) {
	if !gsp.RSAPeers.Contains(rumor.Origin) {
		return
	}
	plain, err := crypto.RSADecrypt(rumor.RSAEncryptedMessage, gsp.RSAPrivateKey)
	if err != nil {
		log.Printf("Received encryted rumor but is not for us. \n")
		return
	}
	encMessage := &message.EncryptedMessage{}
	protobuf.Decode(plain, encMessage)
	log.Printf("successfully decrypted RSA encrypted message. Starting key exchange. \n", plain)
	encRumor := &message.RumorMessage{
		Origin:           rumor.Origin,
		ID:               rumor.ID,
		Text:             rumor.Text,
		EncryptedMessage: encMessage,
	}
	go gsp.handleEncryptedMessage(encRumor)
}
