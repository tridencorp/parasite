package fetchers

import (
	"parasite/p2p"
)

// Fetchers are responsible for fetching and validating data from the network.
//
// Since the blockchain is a decentralized and public network, we cannot be 100%
// sure that the data we receive is valid. We must cross-check it with other peers,
// and this will be the fetcher's job.
type Fetcher[Sender p2p.Sender] struct {
	PeerCount  uint8
	MatchCount uint8

	// Peers to which we will send requests. Require 'Send' function.
	Peers []Sender

	Response chan *p2p.Msg // Response from peer.
	Handler  chan *p2p.Msg // Handler to which we will send validated data.
}

func NewFetcher[T p2p.Sender](peerCount, matchCount uint8, peers []T) *Fetcher[T] {
  return &Fetcher[T] {
		PeerCount:  peerCount,
		MatchCount: matchCount,
		Peers:      peers,
		Response:   make(chan *p2p.Msg, 10),
		Handler:    make(chan *p2p.Msg, 10), 
	}
}

func (fetcher *Fetcher[T]) FetchBlockHeaders(numbers []uint64) []*p2p.Msg {
  // Send request to peers and set handler to this fetcher.
  for _, peer := range fetcher.Peers {
    msg, _ := p2p.GetBlockHeaders(14_678_700, 1)
    msg.Handler = fetcher.Response

    peer.Send(msg)
  }

  // Collect responses.
  responses, _ := fetcher.collectResponses()
  return responses
}

// Collecting responses.
func (fetcher *Fetcher[T]) collectResponses() ([]*p2p.Msg, error) {
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
