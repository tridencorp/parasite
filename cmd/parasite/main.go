package main

import (
	"fmt"
	"os"
	"parasite/config"
	"parasite/log"
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
