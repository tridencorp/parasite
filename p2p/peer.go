package p2p

import (
	"parasite/log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
)

// Main struct handling P2P communication.
type Peer struct {
	conn     *rlpx.Conn
	messages chan *Msg

	Response chan *Msg // Default channel to which we will send the response.
	Failure  chan *Msg // Default channel to which we will send the failures.

	// Track requested messages so we can verify if
	// an incoming message is the one we requested.
	// 
	// It will also enable us to find requested messages
	// and use correct handler for processing.
	// 
	// We will use requestId for finding req/res match.
	RequestedMsgs map[uint64]*Msg
}

// Return new peer.
func NewPeer(conn *rlpx.Conn) *Peer {
	return &Peer{
		conn: conn,
		messages: make(chan *Msg, 100),
		RequestedMsgs: make(map[uint64]*Msg),
	}
}

func (peer *Peer) GetBlockBodies(hashes []common.Hash) error {
	msg, err := EncodeMsg(GetBlockBodiesMsg, hashes)
	if err != nil {
		return err
	}

	peer.Send(msg)
	return nil
}

func (peer *Peer) GetReceipts(hashes []common.Hash) error {
	msg, err := EncodeMsg(GetReceiptsMsg, hashes)
	if err != nil {
		return err
	}

	peer.Send(msg)
	return nil
}

func (peer *Peer) GetBlockHeaders(number, amount uint64) error {
	msg, err := EncodeMsg(GetBlockHeadersMsg, getBlockHeadersMsg{number, amount, 0, false})
	if err != nil {
		return err
	}

	peer.Send(msg)
	return nil
}

// Start writer and reader goroutines.
func (p *Peer) Start() {
	p.Response, p.Failure = make(chan *Msg, 1), make(chan *Msg, 1)

	go p.StartWriter()
	go p.StartReader(NewDispatcher(p.Response, p.Failure))
}

// Close peer connection.
// 
// TODO: Goroutines also should be closed.
// TODO: A lot of goroutines can use this,
// close it safely.
func (p *Peer) Close() error {
	return p.conn.Close()
}

// Reads message from a connected peer.
// BlocksBlocks until data is available.
func (p *Peer) Read() (*Msg, error) {
	code, data, _, err := p.conn.Read()
	if err != nil {
		return nil, err
	}

	return NewMsg(int(code), data), nil
}

// Send msg to peer messages channel.
func (p *Peer) Send(msg *Msg) {
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

// Start the peer reader, which will read messages 
// sequentially and send them to dispatcher.
func (p *Peer) StartReader(dispatcher Dispatcher) { 
	for {
		msg, err := p.Read()
		if err != nil {
			log.Error("%v", err)
			break
		}

		p.Response <- msg
	}
}

