package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/dedis/protobuf"
	"github.com/vquelque/SecuriChat/message"
)

func main() {

	uiPort := flag.Int("UIPort", 8080, "Port for the UI client (default 8080)")
	text := flag.String("msg", "", "message to be sent; if the -dest flag is present, this is a private message, otherwise it’s a rumor message")
	destinationName := flag.String("destName", "", "destinationName for the private message. can be omitted")
	destinationPubKey := flag.String("destPubKey", "", "destination public key")
	encrypted := flag.Bool("encrypted", false, "encrypted message")
	authAnswer := flag.String("authAnswer", "", "use this flag to the answer to the question")
	authQuestion := flag.String("authQuestion", "", "use this flag to indicate that you want to authenticate by providing this question ")

	flag.Parse()

	addr := fmt.Sprintf("127.0.0.1:%d", *uiPort) //localhost gossiper address
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	udpAddrCli, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	udpConn, err := net.ListenUDP("udp4", udpAddrCli)
	if err != nil {
		log.Fatalln(err)
	}

	msg := &message.Message{}
	msg.Destination = *destinationName
	if *destinationPubKey != "" {
		msg.Destination = msg.Destination + "," + *destinationPubKey

	}

	msg.Text = *text
	msg.Encrypted = *encrypted
	msg.AuthAnswer = *authAnswer
	msg.AuthQuestion = *authQuestion
	pkt, err := protobuf.Encode(msg)

	if err != nil {
		log.Fatalln(err)
	}

	_, err = udpConn.WriteToUDP(pkt, udpAddr)
	if err == nil {
		fmt.Printf("CLIENT MESSAGE sent to %s \n", udpAddr.String())
	}
}
