package fetchers

import (
	"parasite/p2p"
	"parasite/test"
	"testing"
)

func TestFetchBlockHeaders(t *testing.T) {
	peers := []*test.Peer{}
	fetcher := NewFetcher(3, 3, peers)

	// Simulate response from peers.
	go func(res chan *p2p.Msg) {
		for i:=0; i < 3; i++ {
			res <- p2p.NewMsg(1, []byte("res from peer"))
		}
	}(fetcher.Response)

	res := fetcher.FetchBlockHeaders([]uint64{0})
	if len(res) < 3 {
		t.Errorf("Expected len to be %d, got %d", 3, len(res))
	}
}
