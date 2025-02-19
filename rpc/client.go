package rpc

import (
	"strconv"
)

// eth_blockNumber
//
// Get the number of most recent block.
func (node *Node) BlockNumber() (int64, error) {
	res := ""

	err := node.Send("eth_blockNumber", nil, &res)
	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(res[2:], 16, 32)
}

// eth_getBalance
// 
// Return the balance of the account of given address.
func (node *Node) Balance(address string) (int64, error) {
	params := []any{address, "latest"}
	res := ""

	err := node.Send("eth_getBalance", params, &res)
	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(res[2:], 16, 32)
}

