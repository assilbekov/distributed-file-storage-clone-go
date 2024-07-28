package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
)

func makeServer(listenAddr string, nodes ...string) *FileServer {
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
	s1 := makeServer("localhost:8080", "")
	go func() {
		log.Fatal(s1.Start())
	}()
}
