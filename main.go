package main

import (
	"bytes"
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

	s := NewFileServer(fileServerOpts)
	tcpTransportOpts.OnPeer = s.onPeer

	return s
}

func main() {
	s1 := makeServer("localhost:8080", "")
	s2 := makeServer("localhost:8081", ":8080")
	go func() {
		log.Fatal(s1.Start())
	}()

	s2.Start()

	data := bytes.NewReader([]byte("My new big data file"))
	s2.StoreFile("key", data)
}
