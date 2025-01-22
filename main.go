package main

import (
	"fileStorage/p2p"
	"log"
)

func makeServer(listenAddr string, nodes ...string)*FileServer  {
	tcptransportOpts := p2p.TCPTransportOps{
		ListenAdder: listenAddr,
		HandshakeFunc: p2p.NOPHandshakeFunc,
		Decoder: p2p.DefaultDecoder{},
	}
	tcpTransport := p2p.NewTCPTransport(tcptransportOpts)
	fileServerOpts := FileServerOpts{
		StorageRoot: listenAddr + "_network",
		PathTransformFunc: CASPathTransformfunc,
		Transport: tcpTransport,
		bootstrapnodes: nodes,
	}
	return NewFileServer(fileServerOpts)
}
func main() {
	s1 := makeServer(":3000", "")
	go func() {s1.Start()}()

}
