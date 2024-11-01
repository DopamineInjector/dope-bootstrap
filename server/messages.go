package server

type NewConnectionMessage struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type UpdateNodesMessage struct {
	Address string `json:"address"`
}
