package main

import (
	"flag"

	"dope.bootstrap/server"
)

func main() {
	addr := flag.String("address", "127.0.0.1", "Ip address to the bootstrap server")
	port := flag.Int("port", 7312, "Port to run boostrap server on")
	flag.Parse()

	server.Run(addr, port)
}
