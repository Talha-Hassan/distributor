package main

import (
	"distributer/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(":3000")

	if err := tr.ListerAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
