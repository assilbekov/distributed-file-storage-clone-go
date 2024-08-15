package main

import (
	"bytes"
	"encoding/gob"
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

func (s *FileServer) broadcast(msg *DataMessage) error {
	peers := []io.Writer{}
	for _, peer := range s.peers {
		peers = append(peers, peer)
	}

	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)
}

type Message struct {
	Payload any
}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	// 1. Write the data to the store.
	// 2. Broadcast the data to all connected peers.

	buf := new(bytes.Buffer)
	msg := &Message{
		Payload: []byte("storagekey"),
	}
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.peers {
		if err := peer.Send(buf.Bytes()); err != nil {
			return err
		}
	}

	payload := []byte("THIS LARGE FILE")
	for _, peer := range s.peers {
		if err := peer.Send(payload); err != nil {
			return err
		}
	}

	return nil

	/*buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)

	if err := s.store.Write(key, tee); err != nil {
		return err
	}

	p := &DataMessage{
		Key:  key,
		Data: buf.Bytes(),
	}

	return s.broadcast(&Message{
		From:    "Todo",
		Payload: p,
	})*/
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
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Printf("failed to decode message: %v\n", err)
			}

			peer, ok := s.peers[rpc.From]
			if !ok {
				log.Printf("peer not found: %v\n", rpc.From)
				continue
			}

			b := make([]byte, 1024)
			if _, err := peer.Read(b); err != nil {
				log.Fatal("Couldn't read a buffer")
			}
			panic("panic to test")

			fmt.Printf("peer %+v\n", peer)

			fmt.Printf("received message %+v\n", msg)

			/*if err := s.handleMessage(&m); err != nil {
				log.Printf("failed to handle message: %v\n", err)
			}*/
		case <-s.quitch:
			return
		}
	}
}

/*func (s *FileServer) handleMessage(msg *Message) error {
	switch v := msg.Payload.(type) {
	case *DataMessage:
		fmt.Printf("received data %+v\n", v)
	}

	return nil
}*/

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
