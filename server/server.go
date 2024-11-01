package server

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var knownNodeAddresses []string

const (
	NODES_ENDPOINT = "/nodes"
)

func Run(addr *string, port *int) {
	serverAddress := fmt.Sprintf("%s:%d", *addr, *port)

	http.HandleFunc(NODES_ENDPOINT, nodesHandler)

	log.Printf("Running bootstrap server on %s:%d", *addr, *port)
	http.ListenAndServe(serverAddress, nil)
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("Error upgrading to WebSocket:", err)
	}

	return conn, err
}

func sendWsMessage(targetAddress string, message []byte) {
	u := url.URL{Scheme: "ws", Host: targetAddress, Path: "/node"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Warnf("Unable to establish connection to %s", targetAddress)
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Warnf("Unable to write message to %s", targetAddress)
	}

	log.Infof("Message sent to %s", targetAddress)
}
