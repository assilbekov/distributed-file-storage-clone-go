package p2p

import "net"

// Message holds any arbitrary data that can be sent over the network.
type Message struct {
	From    net.Addr
	Payload []byte
}
