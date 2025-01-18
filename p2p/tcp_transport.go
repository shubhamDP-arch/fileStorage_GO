package p2p

import (
	"bytes"
	"fmt"
	"net"
	"sync"
)
type TCPPeer struct{
	conn net.Conn
	outbound bool
}
func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn: conn,
		outbound: outbound,
	}
}
type TCPTransport struct {
	listenAddress string
	listener   net.Listener
	shakeHands HandshakeFunc
	decoder 	Decoder
	mu         sync.RWMutex
	peers      map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport {
	return &TCPTransport{
		shakeHands: NOPHandshakeFunc,
		listenAddress: listenAddr,
	}
}

func (t *TCPTransport)ListenAndAccept()error  {
	var err error
	ln, err := net.Listen("tcp", t.listenAddress)
	if err != nil {
		return err
	}
	t.listener = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport)startAcceptLoop()  {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP transport error accept error: %s\n", err)
		}	
		fmt.Printf("new incoming connection %v\n", conn)

		go t.handleConn(conn)
	}
	
}
type Temp struct{}

func (t* TCPTransport)handleConn(conn net.Conn)  {
	peer := NewTCPPeer(conn, true)
	if err := t.shakeHands(peer); err != nil {
		
	}

	msg := &Temp{}
	for{
		if err := t.decoder.Decode(conn, msg); err != nil{
			fmt.Printf("TCP error in loop %s\n", err)
		}
	}

}
