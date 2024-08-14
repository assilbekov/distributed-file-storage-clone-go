package p2p

import (
	"errors"
	"fmt"
	"net"
)

// TCPPeer represents a node in the network.
type TCPPeer struct {
	// con is an underlying connection of the peer.
	net.Conn

	// if we dial and retrieve a connection, outbound is true
	// if we accept and retrieve a connection, outbound is false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		Conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Send(data []byte) error {
	_, err := p.Conn.Write(data)
	return err
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
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

// Close implements the Transport interface. It closes the underlying listener.
func (t *TCPTransport) Close() error {
	return t.listener.Close()
}

// Dial implements the Transport interface. It dials a connection to the given address.
func (t *TCPTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	go t.handleConnection(conn, true)

	return nil
}

func (t *TCPTransport) ListenAndAccept() error {
	ln, err := net.Listen("tcp", t.ListedAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.startAcceptLoop()

	fmt.Printf("Listening on %v\n", ln.Addr())

	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			// Handle error
			fmt.Printf("TCP error accepting connection: %v\n", err)
			continue
		}

		go t.handleConnection(conn, false)
	}
}

func (t *TCPTransport) handleConnection(conn net.Conn, outbound bool) {
	var err error
	defer func() {
		fmt.Printf("Closing connection from %v\n, with error: %v", conn.RemoteAddr(), err)
		if err := conn.Close(); err != nil {
			fmt.Printf("TCP error closing connection: %v\n", err)
		}
	}()

	peer := NewTCPPeer(conn, outbound)

	if err = t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	fmt.Printf("Handling connection from %v\n", conn.RemoteAddr())

	// Read loop
	rpc := RPC{}
	for {
		err = t.Decoder.Decode(conn, &rpc)
		if err := t.Decoder.Decode(conn, &rpc); err != nil {
			// Handle error
			fmt.Printf("TCP read error decoding message: %v\n", err)
			return
		}

		rpc.From = conn.RemoteAddr().String()
		t.rpcch <- rpc
	}

	// Read the message
	// Decode the message
	// Handle the message
}
