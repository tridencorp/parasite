package main

import (
	"fmt"
	"parasite/key"
	"parasite/node"
	"parasite/p2p"
)

const nodePrv = "02f378c3758a72b90f0cb53e36f7873308ca1d7d450861939e163e615cc65dce"
const nodeID  = "enode://e806157dfc5e11365210e09ad4af4fb129024de4a8c97c7a6c834daf9567200f9e8d03a769c1d53f13286a643e295b6f38073e90a1833c1e06ef23cc402cfecb@127.0.0.1:30303?discport=00"

func main() {
	srcPrv, err := key.Private()
	if err != nil {
		fmt.Print(err)
	}

	peer, err := node.Connect(nodeID, srcPrv)
	if err != nil {
		fmt.Print(err)
	}

	StartPeerReader(peer)
}

// Main peer reader responsible for handling all p2p messages.
// @CLEAN: This could probably be moved to a different place.
func StartPeerReader(peer *p2p.Peer) {
	msg, err := peer.Read()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Print(msg.Code)
	fmt.Print(msg.Data)
	fmt.Print(msg.Size)
}
