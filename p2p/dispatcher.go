package p2p

import (
	"parasite/log"
)

// Dispatcher interface. It should dispatch messages to proper handlers.
type Dispatcher interface {
	Channels() (chan Msg, chan Msg)
	Dispatch(peer *Peer, msg Msg)
}

// Dispatching received messages to designated handlers.
type MsgDispatcher struct {
	Handler chan Msg
	Failure chan Msg
}

// Create new Dispatcher.
func NewDispatcher() *MsgDispatcher {
	return &MsgDispatcher{make(chan Msg, 1), make(chan Msg, 10)}
}

func (dispatcher *MsgDispatcher) Channels() (chan Msg, chan Msg){
	return dispatcher.Handler, dispatcher.Failure
}

// Main dispatcher responsible for dispatching all incomming messages.
// It uses 2 channels: one for normal message handling and another one 
// for sending errors.
// Dispatcher is called by peer each time new message arrives.
func (dispatcher *MsgDispatcher) Dispatch(peer *Peer, msg Msg) { 
	if msg.Code == PingMsg {
		peer.Send(NewMsg(PongMsg, []byte{}))
		return
	}

	if msg.Code == BlockHeadersMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == DiscMsg {
		dispatcher.Failure <- msg
		return
	}

	if msg.Code == NewPooledTransactionHashesMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == TransactionsMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == BlockBodiesMsg { 
		dispatcher.Handler <- msg
		return 
	}

	if msg.Code == ReceiptsMsg { 
		dispatcher.Handler <- msg
		return 
	}

	// We are parasite.
	// We are not responding to Get requests. 
	// If one arrives, we will return an empty response.
	// 
	// Nodes like this are helpfull to decrease load on
	// whole ethereum network - instead of sending rpc
	// requests to other nodes we will be calling our 
	// own node. It's good for dApps.
	// 
	// TODO: return empty response.
	if msg.Code == GetBlockHeadersMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == GetBlockBodiesMsg {
		dispatcher.Handler <- msg
		return
	}

	if msg.Code == GetReceiptsMsg { 
		dispatcher.Handler <- msg
		return 
	}

	// If we are here then we have unsupported message. 
	// Just print it for now.
	log.Error("Unknown msg code: %d\n", msg.Code)
}
