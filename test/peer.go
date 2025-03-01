package test

import (
	"parasite/p2p"
)

// Testing peers.
type Peer struct {
	Handler  chan *p2p.Msg
	Failure  chan *p2p.Msg
}

func (p *Peer) Send(msg *p2p.Msg) {}

// Return N number of peers.
func Peers(num int) (peers []*Peer) {
	for i:=0; i<num; i++ {
		peers = append(peers, &Peer{})
	}
	return peers
}