package main

import (
	"fmt"
	"os"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/node"
	"parasite/p2p"
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
  // go peer.StartReader(handler, handler, Dispatch)/

  node.ConnectNodes(nodes, srcPrv, handler, handler)

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

func displayLogo(file string) error {
  data, err := os.ReadFile(file)
  if err != nil {
    return err
  }

  fmt.Printf(log.Magenta + "%s" + log.Reset, data)
  return nil
}