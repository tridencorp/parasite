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
// Data flow (channels):
//
// +------+                +---------+                 +---------+
// | PEER | ---Response--> | FETCHER | <---Handler---> | HANDLER |
// +------+                +---------+                 +---------+
//
type Fetcher[T p2p.Sender] struct {
	PeerCount  uint8
	MatchCount uint8

	// Peers to which we will send requests. Require 'Send' function.
	Peers []T

	PeerRes chan *p2p.Msg // Response from peer.

	HandlerReq  chan *p2p.Msg // Request from handler.
	HandlerRes  chan *p2p.Msg // Response for handler.

  // Custom fetchers must implement these functions.
  Message  func(params any) *p2p.Msg
  Validate func(msgs []*p2p.Msg) (*p2p.Msg, error)
}

func NewFetcher[T p2p.Sender](peerCount, matchCount uint8, peers []T) *Fetcher[T] {
  return &Fetcher[T] {
		PeerCount:  peerCount,
		MatchCount: matchCount,
		Peers:      peers,
		PeerRes:    make(chan *p2p.Msg, 10),
		HandlerRes: make(chan *p2p.Msg, 10), 
		HandlerReq: make(chan *p2p.Msg, 10), 
	}
}

func (fetcher *Fetcher[T]) FetchBlockHeaders(number uint64) {
  // Prepare message using Message callback.
  params := uint64(0)
  msg := fetcher.Message(params)
  msg.Handler = fetcher.PeerRes

  // Send message to peers.
  for _, peer := range fetcher.Peers {
    peer.Send(msg)
  }
      
  // Collect messages.
  msgs, err := fetcher.collectMessages()
  if err != nil {
    fmt.Println(err)
  }

  // Validate response.
  headers, err := fetcher.Validate(msgs)
  if err != nil {
    fmt.Println(err)
  }

  fetcher.HandlerRes <- &p2p.Msg{Payload: headers}
}

// Collecting responses.
func (fetcher *Fetcher[T]) collectMessages() ([]*p2p.Msg, error) {
  msgs := []*p2p.Msg{}

  for {
    select {
    case msg := <-fetcher.PeerRes:
      msgs = append(msgs, msg)
      if len(msgs) == int(fetcher.PeerCount) {
        return msgs, nil
      }
    }
  }
}