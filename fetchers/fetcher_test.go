package fetchers

import (
	"fmt"
	"parasite/test"
	"testing"
)


func TestFetchBlockHeaders(t *testing.T) {
	// Test with valid response from all peers.
	peers := test.Peers(3)
	fmt.Println(peers)
}
