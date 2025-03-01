package fetchers

import (
	"bytes"
	"parasite/p2p"
	"parasite/test"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)


func TestFetchBlockHeaders(t *testing.T) {
	// Test with valid response from all peers.
	headers := test.Headers(1)
	data    := []any{uint64(1), headers}
	req, _  := rlp.EncodeToBytes(data)

	in  := make(chan *p2p.Msg, 10)
	out := make(chan []*p2p.BlockHeader, 10)

	msg := p2p.NewMsg(1, req)
	msg.Handler = in

	f := NewHeaderFetcher(in, out)
	
	// Add peers to fetcher.
	// TODO: Figure out better way.
	peers := test.Peers(3, msg)
	for _, peer := range peers { f.Peers = append(f.Peers, peer) }

	go f.Fetch(uint64(0), uint64(1))
	res := <- f.Output

	expected := headers[0].Hash().Bytes()
	got      := res[0].Hash().Bytes()

	if !bytes.Equal(expected, got) {
		t.Errorf("Expected headers\nto be %v,\n  got %v", expected, got)
	}
}
