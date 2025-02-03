package main

import (
	"crypto/ecdsa"
	"flag"
	"fmt"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
)

const dstID = "enode://e806157dfc5e11365210e09ad4af4fb129024de4a8c97c7a6c834daf9567200f9e8d03a769c1d53f13286a643e295b6f38073e90a1833c1e06ef23cc402cfecb@127.0.0.1:30303?discport=00"

func main() {
	srcPrv, err := key.Private()
	if err != nil {
		fmt.Print(err)
	}

	err = log.Setup("parasite.log")
	if err != nil {
		fmt.Print(err)
	}

	log.Info("Connecting to peer ...")
	peer, err := node.Connect(dstID, srcPrv)
	if err != nil {
		fmt.Print(err)
	}

	// Sync playground 
	sync := flag.Bool("sync", false, "")
	flag.Parse()

	if *sync {
		StartSync(peer, srcPrv)
	}

	// Normal flow
	StartPeerReader(peer, srcPrv)
}

func GetBlockHeadersByNumber(reqId, start, amount, skip uint64, reverse bool) []any {
	return []any{reqId, []any{start, amount, skip, reverse}}
}

func StartSync(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
	// Request bunch of block headers
	reqId    := uint64(666)
	start    := uint64(14678570)
	amount   := uint64(1)
	skip     := uint64(0)
	reverse  := false

	data, _ := rlp.EncodeToBytes(GetBlockHeadersByNumber(reqId, start, amount, skip, reverse))
	headerMsg := p2p.NewMsg((p2p.BaseProtocolLen+3), data)

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
			fmt.Println(msg.Code)
			fmt.Println(msg.Data)

		default:
			fmt.Printf("Unsupported msg code: %d\n", msg.Code)
			fmt.Printf(string(msg.Data))
		}
	}

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

		case p2p.PingMsg:

		default:
			fmt.Printf("Unsupported msg code: %d", msg.Code)
			fmt.Printf(string(msg.Data))
		}
	}
}
