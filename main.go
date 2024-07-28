package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
	"time"
)

func makeServer(listenAddr, root string, nodes ...string) *FileServer {
	tcpTransportOpts := p2p.TCPTransportOpts{
		ListedAddr:    listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// TODO: onPeer func.
	}
	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: CASPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}

	return NewFileServer(fileServerOpts)
}

func main() {

	go func() {
		time.Sleep(5 * time.Second)
		s.Stop()
	}()

	if err := s.Start(); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
