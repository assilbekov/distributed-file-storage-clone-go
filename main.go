package main

import (
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"log"
)

func main() {
	tr := p2p.NewTCPTransport(":4000")

	if err := tr.ListenAndAccept(); err != nil {
		log.Fatalf("failed to listen and accept: %v", err)
	}

	select {}
}
