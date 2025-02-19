package rpc

import (
	"encoding/hex"
	"fmt"
	"math/big"
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
func (node *Node) GetBalance(address string) (int64, error) {
	params := []any{address, "latest"}
	res := ""

	err := node.Send("eth_getBalance", params, &res)
	if err != nil {
		return -1, err
	}

	return strconv.ParseInt(res[2:], 16, 32)
}

// eth_gasPrice
// 
// Return the balance of the account of given address.
func (node *Node) GasPrice() (*big.Int, error) {
	res := ""

	err := node.Send("eth_gasPrice", nil, &res)
	if err != nil {
		return nil, err
	}

	price := new(big.Int)
	_, ok := price.SetString(res[2:], 16); 
	if !ok {
		return nil, fmt.Errorf("Error converting hex string to big.Int")
	}

	return price, nil
}

// eth_getCode
// 
// Return code at a given address.
func (node *Node) GetCode(address string) ([]byte, error) {
	params := []any{address, "latest"}
	res := ""

	err := node.Send("eth_getCode", params, &res)
	if err != nil {
		return nil, err
	}

	bytecode, err := hex.DecodeString(res[2:])
	return bytecode, nil
}
