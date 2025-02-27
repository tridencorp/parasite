package test

import (
	"parasite/p2p"
)

// Testing peers.
type Peer struct {
	Handler  chan *p2p.Msg
	Failure  chan *p2p.Msg
}

func (p *Peer) Send(msg *p2p.Msg) {

}
