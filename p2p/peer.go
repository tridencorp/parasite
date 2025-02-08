package p2p

import (
	"math/rand/v2"
	"parasite/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
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
	for msg := range p.messages {
		_, err := p.conn.Write(msg.Code, msg.Data)

		if err != nil {
			log.Error("%s", err)
		}
	}
}

// Request bunch of block headers from peer.
func (p *Peer) GetBlockHeaders(start, amount, skip uint64) (uint64, error){
	reqId := rand.Uint64()

	req := GetBlockHeaders{
		Start:   start,
		Amount:  amount,
		Skip:    skip,
		Reverse: false,
	}

	data, err := rlp.EncodeToBytes([]any{reqId, req})
	if err != nil {
		return 0, err
	}

	p.Send(NewMsg(GetBlockHeadersMsg, data))
	return reqId, nil
}

// Request blocks from peer.
func (p *Peer) GetBlocks(headerHashes []common.Hash) (uint64, error) {
	reqId := rand.Uint64()

	data, err := rlp.EncodeToBytes([]any{reqId, headerHashes})
	if err != nil {
		return 0, nil
	}

	p.Send(NewMsg(GetBlockBodiesMsg, data))
	return reqId, nil
}
