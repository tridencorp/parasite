package server

import "parasite/p2p"

// Dispatching received messages to designated handlers.
type Dispatcher struct {
	handler chan p2p.Msg
	failure chan p2p.Msg
}

// Main dispatcher responsible for dispatching all incomming messages.
// It uses 2 channels: one for normal message handling and another one 
// for sending errors.
func Dispatch(msg *p2p.Msg) { 

}
