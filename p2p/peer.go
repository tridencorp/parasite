package p2p

import (
	"github.com/ethereum/go-ethereum/p2p/rlpx"
)

// Main struct handling P2P communication.
type Peer struct {
	conn *rlpx.Conn
}

// Return new peer.
func NewPeer(conn *rlpx.Conn) *Peer {
	return &Peer{conn: conn}
}

// Reads message from a connected peer.
// BlocksBlocks until data is available.
func (p *Peer) Read() (Msg, error) {
	code, data, size, err := p.conn.Read()
	if err != nil {
		return Msg{}, err
	}

	return Msg{Code: code, Size: uint32(size), Data: data}, err
}

// Send message to peer.
func (p *Peer) Send(msg Msg) (uint32, error) {
	return p.conn.Write(msg.Code, msg.Data)
}
