package p2p

import (
	"fmt"
	"net"
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

type TCPTransportOps struct {
	ListenAddress string
	HandshakeFunc Handshake
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPtransport struct {
	TCPTransportOps
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(data TCPTransportOps) *TCPtransport {
	return &TCPtransport{
		TCPTransportOps: data,
		rpcch:           make(chan RPC),
	}
}

func (t *TCPtransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPtransport) ListerAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
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

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// type temp struc t{}

func (t *TCPtransport) HandlerConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("Dropping Connection due to : %s", err)
	}()

	peer := newTCPPeer(conn, true)

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	fmt.Printf("Accpeted %+v\n", conn)
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP Handshake Error %s\n", err)
		return
	}
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			fmt.Printf("TCP Error %s\n", err)
			continue
		}
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

}
