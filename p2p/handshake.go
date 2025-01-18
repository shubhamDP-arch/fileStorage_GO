package p2p

type HandshakeFunc func (Peer)error

func NOPHandshakeFunc(any)error  {
	return nil
}

