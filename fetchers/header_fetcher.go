package fetchers

import (
	"fmt"
	"parasite/p2p"
)

type HeaderFetcher struct {
	Fetcher[[]*p2p.BlockHeader]

	Number uint64 
	Amount uint64
}

func NewHeaderFetcher() *HeaderFetcher {
	fetcher := new(HeaderFetcher)

	fetcher.Input    = make(chan *p2p.Msg, 10)
	fetcher.Output   = make(chan []*p2p.BlockHeader, 10)
	fetcher.Validate = fetcher.validate
	fetcher.Request  = fetcher.request

	return fetcher
}

func (fetcher *HeaderFetcher) validate(responses []*p2p.Msg) ([]*p2p.BlockHeader, error) {
  expected := []*p2p.BlockHeader{}

  for i, res := range responses {
		headers, err := p2p.DecodeBlockHeaders(res)
		if err != nil {
			fmt.Println(err)
		}

    // Initial setup. First header will be our reference point.
    if i == 0 {
      expected = headers
      continue
    }

    // Go through all headers.
    for i, header := range headers {
			if header.Hash() != expected[i].Hash() {
				return nil, fmt.Errorf("Block headers hash do not match") 
			}
		}
  }

  return expected, nil	
}

// This request method will be called each iteration, based
// on given time interval.
func (fetcher *HeaderFetcher) request(args ...any) *p2p.Msg {
	number := args[0].(uint64)
	amount := args[1].(uint64)

	msg, _ := p2p.GetBlockHeaders(number, amount)
	msg.Handler = fetcher.Input

	return msg
}
