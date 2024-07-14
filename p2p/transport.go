package p2p

//Peer is an interface that represent the remote node
type Peer interface {
	Close() error
}

//handles the communications between nodes between network
type Transport interface {
	ListerAndAccept() error
	Consume() <-chan RPC
}
