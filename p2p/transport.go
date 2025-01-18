package p2p



//peer is an interface that represents a remote node
type Peer interface{
	
}

//transport is anything that handles the communication between nodes in the network
type Trasport interface{
	ListenAndAccept() error
	
}