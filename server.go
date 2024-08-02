package main

import (
	"fmt"
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
	"io"
	"log"
	"sync"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	peerLock sync.Mutex
	peers    map[string]p2p.Peer

	store  *Store
	quitch chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServerOpts: opts,
		store:          NewStore(storeOpts),
		quitch:         make(chan struct{}),
		peers:          make(map[string]p2p.Peer),
	}
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Write the data to the store.
	// 2. Broadcast the data to all connected peers.

	return nil
}

func (s *FileServer) Stop() {
	close(s.quitch)
}

func (s *FileServer) onPeer(peer p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()

	s.peers[peer.RemoteAddr().String()] = peer

	log.Println("new peer connected", peer.RemoteAddr().String())

	return nil
}

func (s *FileServer) loop() {
	defer func() {
		fmt.Printf("shutting down server\n")
		s.Transport.Close()
	}()

	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println("received message", msg)
		case <-s.quitch:
			return
		}
	}
}

func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}

		fmt.Printf("attempting to bootstrap network with %s\n", addr)
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				fmt.Printf("failed to dial %s: %v\n", addr, err)
			}
		}(addr)
	}

	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()
	s.loop()

	return nil
}
