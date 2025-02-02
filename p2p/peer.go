package p2p

import "github.com/ethereum/go-ethereum/p2p/rlpx"

// Main struct handling P2P communication.
type Peer struct {
	conn *rlpx.Conn
}

// Return new peer.
func NewPeer(conn *rlpx.Conn) *Peer {
	return &Peer{conn: conn}
}
