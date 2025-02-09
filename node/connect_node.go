package node

import (
	"crypto/ecdsa"
	"parasite/log"
	"parasite/p2p"
)

// Establish connections with nodes.
func ConnectNodes(nodes []string, prv *ecdsa.PrivateKey, dispatch p2p.Dispatch, handler chan p2p.Msg, failure chan p2p.Msg) {
	for _, n := range nodes {
		peer, err := Connect(n, prv)
		if err != nil {
			log.Error(err.Error())
			continue
		}

		log.Info("Peer Connected: %v", peer)

		go peer.StartWriter()
		go peer.StartReader(handler, handler, dispatch)
	}
}
