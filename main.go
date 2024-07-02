package main

import (
	"fmt"
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListedAddr:    ":4000",
		Decoder:       p2p.DefaultDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
		OnPeer: func(peer p2p.Peer) error {
			fmt.Printf("Peer connected: %s\n", peer)
			return nil
		},
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("Received message: %s\n", string(msg.Payload))
		}
	}()

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %v", err)
	}

	select {}
}
