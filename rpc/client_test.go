package rpc

import (
	"testing"
)

const localAddress   = "http://127.0.0.1:8545"
const accountAddress = "0xdb2e0aaa2786bd19236aaccd9998452f72dd2b73"
const localBlock     = 14_000_000 
const localTx        = "0x3dac2080b4c423029fcc9c916bc430cde441badfe736fc6d1fe9325348af80fd"

func TestBlockNumber(t *testing.T) {
	node := NewNode(localAddress)
	res, err := node.BlockNumber()

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
	if res <= 0   { t.Errorf("Expected response to be > 0, got %d", res) }
}

func TestBalance(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GetBalance(accountAddress)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGasPrice(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GasPrice()

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetCode(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GetCode(accountAddress)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetBlockByNumber(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GetBlockByNumber(localBlock)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetTransactionByHash(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GetTransactionByHash(localTx)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetTransactionReceipt(t *testing.T) {
	node := NewNode(localAddress)
	_, err := node.GetTransactionReceipt(localTx)

	if err != nil { t.Errorf("Expected no errors, got %s", err) }
}
