package gossiper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/message"
	"github.com/vquelque/SecuriChat/vector"
)

// Procecces incoming rumor message.
func (gsp *Gossiper) processRumorMessage(msg *message.RumorMessage, sender string) {

	next := gsp.VectorClock.NextMessageForPeer(msg.Origin)
	if sender != "" && msg.ID >= next && msg.Origin != gsp.Name {
		gsp.Routing.UpdateRoute(msg, sender) //update routing table
		if msg.Text != "" {
			fmt.Println(gsp.Routing.PrintUpdate(msg.Origin))
		}
	}

	if next == msg.ID {
		// we were waiting for this message
		// increase mID for peer and store message
		gsp.VectorClock.IncrementMIDForPeer(msg.Origin)
		gsp.RumorStorage.Store(msg)
		//if sender is nil then it is a client message
		if sender != "" {
			if msg.Origin != gsp.Name {
				fmt.Println(msg.PrintRumor(sender))
				fmt.Println(gsp.Peers.PrintPeers())
			}
		}
		//TODO CLEANLY STORE IN UIStorage
		//pick random peer and rumormonger
		randPeer := gsp.Peers.PickRandomPeer(sender)
		if randPeer != "" {
			gsp.rumormonger(msg, randPeer)
		}
	}

	// acknowledge the packet if not sent by client
	if sender != "" {
		gsp.sendStatusPacket(sender)
	}
}

// Handle the rumormongering process and launch go routine that listens for ack or timeout.
func (gsp *Gossiper) rumormonger(rumor *message.RumorMessage, peerAddr string) {
	go gsp.listenForAck(rumor, peerAddr)
	gsp.sendRumorMessage(rumor, peerAddr)
	fmt.Printf("MONGERING with %s \n", peerAddr)
}

// Listen and handle ack or timeout.
func (gsp *Gossiper) listenForAck(rumor *message.RumorMessage, peerAddr string) {
	// register this channel inside the map of channels waiting for an ack (observer).
	id := peerAddr + fmt.Sprintf("%s : %d", rumor.Origin, rumor.ID)
	channel := gsp.WaitingForAck.Register(id)
	timer := time.NewTicker(constant.AckTimeout * time.Second)
	defer func() {
		timer.Stop()
		gsp.WaitingForAck.Unregister(id)
	}()

	//keep running while channel open with for loop assignment
	for {
		select {
		case <-timer.C:
			gsp.coinFlip(rumor, peerAddr)
			return
		case ack := <-channel:
			if ack {
				gsp.coinFlip(rumor, peerAddr)
			}
			return
		}
	}
}

// Send rumor to peerAddr.
func (gsp *Gossiper) sendRumorMessage(msg *message.RumorMessage, peerAddr string) {
	gp := GossipPacket{RumorMessage: msg}
	gsp.send(&gp, peerAddr)
}

// CoinFlip tosses a coin. If head, we rumormonger the rumor to a random peer. We exclude the sender
// from the randomly chosen peer.
func (gsp *Gossiper) coinFlip(rumor *message.RumorMessage, sender string) {
	head := rand.Int() % 2
	if head == 0 {
		// exclude the sender of the rumor from the set where we pick our random peer to prevent a loop.
		peer := gsp.Peers.PickRandomPeer(sender)
		if peer != "" {
			fmt.Printf("FLIPPED COIN sending rumor to %s\n", peer)
			gsp.rumormonger(rumor, peer)
		}
	}
}

// Check if we are in sync with peer. Else, send the missing messages to the peer.
func (gsp *Gossiper) synchronizeWithPeer(same bool, toAsk []vector.PeerStatus, toSend []vector.PeerStatus, peerAddr string) {
	if same {
		fmt.Printf("IN SYNC WITH %s \n", peerAddr)
		return
	}
	if len(toSend) > 0 {
		// we have new messages to send to the peer : start mongering
		//get the rumor we need to send from storage
		rumorMsg := gsp.RumorStorage.Get(toSend[0].Identifier, toSend[0].NextID)
		if rumorMsg != nil {
			gsp.rumormonger(rumorMsg, peerAddr)
		}
	} else if len(toAsk) > 0 {
		// send status for triggering peer mongering
		gsp.sendStatusPacket(peerAddr)
	}
}
