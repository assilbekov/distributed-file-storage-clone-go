package p2p

import "net"

// Peer is an interface that represents a node in the network.
type Peer interface {
	net.Conn
	Send([]byte) error
}

// Transport is an interface that handles the communication between the nodes in the network.
// This can be a (TCP, UDP, websockets, etc.) connection.
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
