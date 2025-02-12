package p2p

import (
	"crypto/ecdsa"
)

// PeerHub is responsible for managing all peers.
type PeerHub struct {
	Peers []*Peer

	// Main dispatcher for handling peer messages.
	dispatcher Dispatcher
}

func NewPeerHub(peerList []string, dispatcher Dispatcher, prv *ecdsa.PrivateKey) *PeerHub {
	return &PeerHub{dispatcher: dispatcher}
}
