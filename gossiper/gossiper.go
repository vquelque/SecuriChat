package gossiper

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/crypto"
	"github.com/vquelque/SecuriChat/encConversation"
	"log"
	"sync"
	"time"

	"github.com/dedis/protobuf"
	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/message"
	"github.com/vquelque/SecuriChat/observer"
	"github.com/vquelque/SecuriChat/peers"
	"github.com/vquelque/SecuriChat/routing"
	"github.com/vquelque/SecuriChat/socket"
	"github.com/vquelque/SecuriChat/storage"
	"github.com/vquelque/SecuriChat/vector"
)

// Gossiper main structure
type Gossiper struct {
	Name                  string
	Peers                 *peers.Peers
	Simple                bool
	PeersSocket           socket.Socket
	UISocket              socket.Socket
	VectorClock           *vector.Vector        //current status of this peer.
	RumorStorage          *storage.RumorStorage //store all previously received rumors.
	PrivateStorage        *storage.PrivateStorage
	Active                sync.WaitGroup     //Active go routines.
	WaitingForAck         *observer.Observer //registered go routines channels waiting for an ACK.
	AntiEntropyTimer      int
	ResetAntiEntropyTimer chan bool
	Routing               *routing.Routing
	Rtimer                int
	convStateMap          *encConversation.ConvStateMap
	OTRprivateKey         *otr3.DSAPrivateKey
	RSAPrivateKey         *rsa.PrivateKey
	RSAPublickKey         *rsa.PublicKey
}

// GossipPacket is the only type of packet sent to other peers.
type GossipPacket struct {
	Simple       *message.SimpleMessage
	RumorMessage *message.RumorMessage
	StatusPacket *vector.StatusPacket
	Private      *message.PrivateMessage
	DataRequest  *message.DataRequest
	DataReply    *message.DataReply
}

// Encapsulate received messages from peers/client to put in the queue
type receivedPackets struct {
	data   []byte
	sender string
}

// NewGossiper creates and returns a new gossiper running at given address, port with given name.
func NewGossiper(address string, uiPort int, peersList string, simple bool, antiEntropyTimer int, rtimer int) *Gossiper {
	peersSocket := socket.NewUDPSocket(address)
	uiSocket := socket.NewUDPSocket(fmt.Sprintf("127.0.0.1:%d", uiPort))

	peersSet := peers.NewPeersSet(peersList)
	vectorClock := vector.NewVector()
	rumorStorage := storage.NewRumorStorage()
	privateStorage := storage.NewPrivateStorage()
	waitingForAck := observer.Init()
	resetAntiEntropyChan := make(chan (bool))
	routing := routing.NewRoutingTable()
	name := encConversation.GetRandomName()

	OTRpriv := &otr3.DSAPrivateKey{}
	_ = OTRpriv.Generate(rand.Reader)
	RSAPriv, RSAPub := crypto.GenerateRSAKeypair()

	return &Gossiper{
		Name:                  name,
		Peers:                 peersSet,
		Simple:                simple,
		PeersSocket:           peersSocket,
		UISocket:              uiSocket,
		VectorClock:           vectorClock,
		RumorStorage:          rumorStorage,
		PrivateStorage:        privateStorage,
		WaitingForAck:         waitingForAck,
		Active:                sync.WaitGroup{},
		AntiEntropyTimer:      antiEntropyTimer,
		ResetAntiEntropyTimer: resetAntiEntropyChan,
		Routing:               routing,
		Rtimer:                rtimer,
		convStateMap:          encConversation.InitConvStateMap(),
		OTRprivateKey:         OTRpriv,
		RSAPrivateKey:         RSAPriv,
		RSAPublickKey:         RSAPub,
	}
}

////////////////////////////
// Packets, GossipPacket //
////////////////////////////

// serialize with protobuf and send the gossipPacket to the provided UDP addr using the provided gossiper
func (gsp *Gossiper) send(gossipPacket *GossipPacket, addr string) {
	pkt, err := protobuf.Encode(gossipPacket)
	if err != nil {
		log.Print(err)
	}
	gsp.PeersSocket.Send(pkt, addr)
}

func (gsp *Gossiper) broadcastPacket(pkt *GossipPacket, sender string) {
	for _, peer := range gsp.Peers.GetAllPeers() {
		gsp.send(pkt, peer)
	}
}

