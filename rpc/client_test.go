package rpc

import (
	"testing"
)

const localAddress   = "http://127.0.0.1:8545"
const accountAddress = "0x407d73d8a49eeb85d32cf465507dd71d507100c1"

func TestBlockNumber(t *testing.T) {
	node := NewNode(localAddress)
	res, err := node.BlockNumber()

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
	if res <= 0   { t.Errorf("Expected response to be > 0, got %d", res) }
}

func TestBalance(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.Balance(accountAddress)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}
