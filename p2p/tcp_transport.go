package p2p

import (
	"fmt"
	"net"
	"sync"
)

type TCPPeer struct {
	conn net.Conn
	//if we dailed a connection => outbound = true
	//if we accept and retrieve a connection => outbound = false
	outbound bool
}

func newTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPtransport struct {
	listenAddress string
	listener      net.Listener
	handshakeFunc Handshake
	decoder       Decoder
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPtransport {
	return &TCPtransport{
		handshakeFunc: NoHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPtransport) ListerAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	go t.AcceptLoop()

	return nil
}

func (t *TCPtransport) AcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}

		go t.HandlerConn(conn)
	}
}

type temp struct{}

func (t *TCPtransport) HandlerConn(conn net.Conn) {
	peer := newTCPPeer(conn, true)

	if err := t.handshakeFunc(peer); err != nil {
		
	}
	msg := &temp{}
	for {
		if err := t.decoder.Decoder(conn, msg); err != nil {
			fmt.Printf("TCP Error %v\n", err)
			continue
		}
	}

}
