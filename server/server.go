package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vquelque/SecuriChat/gossiper"
	"github.com/vquelque/SecuriChat/message"
)

type frontendData struct {
	PeerId    string
	PubRSAKey string
}

// Upgrader to upgrade the http connection to a websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	// Check for the origin of incoming connection.
	// For now let everyone connect.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// define a reader which will listen for new messages on the webSocket
func reader(conn *websocket.Conn, gsp *gossiper.Gossiper) {
	for {
		// read in a message
		cliMsg := &message.Message{}
		err := conn.ReadJSON(cliMsg)
		if err != nil {
			log.Println(err)
			return
		}
		sendMessage(conn, gsp, cliMsg)
	}
}
func sendMessage(conn *websocket.Conn, gsp *gossiper.Gossiper, cliMsg *message.Message) {
	cliMsg.Origin = gsp.Name
	go gsp.ProcessClientMessage(cliMsg)
	if err := conn.WriteJSON(cliMsg); err != nil {
		log.Println(err)
		return
	}
}

// define our WebSocket endpoint
func serveWs(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Host)

		// upgrade this connection to a WebSocket
		// connection
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		// listen indefinitely for new messages coming
		// through on our WebSocket connection
		reader(ws, gsp)
	})
}

func initHandler(gsp *gossiper.Gossiper) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		switch r.Method {
		case "GET":
			frontendData := &frontendData{
				PeerId:    gsp.Name,
				PubRSAKey: "RSA PUB KEY",
			}
			frontendDataJSON, err := json.Marshal(frontendData)
			if err != nil {
				log.Print("error encoding json frontend data")
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(frontendDataJSON)
		}
	})
}

func setupRoutes(gsp *gossiper.Gossiper) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	// mape our `/ws` endpoint to the `serveWs` function
	http.HandleFunc("/ws", serveWs(gsp))
	http.HandleFunc("/init", initHandler(gsp))
}

func StartReactServer(gsp *gossiper.Gossiper) {
	fmt.Println("Running SecuriChat websocket server")
	setupRoutes(gsp)
	http.ListenAndServe(":8080", nil)
}
