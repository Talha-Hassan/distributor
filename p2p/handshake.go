package p2p

type Handshake func(Peer) error

func NoHandshakeFunc(Peer) error { return nil }
