package main

import (
	"distributer/p2p"
	"log"
)

func main() {
	tcpopts := p2p.TCPTransportOps{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NoHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}
	tr := p2p.NewTCPTransport(tcpopts)
	if err := tr.ListerAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
