package rpc

import "github.com/ethereum/go-ethereum/common"

type Transaction struct {
  BlockHash   common.Hash    `json:"blockHash"`
	BlockNumber string         `json:"blockNumber"`
	Value       string         `json:"value"`
	From        common.Address `json:"from"`
	To          common.Address `json:"to"`
	Gas         string         `json:"gas"`
	GasPrice    string         `json:"gasPrice"`
	Hash        common.Hash    `json:"hash"`
	Input       string         `json:"input"`
	Nonce       string         `json:"nonce"`
	Index       string         `json:"transactionIndex"`
	V           string         `json:"v"`
	R           string         `json:"r"`
	S           string         `json:"s"`
}