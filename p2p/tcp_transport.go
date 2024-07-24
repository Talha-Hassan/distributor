package p2p

import (
	"errors"
	"fmt"
	"log"
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

// remote addr implements the peer interface and return the
// remote addr of it's underlying connection
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (t *TCPtransport) ListerAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddress)
	if err != nil {
		return err
	}
	go t.AcceptLoop()

	fmt.Printf("Server is Listening on Port => %s\n", t.ListenAddress)

	return nil
}

func (t *TCPtransport) AcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}
		if err != nil {
			fmt.Printf("TCP accept error: %s\n", err)
		}
		fmt.Printf("New Incoming Connection %+v\n", conn)
		go t.HandlerConn(conn, false)
	}
}

func (t *TCPtransport) Close() error {
	return t.listener.Close()
}

func (p *TCPPeer) send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}
func (t *TCPtransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal(err)
		return err
	}
	go t.HandlerConn(conn, true)
	return nil
}
func (t *TCPtransport) HandlerConn(conn net.Conn, outbound bool) {
	var err error

	defer func() {
		fmt.Printf("Dropping Connection due to : %s", err)
		conn.Close()
	}()

	peer := newTCPPeer(conn, outbound)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
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
