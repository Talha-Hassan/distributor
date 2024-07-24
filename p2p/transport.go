package p2p

import "net"

//Peer is an interface that represent the remote node
type Peer interface {
	send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

//handles the communications between nodes between network
type Transport interface {
	ListerAndAccept() error
	Consume() <-chan RPC
	Close() error
	Dial(string) error
}
