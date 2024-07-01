package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTCPTransport(t *testing.T) {
	listenaddr := ":4000"
	tr := NewTCPTransport(listenaddr)
	assert.Equal(t, tr.listenAddress, listenaddr)

	assert.Nil(t, tr.ListerAndAccept())
}
