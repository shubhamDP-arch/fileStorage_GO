package main

import (
	"fileStorage/p2p"
	"fmt"
	"log"
)

func OnPeer(peer p2p.Peer)error{
	peer.Close()
	return nil
}
func main() {
	tcpOpts := p2p.TCPTransportOps{
		ListenAdder:   ":4000",
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer: OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}
	select {}

}
