package main

import (
	"distributer/p2p"
	"log"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOps{
		ListenAddress: listenAddr,
		HandshakeFunc: p2p.NoHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// OnPeer: ,
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileserverOpts := FileServeropt{
		StorageRoot:       listenAddr + "_NETWORK",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNetwork:  nodes,
	}
	s := newFileServer(fileserverOpts)
	tcpTransport.OnPeer = s.OnPeer

	return s
}
func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")

	go func() {
		log.Fatal(s1.start())
	}()
	s2.start()
}
