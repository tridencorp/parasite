package p2p

import (
	"parasite/log"

	"github.com/ethereum/go-ethereum/p2p/rlpx"
)

// Main struct handling P2P communication.
type Peer struct {
	conn     *rlpx.Conn
	messages chan Msg
}

// Return new peer.
func NewPeer(conn *rlpx.Conn) *Peer {
	return &Peer{conn: conn, messages: make(chan Msg, 100)}
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

// Send msg to peer messages channel.
func (p *Peer) Send(msg Msg) {
	p.messages <- msg
}

// Start peer writer goroutine. There should be only one writer per peer.
func (p *Peer) StartWriter() {
	// We are iterating over each message and writing it sequentially 
	// to the TCP socket.
	for msg := range p.messages {
		_, err := p.conn.Write(msg.Code, msg.Data)

		if err != nil {
			log.Error("%s", err)
		}
	}
}
