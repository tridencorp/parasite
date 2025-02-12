package main

import (
	"fmt"
	"os"
	"parasite/config"
	"parasite/key"
	"parasite/log"
	"parasite/p2p"
	"parasite/server"
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
	handler := make(chan p2p.Msg) 
	failure := make(chan p2p.Msg, 10) 

	dispatcher := server.NewDispatcher(nil, handler, failure)

	peerHub := p2p.NewPeerHub(nodes, dispatcher, failure, srcPrv)
	peerHub.ConnectAll()
	go peerHub.Start()

  for _ = range handler {
    fmt.Println("!!! got headers via handler !!!")
    // fmt.Printf("msg: %v", msg)
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