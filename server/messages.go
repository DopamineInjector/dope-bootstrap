package server

type NewConnectionMessage struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type AvailableNodesAddresses struct {
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}
