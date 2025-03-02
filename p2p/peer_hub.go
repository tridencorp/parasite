package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"parasite/log"
)

// PeerHub is responsible for managing all peers.
type PeerHub struct {
	Peers    []*Peer
	PeerList []string
	PrvKey *ecdsa.PrivateKey

	// Failure channel, on which errors from peers are received.
	// Based on the message type PeerHub will decide what to do,
	// ex: disconnect, reconnect, ...
	failure chan *Msg

	// Main dispatcher for handling peer messages.
	dispatcher Dispatcher
}

func NewPeerHub(peerList []string, prv *ecdsa.PrivateKey) *PeerHub {
	response, failure := make(chan *Msg), make(chan *Msg) 
	dispatcher := NewDispatcher(response, failure) 

	return &PeerHub{
		PeerList: peerList, 
		PrvKey: prv, 
		dispatcher: dispatcher,
		failure: failure,
	}
}

// Connect all peers.
func (hub *PeerHub) ConnectAll() {
	for _, address := range hub.PeerList {
		peer, err := Connect(address, hub.PrvKey)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		hub.Peers = append(hub.Peers, peer)

		// Run connected peer in goroutine and listen for incomming messages.
		go peer.StartWriter()
		go peer.StartReader(hub.dispatcher)
	}

	log.Debug("Connected Peers: %d", len(hub.Peers))
}

// Start the main PeerHub goroutine which is responsible
// for listening to messages from peers, ex: failure messages.
func (hub *PeerHub) Start() {
	select {
	case msg := <- hub.failure:
		fmt.Printf("FAILURE MSG: %v", msg)
	}
}
