package main

import (
	"flag"

	"dope.bootstrap/server"
)

func main() {
	addr := flag.String("address", "0.0.0.0", "Ip address to the bootstrap server")
	port := flag.Int("port", 7312, "Port to run boostrap server on")
	flag.Parse()

	server.Run(addr, port)
}
