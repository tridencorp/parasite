package main

import (
	"crypto/ecdsa"
	"fmt"
	"parasite/key"
	"parasite/node"
	"parasite/p2p"
)

const dstID = "enode://e806157dfc5e11365210e09ad4af4fb129024de4a8c97c7a6c834daf9567200f9e8d03a769c1d53f13286a643e295b6f38073e90a1833c1e06ef23cc402cfecb@127.0.0.1:30303?discport=00"

func main() {
	srcPrv, err := key.Private()
	if err != nil {
		fmt.Print(err)
	}

	peer, err := node.Connect(dstID, srcPrv)
	if err != nil {
		fmt.Print(err)
	}

	StartPeerReader(peer, srcPrv)
}

// Main peer reader responsible for handling all p2p messages.
func StartPeerReader(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
	for {
		msg, err := peer.Read()
		if err != nil {
			fmt.Print(err)
			break
		}

		switch msg.Code {
		case p2p.HandshakeMsg:
			err := p2p.HandleHandshake(msg, peer, &srcPrv.PublicKey)
			if err != nil {
				fmt.Print(err)
			}

		case p2p.StatusMsg:
			err := p2p.HandleStatus(msg, peer)
			if err != nil {
				fmt.Print(err)
			}

		default:
			fmt.Printf("Unsupported msg code: %d", msg.Code)
			fmt.Printf(string(msg.Data))
		}
	}
}
