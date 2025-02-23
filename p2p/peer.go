package p2p

import (
	"parasite/log"

	"github.com/ethereum/go-ethereum/p2p/rlpx"
)

// Main struct handling P2P communication.
type Peer struct {
	conn     *rlpx.Conn
	messages chan Msg

	Response chan Msg // Default channel to which we will send the response.
	Failure  chan Msg // Default channel to which we will send the failures.

	// Track requested messages so we can verify if
	// an incoming message is the one we requested.
	// 
	// It will also enable us to find requested messages
	// and use correct handler for processing.
	// 
	// We will use requestId for finding req/res match.
	RequestedMsgs map[uint64]Msg
}

// Return new peer.
func NewPeer(conn *rlpx.Conn) *Peer {
	return &Peer{
		conn: conn,
		messages: make(chan Msg, 100),
		RequestedMsgs: make(map[uint64]Msg),
	}
}

func (peer *Peer) GetBlockHeaders(number, amount uint64) error {
	msg, err := BlockHeadersReqMsg(number, amount, 0, false)
	if err != nil {
		return err
	}

	peer.Send(msg)
	return nil
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
	p.RequestedMsgs[msg.ReqId] = msg
	p.messages <- msg
}

// Start peer writer goroutine. 
// There should be only one writer per peer.
func (p *Peer) StartWriter() {
	// We are iterating over each message 
	// and writing it sequentially to socket.
	for msg := range p.messages {
		_, err := p.conn.Write(msg.Code, msg.Data)

		// @TODO: Figure out how to handle errors.
		if err != nil {
			log.Error("%s", err)
		}
	}
}

// Dispatcher interface. It should dispatch messages to proper handlers.
type Dispatcher interface {
	Channels() (chan Msg, chan Msg)
	Dispatch(peer *Peer, msg Msg)
}

// Start the peer reader, which will read messages 
// sequentially and send them to dispatcher.
func (p *Peer) StartReader(dispatcher Dispatcher) { 
	for {
		msg, err := p.Read()
		if err != nil {
			log.Error("%v", err)
			break
		}

		dispatcher.Dispatch(p, msg)
	}
}
