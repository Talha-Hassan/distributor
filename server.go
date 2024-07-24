package main

import (
	"distributer/p2p"
	"fmt"
	"log"
	"sync"
)

type FileServeropt struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNetwork  []string
}

type FileServer struct {
	FileServeropt
	store    *Store
	quitch   chan struct{}
	peers    map[string]p2p.Peer
	peerLock sync.Mutex
}

func newFileServer(opts FileServeropt) *FileServer {
	storeOpts := StoreOps{
		root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}
	return &FileServer{
		FileServeropt: opts,
		store:         NewStore(storeOpts),
		quitch:        make(chan struct{}),
		peers:         make(map[string]p2p.Peer),
	}
}

func (f *FileServer) stop() {
	close(f.quitch)
}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.peerLock.Lock()
	defer s.peerLock.Unlock()
	s.peers[p.RemoteAddr().String()] = p
	log.Printf("Connected with remote %s", p.RemoteAddr())
	return nil
}

func (f *FileServer) loop() {
	defer func() {
		log.Panicln("Server is Quiting")
		f.Transport.Close()
	}()
	for {
		select {
		case msg := <-f.Transport.Consume():
			fmt.Println(msg)
		case <-f.quitch:
			return
		}
	}
}

func (s *FileServer) bootStrapNetwork() error {
	for _, addr := range s.BootstrapNetwork {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("Attempting to connect with remote: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("Dial Error: ", err)
			}
		}(addr)
	}
	return nil
}

func (s *FileServer) start() error {
	if err := s.Transport.ListerAndAccept(); err != nil {
		return err
	}
	if len(s.BootstrapNetwork) != 0 {
		s.bootStrapNetwork()
	}
	s.loop()
	return nil
}
