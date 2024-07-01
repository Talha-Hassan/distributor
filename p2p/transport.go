package p2p

//Peer is an interface that represent the remote node
type Peer interface {
}

//handles the communications between nodes between network
type Transport interface {
	ListerAndAccept() error
}
