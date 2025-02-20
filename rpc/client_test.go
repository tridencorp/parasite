package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

const localAddress   = "http://127.0.0.1:8545"
const accountAddress = "0xdb2e0aaa2786bd19236aaccd9998452f72dd2b73"
const localBlock     = 14_000_000 
const localTx        = "0x3dac2080b4c423029fcc9c916bc430cde441badfe736fc6d1fe9325348af80fd"
const rawTx          = "0x3dac2080b4c423029fcc9c916bc430cde441badfe736fc6d1fe9325348af80fd"

func TestBlockNumber(t *testing.T) {
  node := NewNode(localAddress)
  res := hexutil.Uint64(0)
  err := node.BlockNumber(&res)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
  if res <= 0   { t.Errorf("Expected response to be > 0, got %d", res) }
}

func TestBalance(t *testing.T) {
  node := NewNode(localAddress)
  res := ""
  err := node.GetBalance(&res, accountAddress)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGasPrice(t *testing.T) {
  node := NewNode(localAddress)
  res := ""
  err := node.GasPrice(&res)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetCode(t *testing.T) {
  node := NewNode(localAddress)
  res := ""
  err := node.GetCode(&res, accountAddress)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetBlockByNumber(t *testing.T) {
  node := NewNode(localAddress)
  res := &Block{}
  err := node.GetBlockByNumber(res, localBlock)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetTransactionByHash(t *testing.T) {
  node := NewNode(localAddress)
  res := &Transaction{}
  err := node.GetTransactionByHash(res, localTx)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestGetTransactionReceipt(t *testing.T) {
  node := NewNode(localAddress)
  res := &Receipt{}
  err := node.GetTransactionReceipt(res, localTx)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}

func TestSendRawTransaction(t *testing.T) {
  node := NewNode(localAddress)
  res := ""
  err := node.SendRawTransaction(&res, rawTx)

  if err != nil { t.Errorf("Expected no errors, got %s", err) }
}
