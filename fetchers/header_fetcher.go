package fetchers

import (
	"fmt"
	"parasite/p2p"
)

type HeaderFetcher struct {
	Fetcher[[]*p2p.BlockHeader]
}

func NewHeaderFetcher(in chan *p2p.Msg, out chan []*p2p.BlockHeader) *HeaderFetcher {
	fetcher := new(HeaderFetcher)

	fetcher.Input    = make(chan *p2p.Msg, 10)
	fetcher.Output   = out
	fetcher.Validate = fetcher.validate
	fetcher.Request  = fetcher.request

	return fetcher
}

func (fetcher *HeaderFetcher) validate(msgs []*p2p.Msg) ([]*p2p.BlockHeader, error) {
  expected := []*p2p.BlockHeader{}

  for i, msg := range msgs {
		headers, _ := p2p.DecodeBlockHeaders(msg)
		
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

  return expected, nil	
}

func (fetcher *HeaderFetcher) request(args ...any) *p2p.Msg {
	number := args[0].(uint64)
	amount := args[1].(uint64)

	msg, _ := p2p.GetBlockHeaders(number, amount)
	msg.Handler = fetcher.Input

	return msg
}
