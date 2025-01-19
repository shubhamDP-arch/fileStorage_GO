package p2p

import (
	"fmt"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}
func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOps struct {
	ListenAdder string

	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}
type TCPTransport struct {
	TCPTransportOps
	rpcch   chan RPC
	listner net.Listener
}

func NewTCPTransport(opts TCPTransportOps) *TCPTransport {
	return &TCPTransport{
		TCPTransportOps: opts,
		rpcch:           make(chan RPC),
	}
}

func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcch
}
func (t *TCPTransport) ListenAndAccept() error {
	var err error
	ln, err := net.Listen("tcp", t.ListenAdder)
	if err != nil {
		return err
	}
	t.listner = ln
	go t.startAcceptLoop()
	return nil
}

func (t *TCPTransport) startAcceptLoop() {
	for {
		conn, err := t.listner.Accept()
		if err != nil {
			fmt.Printf("TCP transport error accept error: %s\n", err)
		}
		fmt.Printf("new incoming connection %v\n", conn)

		go t.handleConn(conn)
	}

}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error
	defer func() {
		fmt.Printf("dropping peer connection: %s\n", err)
		conn.Close()
	}()
	peer := NewTCPPeer(conn, true)
	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}
	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
	}

	rpc := RPC{}
	// buf := make([]byte, 2000)
	for {
		// n, err := conn.Read(buf)
		// if err != nil {
		// 	fmt.Printf("TCP error %s\n", err)
		// }
		err := t.Decoder.Decode(conn, &rpc); 
			if err != nil {
				return
			}
			fmt.Printf("TCP error in loop %s\n", err)
		
		rpc.From = conn.RemoteAddr()
		t.rpcch <- rpc

	}

}
