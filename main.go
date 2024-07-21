package main

import "github.com/assilbekov/distributed-file-storage-clone-go/p2p"

func main() {
	tcpTransportOpts := p2p.TCPTransportOpts{}
	tcpTransport := p2p.NewTCPTransport()
	fileServerOpts := FileServerOpts{
		ListenAddr:        ":8080",
		StorageRoot:       "8080_network",
		PathTransformFunc: CASPathTransformFunc,
	}
}
