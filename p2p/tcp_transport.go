package p2p

import (
	"fmt"
	"net"
	"sync"
)

// TCPPeer represents a node in the network.
type TCPPeer struct {
	// con is an underlying connection of the peer.
	conn net.Conn

	// if we dial and retrieve a connection, outbound is true
	// if we accept and retrieve a connection, outbound is false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

type TCPTransport struct {
	listenAddr    string
	listener      net.Listener
	handshakeFunc HandshakeFunc

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		listenAddr: listenAddr,
		handshakeFunc: func(any) error {
			return nil
		},
		peers: make(map[net.Addr]Peer),
	}
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.startAcceptLoop()

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			// Handle error
			fmt.Printf("TCP error accepting connection: %v\n", err)
			continue
		}

		go t.handleConnection(conn)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	fmt.Println("New incoming peer: ", peer)
	fmt.Printf("Handling connection from %v\n", conn.RemoteAddr())
	// Read the message
	// Decode the message
	// Handle the message
}
