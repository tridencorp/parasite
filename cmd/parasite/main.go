package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	// Load nodes
	var nodes []string
	err := config.Load("./eth_nodes.json", "nodes", &nodes)
	if err != nil {
		fmt.Print(err)
	}

	srcPrv, err := key.Private()
	if err != nil {
		fmt.Print(err)
	}

	err = log.Setup("parasite.log")
	if err != nil {
		fmt.Print(err)
	}

	log.Info("Connecting to peer ...")
	peer, err := node.Connect(nodes[0], srcPrv)
	if err != nil {
		fmt.Print(err)
	}

	// Sync playground 
	// snc := flag.Bool("sync", false, "")
	// flag.Parse()

	StartPeer(peer, srcPrv)
}

func StartPeer(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
	bytes, err := hex.DecodeString("e55f48b9e8733758dc96abc807f0e398780c95e221053d1664a364b6ff770b40")
	fmt.Print(err)
	fmt.Print(bytes)
	fmt.Print(len(bytes))

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

		case p2p.PingMsg:
			peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))

			fmt.Println("Sending headers ...")
			reqId, err := peer.GetBlockHeaders(14678570, 1, 0)
			if err != nil {
				fmt.Print(err)
			}

			fmt.Print(reqId)

		case p2p.GetBlockHeadersMsg:
			fmt.Println(msg.Code)
			fmt.Println(msg.Data)

		case p2p.BlockHeadersMsg:
			log.Info("Get headers")

			headers, err := p2p.HandleBlockHeaders(msg)
			if err != nil {
				fmt.Print(err)
			}

			log.Info("Headers:\n%V", headers)
			hh, err := headers[0].Hash()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Print(hh)
			log.Info("Header Hash:\n%v", hh)

			// Request block
			log.Info("Requesting blocks")
			_, err = peer.GetBlocks([]common.Hash{hh})
			if err != nil {
				fmt.Println(err)
			}

		case p2p.BlockBodiesMsg:
			log.Info("!!! Get Blocks !!!")
			fmt.Println(msg.Code)
			fmt.Println(msg.Data)

		default:
			fmt.Printf("Unsupported msg code: %d\n", msg.Code)
			fmt.Printf(string(msg.Data))
		}
	}
}
