package main

import (
	"fmt"
	"os"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
)

func main() {
  // Display parasite logo
  displayLogo("./logo.txt")

  // Load nodes
  var nodes []string
  err := config.Load("./eth_nodes.json", "nodes", &nodes)
  if err != nil {
    fmt.Print(err)
  }

  err = log.Configure(&log.Config{})
  if err != nil {
    fmt.Print(err)
  }

  // Starting log. All logs will go to it.
  go log.Start()
  
  srcPrv, err := key.Private()
  if err != nil {
    fmt.Print(err)
  }

  
  // !!! TESTING PLAYGROUND !!!
  // 
  // log.Info("Connecting to peer ...")
  // peer, err := node.Connect(nodes[0], srcPrv)
  // if err != nil {
  //   log.Error("Cannot connect to peer:\n%s. \nError: %s\n", nodes[0], err)
  // }
  
  // msg, err := p2p.BlockHeadersReq(14678570, 1, 0, false)
  // if err != nil {
  //   log.Error("%s", err)
  // }
  
  handler := make(chan p2p.Msg) 
  // msg.Handler = handler
  
  // go peer.StartWriter()
  // go peer.StartReader(handler, handler, Dispatch)
  
  node.ConnectNodes(nodes, srcPrv, Dispatch, handler, handler)

  // Lets ask for block headers
  // log.Info("Block Headers request")
  // peer.Send(msg)

  for msg := range handler {
    fmt.Println("!!! got headers via handler !!!")
    fmt.Printf("msg: %v", msg)
  }

  // Let's wait indefinitely for now.
  dummy := make(chan bool)
  <- dummy 
}

// Dispatch incomming messages. Most of them will go to
// handler, failures can be redirect to different place, 
// like PeerHub.
func Dispatch(msg p2p.Msg, peer *p2p.Peer, handler chan p2p.Msg, failure chan p2p.Msg) {
  // (2) PingMsg
  if msg.Code == p2p.PingMsg {
    log.Info("!!! Got Ping !!!")
    peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))
    return
  }

  // (19) GetBlockHeadersMsg
  if msg.Code == p2p.GetBlockHeadersMsg {
    log.Info("%d", msg.Code)
    log.Info("%v", msg.Data)    
    return
  }

  // (20) BlockHeadersMsg
  if msg.Code == p2p.BlockHeadersMsg {
    log.Info("!!! Got headers !!!")

    res, err := p2p.BlockHeadersRes(msg)
    if err != nil {
      log.Error("%v", err)
    }

    // Send msg with requsted headers to our handler.
    reqMsg, exists := peer.RequestedMsgs[res.ReqId]
    if exists {
      reqMsg.Handler <- msg
    }

    return
  }

  // (1) DiscMsg
  if msg.Code == p2p.DiscMsg {
    log.Error("!!! DISCONECT FROM NODE !!!")

    type DiscReason uint8
    var disc []DiscReason

    rlp.DecodeBytes(msg.Data, &disc)
    log.Error("Disconnect from peer: %s", p2p.DiscReasons[disc[0]])
    return
  }

  if msg.Code == p2p.NewPooledTransactionHashesMsg { 
    log.Info("Request: %d : NewPooledTransactions", msg.Code)

    pooledTx, err := p2p.PooledTransactions(msg)
    if err != nil {
      fmt.Println(err)
    }
    
    fmt.Printf("%v", pooledTx)
    return
  }

  if msg.Code == p2p.TransactionsMsg {
    log.Info("Request: %d : p2p.TransactionsMsg", msg.Code)

    txs, err := p2p.NewTransactions(msg)
    if err != nil {
      fmt.Println(err)
    }

    fmt.Printf("%v", txs)
    return
  }

  // @TODO: Needs to be implemented
  if msg.Code == p2p.BlockBodiesMsg    { log.Error("Implement %d", p2p.BlockBodiesMsg)    ;return }
  if msg.Code == p2p.GetBlockBodiesMsg { log.Error("Implement %d", p2p.GetBlockBodiesMsg) ;return }
  if msg.Code == p2p.GetReceiptsMsg    { log.Error("Implement %d", p2p.GetReceiptsMsg)    ;return }
  if msg.Code == p2p.ReceiptsMsg       { log.Error("Implement %d", p2p.ReceiptsMsg)       ;return }

  // If we are here then we have unsupported message. 
  // Just print it for now.
  log.Error("Unknown msg code: %d\n", msg.Code)
}

func displayLogo(file string) error {
  data, err := os.ReadFile(file)
  if err != nil {
    return err
  }

  fmt.Printf(log.Magenta + "%s" + log.Reset, data)
  return nil
}