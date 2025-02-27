package fetchers

import (
	"fmt"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
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
type Fetcher[Sender p2p.Sender] struct {
	PeerCount  uint8
	MatchCount uint8

	// Peers to which we will send requests. Require 'Send' function.
	Peers []Sender

	Response chan *p2p.Msg // Response from peer.
	Handler  chan any      // Handler to which we will send validated data.
}

func NewFetcher[T p2p.Sender](peerCount, matchCount uint8, peers []T) *Fetcher[T] {
  return &Fetcher[T] {
		PeerCount:  peerCount,
		MatchCount: matchCount,
		Peers:      peers,
		Response:   make(chan *p2p.Msg, 10),
		Handler:    make(chan any, 10), 
	}
}

func (fetcher *Fetcher[T]) FetchBlockHeaders(numbers []uint64) {
  // Send request to peers and set handler to this fetcher.
  for _, peer := range fetcher.Peers {
    msg, _ := p2p.GetBlockHeaders(14_678_700, 1)
    msg.Handler = fetcher.Response

    peer.Send(msg)
  }

  // Collect messages.
  msgs, err := fetcher.collectMessages()
  if err != nil {
    fmt.Println(err)
  }

  // Validate response.
  headers, err := fetcher.validateHeaders(msgs)
  if err != nil {
    fmt.Println(err)
    // Data is not valid.
    // TODO: send error response to handler.
  }

  fetcher.Handler <- headers
}

// Collecting responses.
func (fetcher *Fetcher[T]) collectMessages() ([]*p2p.Msg, error) {
  msgs := []*p2p.Msg{}

  for {
    select {
    case msg := <-fetcher.Response:
      msgs = append(msgs, msg)
      if len(msgs) == int(fetcher.PeerCount) {
        return msgs, nil
      }
    }
  }
}

// We are verifying whether the header hash is the same across
// all responses received from peers.
func (fetcher *Fetcher[Sender]) validateHeaders(msgs []*p2p.Msg) ([]*p2p.BlockHeader, error) {
  expected := []*p2p.BlockHeader{}

  for i, msg := range msgs {
    headers := []*p2p.BlockHeader{}
    err := rlp.DecodeBytes(msg.Data, &headers)
    if err != nil {
      return nil, err
    }

    // Initial setup. First header will be our reference point.
    if i == 0 {
      expected = headers
      continue
    }

    // Go through all sub headers.
    for i, hh := range headers {
      if hh.Hash() != expected[i].Hash() {
        return nil, fmt.Errorf("Block headers do not match")
      }
    }

    expected = headers
  }

  return expected, nil
}
