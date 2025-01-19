package p2p

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestTCPTransport(t *testing.T)  {
	opts := TCPTransportOps{
		ListenAdder: ":6000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder: DefaultDecoder{},
	}

	tr := NewTCPTransport(opts)
	
	assert.Equal(t, tr.ListenAdder, ":6000")
	assert.Nil(t, tr.ListenAndAccept())
}