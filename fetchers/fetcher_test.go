package fetchers

import (
	"fmt"
	"parasite/p2p"
	"parasite/test"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

func TestFetchBlockHeaders(t *testing.T) {
	// Test with valid response from all peers.
	peers := []*test.Peer{}
	fetcher := HeaderFetcher(3, 3, peers)

	headers := test.Headers(1)
	raw, _ := rlp.EncodeToBytes(headers)

	// Simulate valid response from peers.
	go func(res chan *p2p.Msg) {
		for i:=0; i < 3; i++ {
			res <- p2p.NewMsg(1, raw)
		}
	}(fetcher.PeerRes)
	
	fetcher.FetchBlockHeaders(uint64(0))
	msg := <- fetcher.HandlerRes
	res := msg.Payload.(*p2p.Msg)
	
	fmt.Println(res)
}
