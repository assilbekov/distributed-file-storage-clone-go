package p2p

import (
	"fmt"
	"net"
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

// Close implements the Peer interface. It closes the underlying connection.
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListedAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	onPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcch    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcch:            make(chan RPC),
	}
}

// Consume implements the Transport interface. Will return a read-only channel of RPC messages.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListedAddr)
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
	var err error
	defer func() {
		fmt.Printf("Closing connection from %v\n, with error: %v", conn.RemoteAddr(), err)
		if err := conn.Close(); err != nil {
			fmt.Printf("TCP error closing connection: %v\n", err)
		}
	}()

	peer := NewTCPPeer(conn, true)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.onPeer != nil {
		if err = t.onPeer(peer); err != nil {
			return
		}
	}

	fmt.Printf("Handling connection from %v\n", conn.RemoteAddr())

	// Read loop
	rpc := RPC{}
	for {
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			// Handle error
			fmt.Printf("TCP error decoding message: %v\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc
	}

	// Read the message
	// Decode the message
	// Handle the message
}
