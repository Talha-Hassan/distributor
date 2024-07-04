package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOps{
		ListenAddress: ":3000",
		HandshakeFunc: NoHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)
	assert.Equal(t, tr.ListenAddress, ":3000")

	assert.Nil(t, tr.ListerAndAccept())
}
