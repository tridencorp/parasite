package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/rand"
)

// Cluster will manage peers.
type Cluster struct {
	Peers []*Peer
	PeerList []string
}

// Add peers to cluster.
func (cluster *Cluster) Add(peers []*Peer) {
	cluster.Peers = append(cluster.Peers, peers...)
}

// Get random peers from cluster.
func (cluster *Cluster) Get(num int) (peers []*Peer) {
	if len(cluster.Peers) == 0 {return}
	if num > len(cluster.Peers) {num = len(cluster.Peers)}

	for i := 0; i < num; i++ {
		rand := rand.Intn(len(cluster.Peers))

		// TODO: check if peers are unique.
		peers = append(peers, cluster.Peers[rand])
	}

	return peers
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
