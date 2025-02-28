package fetchers

import (
	"fmt"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
)

func HeaderFetcher[T p2p.Sender](peerCount, matchCount uint8, peers []T) *Fetcher[T] {
	fetcher := NewFetcher(2, 2, peers)

	fetcher.Message  = HeaderMessage
	fetcher.Validate = ValidateHeaders

	return fetcher
}

func HeaderMessage(params any) *p2p.Msg {
	number := params.(uint64)
	msg, _ := p2p.GetBlockHeaders(number, 1)
	return msg
}

func ValidateHeaders(msgs []*p2p.Msg) (*p2p.Msg, error) {
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
				return nil, fmt.Errorf("Block headers hash do not match") 
			}
		}

    expected = headers
  }

  return msgs[0], nil	
}
