package gossiper

import (
	"bytes"
	"fmt"
	"github.com/vquelque/SecuriChat/crypto"
	"github.com/vquelque/SecuriChat/message"
	"log"
)

func (gsp *Gossiper) handleRSAEncryptedMessage(rumor *message.RumorMessage) {
	//log.Println("Handling RSA msg from", rumor.Origin)
	if !gsp.RSAPeers.Contains(rumor.Origin) {
		log.Println("Origin isn't known by the gossiper")
		return
	}
	encHash, err := crypto.RSADecrypt(rumor.RSAEncryptedMessage, gsp.RSAPrivateKey)
	if err != nil {
		log.Printf("Received encryted rumor but is not for us. \n")
		return
	}
	hash := GetHashOfEncryptedMessage(rumor.EncryptedMessage)
	i := hash[:]
	isEq := bytes.Compare(encHash, i)
	if isEq != 0 {
		log.Println(isEq)
		log.Println(len(i))
		log.Println(len(encHash))
		log.Println("Hash isn't equal, something is wrong")
		log.Println(hash)
		log.Println(encHash)
		return
	}

	fmt.Printf("successfully decrypted RSA encrypted message. Starting key exchange. \n")
	gsp.handleEncryptedMessage(rumor)
}
