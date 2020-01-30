package message

import (
	"fmt"

	"github.com/coyim/otr3"
	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/pow"
	"github.com/vquelque/SecuriChat/utils"
)

// SimpleMessage represents a type of SecuriChat message containing only text.
type SimpleMessage struct {
	OriginalName  string
	RelayPeerAddr string
	Contents      string
}

//RumorMessage represents a type of SecuriChat message to be gossiped.
type RumorMessage struct {
	Origin              string
	ID                  uint32
	Text                string
	PoW                 pow.ProofOfWork
	EncryptedMessage    *EncryptedMessage
	RSAEncryptedMessage RSAEncryptedMessage
}

//PrivateMessage between 2 peers
type PrivateMessage struct {
	Origin      string
	ID          uint32
	Text        string
	Destination string
	HopLimit    uint32
	msg         []otr3.ValidMessage
}

type DataRequest struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte //hash of chunk or metafile if file request
}

type DataReply struct {
	Origin      string
	Destination string
	HopLimit    uint32
	HashValue   []byte
	Data        []byte
}

type EncryptedMessage struct {
	Message otr3.ValidMessage
	Step    int
	Dest    string
}

type RSAEncryptedMessage []byte

//NewSimpleMessage creates a new simpleMessage.
func NewSimpleMessage(contents string, originalName string, relayPeerAddr string) *SimpleMessage {
	return &SimpleMessage{
		OriginalName:  originalName,
		RelayPeerAddr: relayPeerAddr,
		Contents:      contents,
	}
}

//NewRumorMessage creates a new rumorMessage.
func NewRumorMessage(origin string, ID uint32, text string) *RumorMessage {
	return &RumorMessage{
		Origin: origin,
		ID:     ID,
		Text:   text,
	}
}

//NewRumorMessage creates a new EncryptedRumorMessage.
func NewRumorMessageWithEncryptedData(origin string, ID uint32, message *EncryptedMessage) *RumorMessage {
	return &RumorMessage{
		Origin:              origin,
		ID:                  ID,
		EncryptedMessage:    message,
		RSAEncryptedMessage: nil,
	}
}

func NewRSARumorMessage(origin string, ID uint32, message RSAEncryptedMessage, encryptedMessage *EncryptedMessage) *RumorMessage {
	return &RumorMessage{
		Origin:              origin,
		ID:                  ID,
		RSAEncryptedMessage: message,
		EncryptedMessage:    encryptedMessage,
	}
}

func (msg *RumorMessage) Encode() []byte {
	b := utils.EncodeUint64(uint64(msg.ID))
	b = append(b, []byte(msg.Origin)...)
	if msg.EncryptedMessage != nil {
		b = append(b, msg.EncryptedMessage.Encode()...)
	}
	b = append(b, msg.Text...)
	b = append(b, msg.RSAEncryptedMessage...)
	return b
}

//NewPrivateMessage creates a new private message for peer dest (dest is peer identifier not address).
// Set hop limit to 0 for default value
func NewPrivateMessage(origin string, text string, destination string, hoplimit uint32) *PrivateMessage {
	if hoplimit == 0 {
		hoplimit = constant.DefaultHopLimit //default hoplimit
	}
	return &PrivateMessage{
		Origin:      origin,
		ID:          0, //no sequencing for private messages
		Text:        text,
		Destination: destination,
		HopLimit:    hoplimit,
	}
}

// NewRouteRumorMessage creates a route rumor message used to updating routing table
// entries of a peer. It is simply a rumor message with empty text field
func NewRouteRumorMessage(origin string, ID uint32) *RumorMessage {
	return &RumorMessage{
		Origin: origin,
		ID:     ID,
		Text:   "",
	}
}

// Set hop limit to 0 for default value (10)
func NewDataReply(origin string, hoplimit uint32, request *DataRequest, data []byte) *DataReply {
	if hoplimit <= 0 {
		hoplimit = constant.DefaultHopLimit //default hoplimit
	}
	d := &DataReply{
		Origin:      origin,
		Destination: request.Origin,
		HopLimit:    hoplimit,
		HashValue:   request.HashValue,
		Data:        data,
	}
	return d
}

// Set hop limit to 0 for default value (10)
func NewDataRequest(origin string, destination string, hoplimit uint32, hashValue utils.SHA256) *DataRequest {
	if hoplimit <= 0 {
		hoplimit = 10 //default hoplimit
	}
	dr := &DataRequest{
		Origin:      origin,
		Destination: destination,
		HopLimit:    hoplimit,
		HashValue:   hashValue[:],
	}
	return dr
}

//Prints a RumorMessage
func (msg *RumorMessage) PrintRumor(relay string) string {
	return fmt.Sprintf("RUMOR origin %s from %s ID %d contents %s", msg.Origin, relay, msg.ID, msg.Text)
}

//Prints simpleMessage.
func (msg *SimpleMessage) String() string {
	return fmt.Sprintf("SIMPLE MESSAGE origin %s from %s contents %s", msg.OriginalName,
		msg.RelayPeerAddr, msg.Contents)
}

// Printes privateMessage
func (msg *PrivateMessage) String() string {
	return fmt.Sprintf("PRIVATE origin %s hop-limit %d contents %s",
		msg.Origin, msg.HopLimit, msg.Text)
}

func (enc *EncryptedMessage) Encode() []byte {
	b := utils.EncodeUint64(uint64(enc.Step))
	b = append(b, []byte(enc.Dest)...)
	b = append(b, enc.Message...)
	return b
}
