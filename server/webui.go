package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/vquelque/SecuriChat/gossiper"
	"github.com/vquelque/SecuriChat/message"
	"github.com/vquelque/SecuriChat/utils"
)

func peersListHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			peerList := gsp.Peers.GetAllPeers()
			peerListJSON, err := json.Marshal(peerList)
			if err != nil {
				log.Print("Error sending peers list as JSON")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(peerListJSON)
		case "POST":
			http.Redirect(w, r, r.Header.Get("/"), 302)
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			peerAddr := r.FormValue("peerAddr")
			if peerAddr == gsp.PeersSocket.Address() {
				return
			}
			peerAddrChecked := utils.ToUDPAddr(peerAddr)
			if peerAddrChecked == nil {
				return
			}
			if !gsp.Peers.Contains(peerAddr) {
				gsp.Peers.Add(peerAddr)
			} else {
				gsp.Peers.Delete(peerAddr)
			}
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	})
}
func msgHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			var rumorMessageList []message.RumorMessage
			rumorMessageList = gsp.RumorStorage.GetAll()
			mmsgListJSON, err := json.Marshal(rumorMessageList)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(mmsgListJSON)
		case "POST":
			http.Redirect(w, r, r.Header.Get("/"), 302)
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid Data", http.StatusBadRequest)
				return
			}
			messageText := r.FormValue("message")
			cliMsg := &message.Message{Text: messageText}
			gsp.ProcessClientMessage(cliMsg)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	})
}

func idHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			peerID := gsp.Name
			peerIDJSON, err := json.Marshal(peerID)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(peerIDJSON)
		}
	})
}

func contactsHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			contacts := gsp.Routing.GetAllRoutes()
			contactsJSON, err := json.Marshal(contacts)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(contactsJSON)
		}
	})
}

func privateMsgHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			peer := r.FormValue("peer")
			m := gsp.PrivateStorage.GetAllForPeer(peer)
			mJSON, err := json.Marshal(m)
			if err != nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(mJSON)
		case "POST":
			if err := r.ParseForm(); err != nil {
				http.Error(w, "Invalid Data", http.StatusBadRequest)
				return
			}
			peer := r.FormValue("peer")
			messageText := r.FormValue("message")
			print(messageText)
			cliMsg := &message.Message{Text: messageText, Destination: peer}
			gsp.ProcessClientMessage(cliMsg)
			http.Redirect(w, r, r.Header.Get("/privateMsg?peer="+peer), 302)
		default:
			fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
		}
	})
}

// StartUIServer starts the UI server
func StartUIServer(UIPort int, gsp *gossiper.Gossiper) *http.Server {

	UIPortStr := ":" + strconv.Itoa(UIPort)
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("server/")))
	mux.HandleFunc("/id", idHandler(gsp))
	mux.HandleFunc("/peers", peersListHandler(gsp))
	mux.HandleFunc("/message", msgHandler(gsp))
	mux.HandleFunc("/contacts", contactsHandler(gsp))
	mux.HandleFunc("/privateMsg", privateMsgHandler(gsp))
	server := &http.Server{Addr: UIPortStr, Handler: mux}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
			return
		}
		fmt.Printf("UI server started at address 127.0.0.1:%s", UIPortStr)
	}()
	return server
}
