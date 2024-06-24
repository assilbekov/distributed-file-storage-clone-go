package p2p

import "errors"

// ErrInvalidHandshake is returned when the handshake between
// the local and remote node could not be established.
var ErrInvalidHandshake = errors.New("invalid handshake")

// HandshakeFunc is a function that is called when a new connection is established.
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(Peer) error {
	return nil
}
