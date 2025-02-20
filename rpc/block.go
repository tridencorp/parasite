package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Block struct {
	ParentHash  common.Hash    `json:"parentHash"`
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"`
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       string         `json:"logsBloom"`
	Difficulty  hexutil.Big    `json:"difficulty"`
	Number      hexutil.Big    `json:"number"`
	GasLimit    hexutil.Uint64 `json:"gasLimit"`
	GasUsed     hexutil.Uint64 `json:"gasUsed"`
	Time        hexutil.Uint64 `json:"timestamp"`
	Extra       hexutil.Bytes  `json:"extraData"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       hexutil.Uint64 `json:"nonce"`

	Transactions []Transaction `json:"transactions"`
	Uncles       []any         `json:"uncles"`
}
