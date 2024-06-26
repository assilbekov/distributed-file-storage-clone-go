package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
)

func main() {
	tcpOpts := p2p.TCPTransportOpts{
		ListedAddr:    ":4000",
		Decoder:       p2p.GOBDecoder{},
		HandshakeFunc: p2p.NOPHandshakeFunc,
	}
	tr := p2p.NewTCPTransport(tcpOpts)

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %v", err)
	}

	select {}
}
