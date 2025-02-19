package rpc

import (
	"testing"
)

const localAddress = "http://127.0.0.1:8545"

func TestBlockNumber(t *testing.T) {
	node := NewNode(localAddress)
	res, _ := node.BlockNumber()

	if res <= 0 {
		t.Errorf("Expected response to be > 0, got %d", res)
	}
}
