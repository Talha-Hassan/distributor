package main

import (
	"distributer/p2p"
	"fmt"
	"log"
	"time"
)

func main() {
	tcpopts := p2p.TCPTransportOps{
		ListenAddress: ":3000",
		HandshakeFunc: p2p.NoHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tr := p2p.NewTCPTransport(tcpopts)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			msg := <-tr.Consume()
			fmt.Printf("Message: %+v\n", msg)
		}
	}()
	if err := tr.ListerAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
