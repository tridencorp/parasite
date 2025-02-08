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

  go StartPeerReader(peer, srcPrv)
  go peer.StartWriter()

  // Lets ask for block headers
  log.Info("Getting headers request")

  msg, err := p2p.BlockHeadersReq(14678570, 1, 0, false)
  if err != nil {
    log.Error("%s", err)
  }

  handler := make(chan p2p.Msg) 
  msg.Handler = handler
  fmt.Printf("ID: %d", msg.ReqId)

  peer.Send(msg)

  for msg := range handler {
    fmt.Println("!!! got headers from handler !!!")
    fmt.Printf("msg: %v", msg)
  }

  // Let's wait indefinitely for now.
  dummy := make(chan bool)
  <- dummy 
}

// Start peer reader. There should be only one reader per peer.
func StartPeerReader(peer *p2p.Peer, srcPrv *ecdsa.PrivateKey) {
  for {
    msg, err := peer.Read()
    if err != nil {
      log.Error("%v", err)
      break
    }

    // (2) PingMsg
    if msg.Code == p2p.PingMsg {
      log.Info("Got Ping")
      peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))
      continue
    }

    // (19) GetBlockHeadersMsg
    if msg.Code == p2p.GetBlockHeadersMsg {
      log.Info("%d", msg.Code)
      log.Info("%v", msg.Data)    
      continue
    }

    // (20) BlockHeadersMsg
    if msg.Code == p2p.BlockHeadersMsg {
      log.Info("Get headers")

      
      res, err := p2p.BlockHeadersRes(msg)
      if err != nil {
        log.Error("%v", err)
      }

      log.Error("%d", res.ReqId)
      reqMsg, exists := peer.RequestedMsgs[res.ReqId]
      if exists {
        log.Error("exists ........")
        reqMsg.Handler <- msg
      }

      log.Info("Headers:%v", res.Headers)
      hh, err := res.Headers[0].Hash()
      if err != nil {
        log.Error("%v", err)
      }

      // Request block
      log.Info("Requesting blocks")
      msg, err = p2p.BlocksReq([]common.Hash{hh})
      if err != nil {
        log.Error("%v", err)
      }

      peer.Send(msg)
      continue
    }

    // (22) BlockBodiesMsg
    if msg.Code == p2p.BlockBodiesMsg {
      log.Info("!!! Get Blocks !!!")
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
    log.Error("Unsupported msg code: %d\n", msg.Code)
  }
}
