package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	opts := TCPTransportOpts{
		ListedAddr:    ":4000",
		HandshakeFunc: NOPHandshakeFunc,
		Decoder:       DefaultDecoder{},
	}
	tr := NewTCPTransport(opts)

	assert.Equal(t, opts.ListedAddr, tr.ListedAddr)

	// Start the transport
	// tr.Start()
	assert.Nil(t, tr.ListenAndAccept())

	select {
	//
	}
}
