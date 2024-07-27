package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
	"time"
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
		BootstrapNodes:    []string{":4000"},
	}

	s := NewFileServer(fileServerOpts)

	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
