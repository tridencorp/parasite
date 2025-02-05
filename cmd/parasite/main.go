package main

import (
	"crypto/ecdsa"
	"fmt"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
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

func GetBlockHeadersByNumber(reqId, start, amount, skip uint64, reverse bool) []any {	
	return []any{reqId, []any{start, amount, skip, reverse}}
}

func StartPeer(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
	// Request bunch of block headers
	reqId    := uint64(666)
	start    := uint64(14678570)
	amount   := uint64(1)
	skip     := uint64(0)
	reverse  := false

	data, _   := rlp.EncodeToBytes(GetBlockHeadersByNumber(reqId, start, amount, skip, reverse))
	headerMsg := p2p.NewMsg(p2p.GetBlockHeadersMsg, data)

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
			size, err := peer.Send(headerMsg)
			if err != nil {
				fmt.Print(err)
				return
			}
			fmt.Printf("they are send: %d\n", size)

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

		default:
			fmt.Printf("Unsupported msg code: %d\n", msg.Code)
			fmt.Printf(string(msg.Data))
		}
	}
}
