package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func nodesHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := getWebsocketConnection(w, r)
	if err != nil {
		log.Warnf("Failed to establish websocket connection. Reason: %s", err)
		return
	}

	_, receivedMess, err := connection.ReadMessage()
	if err != nil {
		log.Warnf("Failed to read message. Reason: %s", err)
	}

	mess, err := resolveMessage(receivedMess)
	if err != nil {
		log.Warnf("Failed to parse received message. Reason: %s", err)
	}

	newAddr := registerAddres(&mess)
	err = updateNodes(&newAddr)
	if err != nil {
		log.Warnf("Error while updating nodes. Reason: %s", err)
	} else {
		log.Info("Successfully added new node")
	}
}

func resolveMessage(mess []byte) (NewConnectionMessage, error) {
	var resolvedMess NewConnectionMessage
	err := json.Unmarshal(mess, &resolvedMess)

	return resolvedMess, err
}

func registerAddres(newConn *NewConnectionMessage) string {
	nodeAddress := fmt.Sprintf("%s:%d", (*newConn).Ip, (*newConn).Port)
	knownNodeAddresses = append(knownNodeAddresses, nodeAddress)
	log.Infof("Registered node IP: %s", nodeAddress)

	return nodeAddress
}

func updateNodes(newAddr *string) error {
	updateMess := UpdateNodesMessage{Address: *newAddr}
	serializedMess, err := json.Marshal((updateMess))
	if err != nil {
		return err
	}

	for _, addr := range knownNodeAddresses {
		log.Infof("Sending update message to %s about new node at %s", addr, *newAddr)
		sendWsMessage(addr, serializedMess)
	}

	return nil
}
