package rpc

import (
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	ParentHash  common.Hash    `json:"parentHash"`
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"`
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       string         `json:"logsBloom"`
	Difficulty  string         `json:"difficulty"`
	Number      string         `json:"number"`
	GasLimit    string         `json:"gasLimit"`
	GasUsed     string         `json:"gasUsed"`
	Time        string         `json:"timestamp"`
	Extra       string         `json:"extraData"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       string 				 `json:"nonce"`

	Transactions []any         `json:"transactions"`
	Uncles 			 []any         `json:"uncles"`
}
