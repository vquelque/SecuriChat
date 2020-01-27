package gossiper

import (
	"fmt"
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
				if msg.AuthAnswer == "" {
					log.Println("Send encr msg")
					dest := strings.Split(msg.Destination, ",")
					destName := dest[0]
					//destPubKey := dest[1]
					cs, ok := gsp.createOrLoadConversationState(destName)
					if !ok {
						log.Println("Creating conversation")
						gsp.sendEncryptedTextMessage(cs, encConversation.QueryTextMessage, msg.Destination)
					} else {
						log.Println("convo loaded")
					}
					cs.Buffer <- msg.Text
					return
				} else if msg.AuthQuestion != "" {
					log.Printf("Requested authentication with question %s and answer %s", msg.AuthQuestion, msg.AuthAnswer)
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						log.Panic("convo didn't exist")
					}
					toSend, err := cs.Conversation.StartAuthenticate(msg.AuthQuestion, []byte(msg.AuthAnswer))
					if err != nil {
						log.Panic(err.Error())
					}
					cs.Step = encConversation.SMP1
					gsp.sendEncryptedMessage(toSend[0], cs, msg.Destination)
				} else {
					log.Printf("Requested authentication with answer %s", msg.AuthAnswer)
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						log.Panic("convo didn't exist")
					}
					go func() { cs.AnswerChan <- msg.AuthAnswer }()
				}

			}
			mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
			m := message.NewRumorMessage(gsp.Name, mID, msg.Text)
			gsp.processRumorMessage(m, "")
		}
	}
}
