package server

import (
	"parasite/log"
	"parasite/p2p"
)

// Dispatching received messages to designated handlers.
type Dispatcher struct {
	Handler chan p2p.Msg
	Failure chan p2p.Msg
}

// Create new Dispatcher.
func NewDispatcher() *Dispatcher {
	return &Dispatcher{make(chan p2p.Msg, 1), make(chan p2p.Msg, 10)}
}

func (dispatcher *Dispatcher) Channels() (chan p2p.Msg, chan p2p.Msg){
	return dispatcher.Handler, dispatcher.Failure
}

// Main dispatcher responsible for dispatching all incomming messages.
// It uses 2 channels: one for normal message handling and another one 
// for sending errors.
// Dispatcher is called by peer each time new message arrives.
func (dispatcher *Dispatcher) Dispatch(peer *p2p.Peer, msg p2p.Msg) { 
	if msg.Code == p2p.PingMsg {
		peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))
		return
	}

	if msg.Code == p2p.BlockHeadersMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == p2p.DiscMsg {
		dispatcher.Failure <- msg
		return
	}

	if msg.Code == p2p.NewPooledTransactionHashesMsg { 
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == p2p.TransactionsMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == p2p.BlockBodiesMsg { 
		dispatcher.Handler <- msg
		return 
	}

	if msg.Code == p2p.ReceiptsMsg { 
		dispatcher.Handler <- msg
		return 
	}

	if msg.Code == p2p.GetBlockHeadersMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == p2p.GetBlockBodiesMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == p2p.GetReceiptsMsg { 
		dispatcher.Handler <- msg
		return 
	}

	// If we are here then we have unsupported message. 
	// Just print it for now.
	log.Error("Unknown msg code: %d\n", msg.Code)
}
