package main

import (
	"crypto/ecdsa"
	"fmt"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/common"
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

	err = log.Configure(&log.Config{})
	if err != nil {
		fmt.Print(err)
	}

	// Starting log. All logs will go to it.
	go log.Start()

	log.Info("Connecting to peer ...")
	peer, err := node.Connect(nodes[0], srcPrv)
	if err != nil {
		log.Error("Cannot connect to peer:\n%s. \nError: %s\n", nodes[0], err)
	}

	// Sync playground
	// snc := flag.Bool("sync", false, "")
	// flag.Parse()

	StartPeer(peer, srcPrv)
}

func StartPeer(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
  for {
    msg, err := peer.Read()
    if err != nil {
      fmt.Print(err)
      break
    }

    // (0) HandshakeMsg
    if msg.Code == p2p.HandshakeMsg {
      err := p2p.HandleHandshake(msg, peer, &srcPrv.PublicKey)
      if err != nil {
        fmt.Print(err)
      }
      
      continue
    }

    // (16) StatusMsg
    if msg.Code == p2p.StatusMsg {
      err := p2p.HandleStatus(msg, peer)
      if err != nil {
        fmt.Print(err)
      }

      continue
    }

    // (2) PingMsg
    if msg.Code == p2p.PingMsg {
      peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))

      fmt.Println("Sending headers ...")
      _, err := peer.GetBlockHeaders(14678570, 1, 0)
      if err != nil {
        fmt.Print(err)
      }

      continue
    }

    // (19) GetBlockHeadersMsg
    if msg.Code == p2p.GetBlockHeadersMsg {
      fmt.Println(msg.Code)
      fmt.Println(msg.Data)    

      continue
    }

    // (20) BlockHeadersMsg
    if msg.Code == p2p.BlockHeadersMsg {
      log.Info("Get headers")

      headers, err := p2p.HandleBlockHeaders(msg)
      if err != nil {
        fmt.Print(err)
      }

      log.Info("Headers:\n%v", headers)
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

      continue
    }

    // (22) BlockBodiesMsg
    if msg.Code == p2p.BlockBodiesMsg {
      log.Info("!!! Get Blocks !!!")
      fmt.Println(msg.Code)
      fmt.Println(msg.Data)    

      continue
    }

    // (1) DiscMsg
    if msg.Code == p2p.DiscMsg {
      log.Error("!!! DISCONECT FROM NODE !!!")

      type DiscReason uint8
      var disc []DiscReason

      rlp.DecodeBytes(msg.Data, &disc)
      log.Error("Disconnect from peer: %s", p2p.DiscReasons[disc[0]])

      continue
    }

    // If we are here then we have unsupported message. 
    // Just print it for now.
    fmt.Printf("Unsupported msg code: %d\n", msg.Code)
    fmt.Printf(string(msg.Data))
	}
}
