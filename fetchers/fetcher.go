package fetchers

import (
	"fmt"
	"parasite/p2p"
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
	Input  chan *p2p.Msg // Response from peer, requests from handler.
  Output chan T        // Response for handler.

  Peers []p2p.Sender

  // Callbacks.
  Validate func(msgs []*p2p.Msg) (T, error)
  Request  func(args ...any) *p2p.Msg
}

// Fetch data from peers. It do it one time and then terminates.
func (fetcher *Fetcher[T]) Fetch(args ...any) {
  // Prepare message using Request callback.
  req := fetcher.Request(args...)

  // Send message to peers.
  for _, peer := range fetcher.Peers {
    peer.Send(req)
  }
  // Wait for response from all peers and collect messages.
  msgs, err := fetcher.Collect()
  if err != nil {
    fmt.Println(err)
  }
 
  // Validate response.
  // TODO: Handle failures.
  headers, err := fetcher.Validate(msgs)
  if err != nil {
    fmt.Println(err)
  }

  // Send response to output channel. In most cases this will be one of the handler.
  fetcher.Output <- headers
}

// Collecting responses.
func (fetcher *Fetcher[T]) Collect() ([]*p2p.Msg, error) {
  msgs := []*p2p.Msg{}

  for {
    select {
    case msg := <-fetcher.Input:
      msgs = append(msgs, msg)

      // Number of received messages must be equal to number of peers.
      if len(msgs) == len(fetcher.Peers) {
        return msgs, nil
      }
    }
  }
}