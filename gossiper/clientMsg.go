package gossiper

import (
	"fmt"
	"github.com/vquelque/SecuriChat/encConversation"
	"log"

	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/message"
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
			if msg.Encrypted{
				log.Println("Send encr msg")
				msge := []byte("?OTRv3?")
				cs,ok := gsp.convStateMap.Load(msg.Destination)
				if !ok{
					cs = gsp.CreateConversationState()
					gsp.convStateMap.Update(msg.Destination, cs)
				}
				toSend,err := cs.Conversation.Send(msge)
				if err!=nil{
					log.Fatal(err.Error())
				}
				mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
				encMsg := &encConversation.EncryptedMessage{
					Message: toSend[0],
					Step:   cs.Step,
					Dest:    msg.Destination,
				}
				rumor := message.NewRumorMessageWithEncryptedData(gsp.Name,mID,encMsg)
				gsp.processRumorMessage(rumor, "")

				return
			}
			mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
			m := message.NewRumorMessage(gsp.Name, mID, msg.Text)
			gsp.processRumorMessage(m, "")
		}
	}
}
