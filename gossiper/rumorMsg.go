package gossiper

import (
	"bufio"
	"fmt"
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/encConversation"
	"github.com/vquelque/SecuriChat/pow"
	"log"
	"math/rand"
	"os"
	"regexp"
	"time"

	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/message"
	"github.com/vquelque/SecuriChat/vector"
)

// Processes incoming rumor message.
func (gsp *Gossiper) processRumorMessage(msg *message.RumorMessage, sender string) {

	log.Println("rumor received")

	next := gsp.VectorClock.NextMessageForPeer(msg.Origin)
	if sender != "" && msg.ID >= next && msg.Origin != gsp.Name {
		gsp.Routing.UpdateRoute(msg, sender) //update routing table
		if msg.Text != "" {
			fmt.Println(gsp.Routing.PrintUpdate(msg.Origin))
		}
	}

	if next == msg.ID {

		if gsp.Name != msg.Origin {
			if !msg.PoW.Validator(msg.Encode()) {
				log.Println("WARNING : Invalid PoW")
				return
			}
		} else {
			msg.PoW = pow.NewProofOfWork(msg.Encode())
		}

		// we were waiting for this message
		// increase mID for peer and store message
		gsp.VectorClock.IncrementMIDForPeer(msg.Origin)
		gsp.RumorStorage.Store(msg)
		//if sender is nil then it is a client message
		if sender != "" {
			if msg.Origin != gsp.Name {
				gsp.handleEncryptedMessage(msg)
				fmt.Println(msg.PrintRumor(sender))
				fmt.Println(gsp.Peers.PrintPeers())

			}
		}
		//TODO CLEANLY STORE IN UIStorage
		//pick random peer and rumormonger
		log.Println("sending message")
		gsp.sendRumor(sender, msg)
	}

	// acknowledge the packet if not sent by client
	if sender != "" {
		gsp.sendStatusPacket(sender)
	}
}

func (gsp *Gossiper) sendRumor(sender string, msg *message.RumorMessage) {
	randPeer := gsp.Peers.PickRandomPeer(sender)
	if randPeer != "" {
		gsp.rumormonger(msg, randPeer)
	}
}

func (gsp *Gossiper) handleEncryptedMessage(msg *message.RumorMessage) {
	encryptedMessage := msg.EncryptedMessage
	if encryptedMessage != nil {
		log.Println("handling enc msg")
		if encryptedMessage.Dest != gsp.Name {
			log.Println("not me!")
			return
		}
		//Get conversation
		cs, _ := gsp.createOrLoadConversationState(msg.Origin)
		if encryptedMessage.Step == 0 && cs.Step > 0 {
			log.Println("Re-Doing key exchange")
			gsp.convStateMap.Update(msg.Origin, nil)
			cs, _ = gsp.createOrLoadConversationState(msg.Origin)
		}

		plaintxt, toSend, err := cs.Conversation.Receive(encryptedMessage.Message)
		if err != nil {
			log.Fatal(err.Error())
		}


		switch encryptedMessage.Step {
		case encConversation.QueryMsg, encConversation.DHCommit,
			encConversation.DHKey, encConversation.RevealSig:
			log.Println("state is : ", cs.Step)
			log.Printf("Doing key exchange, step %d \n", encryptedMessage.Step+1)
			cs.Step = encryptedMessage.Step + 1
			log.Println("state final is : ", cs.Step)
			gsp.sendEncryptedMessage(toSend[0], cs, msg.Origin)

			if cs.Step == encConversation.Sig {
				gsp.endOfKeyExchange(cs, msg)
			}
		case encConversation.Sig:
			gsp.endOfKeyExchange(cs, msg)
		case encConversation.AkeFinished, encConversation.AuthenticationOK:
			// A message was received
			msgType := " AUTHENTICATED "
			if encConversation.AkeFinished == cs.Step {
				msgType = " UNAUTHENTICATED "
			}
			fmt.Printf("RECEIVED %s ENCR MESSAGE : \n %s \n ", msgType, plaintxt)

		case encConversation.SMP1:
			log.Println("state is : ", cs.Step)
			log.Printf("Doing SMP Protocol, step %d \n", encryptedMessage.Step+1)
			cs.Step = encryptedMessage.Step + 1
			reader := bufio.NewReader(os.Stdin)
			question,_ := cs.Conversation.SMPQuestion()
			fmt.Printf("Enter the secret for the question %s : \n", question)
			secret, _ := reader.ReadString('\n')
			fmt.Printf("Secret is %s",secret)
			secret = removeEndOfLine(secret)
			fmt.Print(secret)
			toSend, err := cs.Conversation.ProvideAuthenticationSecret([]byte(secret))
			if err != nil {
				log.Panic(err.Error())
			}
			gsp.sendEncryptedMessage(toSend[0], cs, msg.Origin)

		case encConversation.SMP2, encConversation.SMP3:

			log.Println("state is : ", cs.Step)
			log.Printf("Doing SMP Protocol, step %d \n", encryptedMessage.Step+1)
			cs.Step = encryptedMessage.Step + 1
			log.Println("state final is : ", cs.Step)
			gsp.sendEncryptedMessage(toSend[0], cs, msg.Origin)

		case encConversation.SMP4, encConversation.SMP5:
			log.Println("Should be ok")
			cs.Step = encConversation.AuthenticationOK

		}
	} else {
		log.Println("rumor message")
	}
}

func removeEndOfLine(secret string) string {
	re := regexp.MustCompile(`\r?\n`)
	secret = re.ReplaceAllString(secret, "")
	return secret
}

func (gsp *Gossiper) endOfKeyExchange(cs *encConversation.ConversationState, msg *message.RumorMessage) {
	log.Println("Key exchange is finished")
	go gsp.sendBufferedEncrRumors(cs, msg)
	cs.Step = encConversation.AkeFinished
}

func (gsp *Gossiper) sendBufferedEncrRumors(cs *encConversation.ConversationState, msg *message.RumorMessage) {
	for textMessage := range cs.Buffer {
		gsp.sendEncryptedTextMessage(cs, textMessage, msg.Origin)
	}
}

func (gsp *Gossiper) sendEncryptedMessage(toSend otr3.ValidMessage, cs *encConversation.ConversationState, dest string) {
	log.Println("Sending encryptedMessage")
	mID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
	encMsg := &message.EncryptedMessage{
		Message: toSend,
		Step:    cs.Step,
		Dest:    dest,
	}
	rumor := message.NewRumorMessageWithEncryptedData(gsp.Name, mID, encMsg)
	gsp.processRumorMessage(rumor, "")
}

func (gsp *Gossiper) sendEncryptedTextMessage(cs *encConversation.ConversationState, text string, dest string) {
	toSend, _ := cs.Conversation.Send(otr3.ValidMessage(text))
	gsp.sendEncryptedMessage(toSend[0], cs, dest)
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
