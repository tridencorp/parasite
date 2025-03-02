package p2p

import (
	"fmt"
	"parasite/key"
	"testing"
)

func TestCluster(t *testing.T) {
	prv, _ := key.Private()

	cluster := Cluster{}
	cluster.Load("./bootnodes.txt")

	cluster.Connect(2, prv)

	fmt.Println(cluster.Get(2))
}
