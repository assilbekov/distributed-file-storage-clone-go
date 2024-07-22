package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
)

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListedAddr:    ":8080",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: onPeer func.
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       "8080_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
	}

	s := NewFileServer(fileServerOpts)

	if err := s.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	select {}
}
