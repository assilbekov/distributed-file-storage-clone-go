package p2p

import "testing"

func TestTCPTranport(t *testing.T) {
	listenAddr := ":4000"
	tt := NewTCPTransport(listenAddr)
}
