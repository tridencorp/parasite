package rpc

import "github.com/ethereum/go-ethereum/common"

type Receipt struct {
	TransactionHash   common.Hash    `json:"transactionHash"`
	TransactionIndex  string         `json:"transactionIndex"`
	BlockHash         common.Hash    `json:"blockHash"`
	BlockNumber       string         `json:"blockNumber"`
	From              common.Address `json:"from"`
	To                common.Address `json:"to"`
	CumulativeGasUsed string         `json:"cumulativeGasUsed"`
	EffectiveGasPrice string         `json:" effectiveGasPrice"`
	GasUsed           string         `json:"gasUsed"`
	ContractAddress   common.Address `json:"contractAddress"`
	Logs              []any          `json:"logs"`
	LogsBloom         string         `json:"logsBloom"`
	Type              string         `json:"type"`
	Root              common.Hash    `json:"root"`
	Status            string         `json:"status"`
}