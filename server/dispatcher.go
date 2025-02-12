package server

import (
	"parasite/log"
	"parasite/p2p"
)

// Dispatching received messages to designated handlers.
type Dispatcher struct {
	peer *p2p.Peer

	handler chan p2p.Msg
	failure chan p2p.Msg
}

// Create new Dispatcher.
func NewDispatcher(peer *p2p.Peer, handler chan p2p.Msg, failure chan p2p.Msg) *Dispatcher {
	return &Dispatcher{peer, handler, failure}
}

// Main dispatcher responsible for dispatching all incomming messages.
// It uses 2 channels: one for normal message handling and another one 
// for sending errors.
// Dispatcher is called by peer each time new message arrives.
func (dispatcher *Dispatcher) Dispatch(msg p2p.Msg) { 
	if msg.Code == p2p.PingMsg {
		dispatcher.peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))
		return
	}

	if msg.Code == p2p.BlockHeadersMsg {
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.DiscMsg {
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.NewPooledTransactionHashesMsg { 
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.TransactionsMsg {
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.BlockBodiesMsg { 
		dispatcher.handler <- msg
		return 
	}

	if msg.Code == p2p.ReceiptsMsg { 
		dispatcher.handler <- msg
		return 
	}

	if msg.Code == p2p.GetBlockHeadersMsg {
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.GetBlockBodiesMsg {
		dispatcher.handler <- msg
		return
	}

	if msg.Code == p2p.GetReceiptsMsg { 
		dispatcher.handler <- msg
		return 
	}

	// If we are here then we have unsupported message. 
	// Just print it for now.
	log.Error("Unknown msg code: %d\n", msg.Code)
}
