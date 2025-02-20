package rpc

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type Receipt struct {
	TransactionHash   common.Hash    `json:"transactionHash"`
	TransactionIndex  hexutil.Uint   `json:"transactionIndex"`
	BlockHash         common.Hash    `json:"blockHash"`
	BlockNumber       hexutil.Big    `json:"blockNumber"`
	From              common.Address `json:"from"`
	To                common.Address `json:"to"`
	CumulativeGasUsed hexutil.Uint64 `json:"cumulativeGasUsed"`
	EffectiveGasPrice hexutil.Big    `json:"effectiveGasPrice"`
	GasUsed           hexutil.Uint64 `json:"gasUsed"`
	ContractAddress   string         `json:"contractAddress"`
	Logs              []any          `json:"logs"`
	LogsBloom         string         `json:"logsBloom"`
	Type              string         `json:"type"`
	Root              common.Hash    `json:"root"`
	Status            hexutil.Uint64 `json:"status"`
}
