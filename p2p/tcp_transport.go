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

type TCPTransportOpts struct {
	ListedAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
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

type Temp struct {
	//
}

func (t *TCPTransport) handleConnection(conn net.Conn) {
	peer := NewTCPPeer(conn, true)

	if err := t.shakeHands(peer); err != nil {
		conn.Close()
		// Handle error
		fmt.Printf("TCP handshake error: %v\n", err)
		return
	}

	// Read loop
	msg := &Temp{}
	for {
		if err := t.decoder.Decode(conn, msg); err != nil {
			// Handle error
			fmt.Printf("TCP error decoding message: %v\n", err)
			return
		}
	}

	fmt.Println("New incoming peer: ", peer)
	fmt.Printf("Handling connection from %v\n", conn.RemoteAddr())
	// Read the message
	// Decode the message
	// Handle the message
}
