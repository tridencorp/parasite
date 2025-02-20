package rpc

import (
	"fmt"
)

// Get the number of most recent block.
func (node *Node) BlockNumber(res any) error {
	return node.Send("eth_blockNumber", res, nil)
}

// Return the balance of the account of given address.
func (node *Node) GetBalance(res any, address string) error {
	return node.Send("eth_getBalance", res, []any{address, "latest"})
}

// Return current gas price.
func (node *Node) GasPrice(res any) error {
  return node.Send("eth_gasPrice", res, nil)
}

// Return code at a given address.
func (node *Node) GetCode(res any, address string) error {
  return  node.Send("eth_getCode", res, []any{address, "latest"})
}

// Return block by given number.
func (node *Node) GetBlockByNumber(res any, number uint32) error {
  hex := fmt.Sprintf("0x%x", number)
	return node.Send("eth_getBlockByNumber", res, []any{hex, true})
}

// Return transaction by given hash.
func (node *Node) GetTransactionByHash(res any, hash string) error {
	return node.Send("eth_getTransactionByHash", res, []any{hash})
}

// Return receipt for given transaction hash.
func (node *Node) GetTransactionReceipt(res any, hash string) error {
  return node.Send("eth_getTransactionReceipt", res, []any{hash})
}

// Send signed transaction (or contract) to ethereum network.
func (node *Node) SendRawTransaction(res any, hex string) error {
  return node.Send("eth_sendRawTransaction", res, []any{hex})
}
