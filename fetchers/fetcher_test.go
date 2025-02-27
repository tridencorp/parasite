package fetchers

import (
	"parasite/p2p"
	"parasite/test"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestFetchBlockHeaders(t *testing.T) {
	// Test with valid response from all peers.
	peers := []*test.Peer{}
	fetcher := NewFetcher(3, 3, peers)

	headers := test.Headers(1)
	bytes, _ := rlp.EncodeToBytes(headers)

	// Simulate response from peers.
	go func(res chan *p2p.Msg) {
		for i:=0; i < 3; i++ {
			res <- p2p.NewMsg(1, bytes)
		}
	}(fetcher.Response)

	fetcher.FetchBlockHeaders([]uint64{0})
	msg := <- fetcher.Handler
	res := msg.([]*p2p.BlockHeader)

	if len(res) != 1 {
		t.Errorf("Expected len to be %d, got %d", 1, len(res))
	}
}
