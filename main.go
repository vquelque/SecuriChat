package main

import (
	"flag"

	"github.com/vquelque/SecuriChat/constant"
	"github.com/vquelque/SecuriChat/gossiper"
	"github.com/vquelque/SecuriChat/server"
)

func main() {
	uiPort := flag.Int("UIPort", 8080, "Port for the UI client (default 8080)")
	gossipAddr := flag.String("gossipAddr", "", "ip:port for the gossiper")
	name := flag.String("name", "", "Name of the gossiper")
	peersList := flag.String("peers", "", "Comma separated list of peers of the form ip:port")
	simple := flag.Bool("simple", false, "Run gossiper in simple broadcast mode")
	antiEntropy := flag.Int("antiEntropy", 10, "Anti entropy timer value in seconds (default to 10sec)")
	startUIServer := flag.Bool("uisrv", false, "set to true to start the UI server on the UI port")
	rtimerFlag := flag.Int("rtimer", 0, "time between sending two route rumor messages")
	flag.Parse()

	antiEntropyTimer := *antiEntropy
	if antiEntropyTimer < 0 {
		antiEntropyTimer = constant.DefaultAntiEntropy
	}
	rtimer := *rtimerFlag
	if rtimer < 0 {
		rtimer = constant.DefaultRTimer
	}

	gossiper := gossiper.NewGossiper(*gossipAddr, *name, *uiPort, *peersList, *simple, antiEntropyTimer, rtimer)

	//starts UI server if flag is set
	if *startUIServer {
		uiServer := server.StartUIServer(*uiPort, gossiper)
		defer uiServer.Shutdown(nil)
	}

	gossiper.Start()
	defer gossiper.KillGossiper()
	gossiper.Active.Wait()

}
