package fetchers

import (
	"fmt"
	"parasite/p2p"
	"time"
)

// Fetchers are responsible for fetching and validating data from peers.
//
// Since the blockchain is a decentralized and public network, we cannot be 100%
// sure that the data we receive is valid. We must cross-check it with other peers,
// and this will be the fetcher's job.
//
// Data flow between peer, fetcher and handler.
//
//    +--------Send()--------+
//    V                      |
// +------+             +---------+ <---Input---- +---------+
// | PEER | ---Input--> | FETCHER |               | HANDLER |
// +------+             +---------+ ---Output---> +---------+

type Fetcher[T any] struct {
  Input  chan *p2p.Msg // Input channel will primarily receive messages from peers.
  Output chan T        // Output channel for sending responses to handlers.

  Peers []p2p.Sender

  // Callbacks.
  Validate func(msgs []*p2p.Msg) (T, error)
  Request  func() *p2p.Msg
}

// Continuously fetch data from peers. Runs until channel is closed.
// Interval is specified in milliseconds.
func (fetcher *Fetcher[T]) Run(interval int) {
  ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	defer ticker.Stop()

  // Send request to peers every given number of milliseconds.
  for range ticker.C {
    // Send message to peers.
    fetcher.Send()

    // Wait for response from peers and collect the messages.
    msgs, err := fetcher.Collect()
    if err != nil {
      fmt.Println(err)
      return
    }

    // Validate response.
    // TODO: Handle failures.
    headers, err := fetcher.Validate(msgs)
    if err != nil {
      fmt.Println(err)
    }

    // Send response to output channel.
    // In most cases this will be one of the handlers.
    fetcher.Output <- headers
  }
}

// Prepare and send message to peers.
func (fetcher *Fetcher[T]) Send() {
  // Call Request callback.
  req := fetcher.Request()

  // Send message to peers.
  for _, peer := range fetcher.Peers {
    peer.Send(req)
  }

}

// Collect response from peers.
// TODO: add timeout and terminate if no response will come.
func (fetcher *Fetcher[T]) Collect() ([]*p2p.Msg, error) {
  msgs := []*p2p.Msg{}

  for {
    select {
    case msg, ok := <-fetcher.Input:
      msgs = append(msgs, msg)

      // Channel was closed.
      if !ok {
        return nil, fmt.Errorf("chanel was closed")
      }

      // Number of received messages must be equal to the number of peers.
      if len(msgs) == len(fetcher.Peers) {
        return msgs, nil
      }
    }
  }
}