package p2p

import (
	"crypto/ecdsa"
	"parasite/log"
)

// PeerHub is responsible for managing all peers.
type PeerHub struct {
	Peers    []*Peer
	PeerList []string
	PrivKey *ecdsa.PrivateKey

	// Main dispatcher for handling peer messages.
	dispatcher Dispatcher
}

func NewPeerHub(peerList []string, dispatcher Dispatcher, prv *ecdsa.PrivateKey) *PeerHub {
	return &PeerHub{ PeerList: peerList, PrivKey: prv, dispatcher: dispatcher }
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

	log.Debug("Nodes Connected: %d", len(hub.PeerList))
}
