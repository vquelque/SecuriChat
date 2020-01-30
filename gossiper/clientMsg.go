package gossiper

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/encConversation"
	"github.com/vquelque/SecuriChat/message"
	"log"
	"strings"
)

// ProcessClientMessage processes client messages
func (gsp *Gossiper) ProcessClientMessage(msg *message.Message) {
	fmt.Println(msg.String())
	if gsp.Simple {
		gp := &GossipPacket{Simple: message.NewSimpleMessage(msg.Text, gsp.Name, gsp.PeersSocket.Address())}
		//broadcast packet
		gsp.broadcastPacket(gp, gsp.PeersSocket.Address())
	} else {
		if msg.Destination != "" && len(msg.Request) == 0 && !msg.Encrypted {
			//private message
			m := message.NewPrivateMessage(gsp.Name, msg.Text, msg.Destination, constant.DefaultHopLimit)
			gsp.processPrivateMessage(m)
		} else {
			//rumor message
			if msg.Encrypted {
				if msg.Text != "" {
					dest := strings.Split(msg.Destination, ",")
					destName := dest[0]
					gsp.registeringPublicKey(dest)
					log.Println("Send encr msg to ", destName)
					cs, ok := gsp.createOrLoadConversationState(destName)
					if !ok {
						log.Println("Creating conversation")
						toSend, _ := cs.Conversation.Send(otr3.ValidMessage(encConversation.QueryTextMessage))
						encMsg := &message.EncryptedMessage{
							Message: toSend[0],
							Step:    cs.Step,
							Dest:    "",
						}
						pub := gsp.RSAPeers.GetPeerPublicKey(destName)
						if pub == nil {
							log.Printf("Can't send to peer %s if there is no public key \n", destName)
							return
						}
						gsp.sendRSAKeyExchangeMessage(encMsg, pub)
					} else {
						log.Println("convo loaded")
					}
					cs.Buffer <- msg.Text
					return
				} else if msg.AuthQuestion != "" && msg.AuthAnswer != "" {
					fmt.Printf("Requested authentication with question %s and answer %s", msg.AuthQuestion, msg.AuthAnswer)
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						msg2 := &message.Message{
							Text:        "Initiating conversation",
							Origin:      msg.Origin,
							Encrypted:   true,
							Destination: msg.Destination,
						}
						gsp.ProcessClientMessage(msg2)

						log.Println("conversation didn't exist ", msg.Destination, " was requested, now initializing new conversation")
					}
					go func() { cs.QuestionChan <- [2]string{msg.AuthQuestion, msg.AuthAnswer} }()
					//gsp.sendEncryptedMessage(toSend[0], cs, msg.Destination)
				} else if msg.AuthAnswer != "" {
					fmt.Printf("Provided answer for SMP key exchange %s", msg.AuthAnswer)
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						log.Printf("Requested for %s by %s \n", msg.Destination, gsp.Name)
						log.Panic("convo didn't exist")
					}
					go func() { cs.AnswerChan <- msg.AuthAnswer }()
				} else if msg.Destination != "" {
					dest := strings.Split(msg.Destination, ",")
					log.Printf("Registering peer %s", dest[0])
					gsp.registeringPublicKey(dest)
				}

			} else {
				mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
				m := message.NewRumorMessage(gsp.Name, mID, msg.Text)
				gsp.processRumorMessage(m, "")
			}
		}
	}
}

func (gsp *Gossiper) registeringPublicKey(dest []string) {
	destName := dest[0]
	if len(dest) == 2 {
		log.Println("Registering public key")
		destPubKey := dest[1]
		gsp.parseAndStoreRSAPublicKey(destPubKey, destName)
	}
}

func (gsp *Gossiper) parseAndStoreRSAPublicKey(destPubKey string, destName string) {
	destPubKey = "-----BEGIN PUBLIC KEY-----\n" + destPubKey + "\n-----END PUBLIC KEY-----"
	block, _ := pem.Decode([]byte(destPubKey))
	if block == nil {
		panic("failed to parse PEM block containing the public key")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic("failed to parse DER encoded public key: " + err.Error())
	}
	switch pub := pub.(type) {
	case *rsa.PublicKey:
		fmt.Println("pub is of type RSA:", pub)
	default:
		panic("unknown type of public key")
	}
	gsp.RSAPeers.Add(destName, pub.(*rsa.PublicKey))
}
