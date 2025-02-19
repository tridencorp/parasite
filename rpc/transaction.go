package rpc

import "github.com/ethereum/go-ethereum/common"

type Transaction struct {
  BlockHash   common.Hash    `json:"blockHash"`
	BlockNumber string         `json:"blockNumber"`
	From        common.Address `json:"from"`
	Gas         string         `json:"gas"`
	GasPrice    string         `json:"gasPrice"`
	Hash        common.Hash    `json:"hash"`
	Input       string         `json:"input"`
	Nonce       string         `json:"nonce"`
	To          common.Address `json:"to"`
	Index       string         `json:"transactionIndex"`
	Value       string         `json:"value"`
	V           string         `json:"v"`
	R           string         `json:"r"`
	S           string         `json:"s"`
}