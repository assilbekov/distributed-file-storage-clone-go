package main

import (
	"fmt"
	"github.com/assilbekov/distributed-file-storage-clone-go/p2p"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

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
	}
}

func (s *FileServer) Stop() {
	close(s.quitch)
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

	if len(s.BootstrapNodes) != 0 {
		fmt.Printf("bootstrapping network with %v\n", s.BootstrapNodes)
		s.bootstrapNetwork()
	}

	s.loop()

	return nil
}
