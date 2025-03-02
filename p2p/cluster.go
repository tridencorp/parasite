package p2p

import (
	"os"
	"strings"
)

// Cluster will manage peers/nodes.
type Cluster[T Sender] struct {
	Peers []T
	PeerList []string
}

func (cluster *Cluster[T]) Add(peers []T) {
	cluster.Peers = peers
}

// Load list with known peers.
func (cluster *Cluster[T]) Load(filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	cluster.PeerList = strings.Split(string(data), "\n")
	return nil
}
