package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func bootstrapHandler(w http.ResponseWriter, r *http.Request) {
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

	registerAddres(&mess)
	err = updateNodes()
	if err != nil {
		log.Warnf("Error while updating nodes. Reason: %s", err)
	} else {
		log.Info("Nodes updarted successfully")
	}
}

func resolveMessage(mess []byte) (NewConnectionMessage, error) {
	var resolvedMess NewConnectionMessage
	err := json.Unmarshal(mess, &resolvedMess)

	return resolvedMess, err
}

func registerAddres(newConn *NewConnectionMessage) {
	nodeAddress := fmt.Sprintf("%s:%d", (*newConn).Ip, (*newConn).Port)
	if !checkIfAddressAlreadyKnown(nodeAddress) {
		knownNodeAddresses = append(knownNodeAddresses, nodeAddress)
		log.Infof("Registered node IP: %s", nodeAddress)
	}
}

func updateNodes() error {
	updateMess := AvailableNodesAddresses{Type: "addresses", Addresses: knownNodeAddresses}
	serializedMess, err := json.Marshal((updateMess))
	if err != nil {
		return err
	}

	for _, addr := range knownNodeAddresses {
		log.Infof("Sending update message to %s about available nodes: %s", addr, knownNodeAddresses)
		sendWsMessage(addr, serializedMess)
	}

	return nil
}

func checkIfAddressAlreadyKnown(addr string) bool {
	for _, v := range knownNodeAddresses {
		if addr == v {
			return true
		}
	}
	return false
}