////////////////////////////
// SimpleMessage //
////////////////////////////
func (gsp *Gossiper) newForwardedMessage(msg *message.SimpleMessage) *message.SimpleMessage {
	msg = message.NewSimpleMessage(msg.Contents, msg.OriginalName, gsp.PeersSocket.Address())
	return msg
}

func (gsp *Gossiper) processSimpleMessage(msg *message.SimpleMessage) {
	fmt.Println(msg.String())
	fmt.Println(gsp.Peers.PrintPeers())
	fwdMsg := gsp.newForwardedMessage(msg)
	packet := &GossipPacket{Simple: fwdMsg}
	gsp.broadcastPacket(packet, msg.RelayPeerAddr)
}

////////////////////////////
// Routing //
////////////////////////////
func (gsp *Gossiper) startRoutingMessageHandler() {
	rTimerDuration := time.Duration(gsp.Rtimer) * time.Second
	timer := time.NewTicker(rTimerDuration)
	//send initial routing message to all neighbors
	for _, peer := range gsp.Peers.GetAllPeers() {
		gsp.sendRouteRumor(peer)
	}
	go func() {
		for {
			select {
			case <-timer.C:
				// timer elapsed : send route rumor packet to randomly chosen peer
				randPeer := gsp.Peers.PickRandomPeer("")
				if randPeer != "" {
					gsp.sendRouteRumor(randPeer)
				}
			}
		}
	}()
}

func (gsp *Gossiper) sendRouteRumor(peer string) {
	rID := gsp.VectorClock.NextMessageForPeer(gsp.Name)
	r := message.NewRouteRumorMessage(gsp.Name, rID)
	gsp.processRumorMessage(r, "")
}

////////////////////////////
// Network //
////////////////////////////
// Handles the incoming packets.
func handleIncomingPackets(socket socket.Socket) <-chan *receivedPackets {
	out := make(chan *receivedPackets, constant.ChannelSize)
	go func() {
		for {
			data, sender := socket.Receive()
			receivedPackets := &receivedPackets{data: data, sender: sender}
			out <- receivedPackets
		}
	}()
	return out
}

// Processes the incoming messages.
func (gsp *Gossiper) processMessages(peerMsgs <-chan *receivedPackets, clientMsgs <-chan *receivedPackets) {
	for {
		select {
		case peerMsg := <-peerMsgs:
			var gp *GossipPacket = &GossipPacket{}
			err := protobuf.Decode(peerMsg.data, gp)
			if peerMsg.sender != "" {
				gsp.Peers.Add(peerMsg.sender)
			}
			if err != nil {
				log.Print(err)
			}
			switch {
			case gp.Simple != nil:
				// received a simple message
				go gsp.processSimpleMessage(gp.Simple)
			case gp.RumorMessage != nil:
				// received a rumorMessage
				gsp.processRumorMessage(gp.RumorMessage, peerMsg.sender)
			case gp.StatusPacket != nil:
				go gsp.processStatusPacket(gp.StatusPacket, peerMsg.sender)
			case gp.Private != nil:
				go gsp.processPrivateMessage(gp.Private)

			}
		case cliMsg := <-clientMsgs:
			msg := &message.Message{}
			err := protobuf.Decode(cliMsg.data, msg)
			if err != nil {
				//	log.Print(err)
			}
			go gsp.ProcessClientMessage(msg)
		}
	}
}

// KillGossiper hills the gossiper
func (gsp *Gossiper) KillGossiper() {
	gsp.PeersSocket.Close()
	gsp.UISocket.Close()
	gsp.Active.Done()
	gsp = nil
}

// Start starts the gossiper
func (gsp *Gossiper) Start() {
	gsp.Active.Add(1)
	peerChan := handleIncomingPackets(gsp.PeersSocket)
	clientChan := handleIncomingPackets(gsp.UISocket)
	go gsp.processMessages(peerChan, clientChan)
	if !gsp.Simple && gsp.AntiEntropyTimer > 0 {
		gsp.startAntiEntropyHandler()
	}
	if gsp.Rtimer > 0 {
		gsp.startRoutingMessageHandler()
	}
	fmt.Printf("Gossiper running at address %s, unique peer id : %s \n", gsp.PeersSocket.Address(), gsp.Name)
}
