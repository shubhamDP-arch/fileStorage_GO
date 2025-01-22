package main

import (
	"fileStorage/p2p"
	"fmt"
	"log"
)

type FileServerOpts struct {
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Trasport
	bootstrapnodes    []string
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

func (s *FileServer) BootStrapNetwork() error {
	for _, addr := range s.bootstrapnodes {
		go func(addr string) {
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial, error", err)
			}
		}(addr)
		return nil
	}
	return nil
}
func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}
	s.BootStrapNetwork()
	s.loop()
	return nil
}

func (s *FileServer) loop() {
	defer func() {
		log.Println("file server stopped user quit action")
		s.Transport.Close()
	}()
	for {
		select {
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)
		case <-s.quitch:
			return
		}
	}

}
