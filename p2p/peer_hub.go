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
	PrivKey *ecdsa.PrivateKey

	// Failure channel, on which errors from peers are received.
	// Based on the message type PeerHub will decide what to do,
	// ex: disconnect, reconnect, ...
	failure chan Msg

	// Main dispatcher for handling peer messages.
	dispatcher Dispatcher
}

func NewPeerHub(peerList []string, dispatcher Dispatcher, prv *ecdsa.PrivateKey) *PeerHub {
	_, failure := dispatcher.Channels()

	return &PeerHub{
		PeerList: peerList, 
		PrivKey: prv, 
		dispatcher: dispatcher,
		failure: failure,
	}
}

// Try to connect to all peers.
func (hub *PeerHub) ConnectAll() {
	for _, address := range hub.PeerList {
		peer, err := Connect(address, hub.PrivKey)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		hub.Peers = append(hub.Peers, peer)

		go peer.StartWriter()
		go peer.StartReader(hub.dispatcher)
	}

	log.Debug("Nodes Connected: %d", len(hub.Peers))
}

// Start the main PeerHub goroutine which is responsible
// for listening to messages from peers, ex: failure messages.
func (hub *PeerHub) Start() {
	select {
	case msg := <- hub.failure:
		fmt.Printf("FAILURE MSG: %v", msg)
	}
}
