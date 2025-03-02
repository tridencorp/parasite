package fetchers

import (
	"bytes"
	"parasite/test"
	"testing"
)

func TestFetchBlockHeaders(t *testing.T) {
	// Test with valid response from all peers.
	msg, headers := test.HeadersResponse(1, 1)
	fetcher := NewHeaderFetcher()

	// Add peers to fetcher.
	// TODO: Figure out better way.
	peers := test.Peers(3, msg, fetcher.Input)
	for _, peer := range peers { fetcher.Peers = append(fetcher.Peers, peer) }

	go fetcher.Fetch(uint64(0), uint64(1))
	res := <- fetcher.Output

	expected := headers[0].Hash().Bytes()
	got      := res[0].Hash().Bytes()

	if !bytes.Equal(expected, got) {
		t.Errorf("Expected headers\nto be %v,\n  got %v", expected, got)
	}
}
