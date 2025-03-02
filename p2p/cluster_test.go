package p2p

import (
	"fmt"
	"testing"
)

func TestCluster(t *testing.T) {
	cluster := Cluster[*Peer]{}
	cluster.Load("./bootnodes.txt")

	fmt.Println(len(cluster.PeerList))
}