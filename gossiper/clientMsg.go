package gossiper

import (
	"fmt"
	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/encConversation"
	"github.com/vquelque/SecuriChat/message"
	"log"
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
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						log.Println("Creating conversation")
						gsp.sendEncryptedTextMessage(cs, encConversation.QueryTextMessage, msg.Destination)
					} else {
						log.Println("convo loaded")
					}
					cs.Buffer <- msg.Text
					return
				} else if msg.AuthQuestion != "" {
					log.Println("Requested authentication with question ", msg.AuthQuestion)
					cs, ok := gsp.createOrLoadConversationState(msg.Destination)
					if !ok {
						log.Panic("convo didn't exist")
					}
					toSend, err := cs.Conversation.StartAuthenticate(msg.AuthQuestion, []byte(msg.AuthAnswer))
					if err != nil {
						log.Panic(err.Error())
					}
					gsp.sendEncryptedMessage(toSend[0], cs, msg.Destination)
				}

			}
			mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
			m := message.NewRumorMessage(gsp.Name, mID, msg.Text)
			gsp.processRumorMessage(m, "")
		}
	}
}
