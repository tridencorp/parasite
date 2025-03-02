package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"
)

// Cluster will manage peers.
type Cluster struct {
	Peers []*Peer
	PeerList []string
}

func (cluster *Cluster) Add(peers []*Peer) {
	cluster.Peers = peers
}

// Load list with known peers.
func (cluster *Cluster) Load(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	cluster.PeerList = strings.Split(string(data), "\n")
	return nil
}

// Connect given number of peers. We must load peer list first (or add them by hand).
func (cluster *Cluster) Connect(num int, key *ecdsa.PrivateKey) {
	for _, enode := range cluster.PeerList {
		peer, err := Connect(enode, key)
		if err != nil {
			fmt.Println(err)
		}

		if peer != nil && len(cluster.Peers) <= num {
			cluster.Peers = append(cluster.Peers, peer)
		}
	}
}
