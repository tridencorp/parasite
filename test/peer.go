package test

import (
	"parasite/p2p"
)

// Testing peers.
type Peer struct {
	Handler  chan *p2p.Msg
	Failure  chan *p2p.Msg
	Response *p2p.Msg
}

// For testing purpose we are immediately sending back
// response.
func (p *Peer) Send(msg *p2p.Msg) {
	p.Handler <- p.Response
}

// Return N number of peers.
func Peers(num int, res *p2p.Msg, handler chan *p2p.Msg) (peers []*Peer) {
	for i:=0; i<num; i++ {
		peers = append(peers, &Peer{Response: res, Handler: handler})
	}
	return peers
}