package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const BloomByteLength = 256

type Bloom [BloomByteLength]byte

type Receipt struct {
	Type              uint8  `json:"type,omitempty"`
	PostState         []byte `json:"root"`
	Status            uint64 `json:"status"`
	CumulativeGasUsed uint64 `json:"cumulativeGasUsed"`
	Bloom             Bloom  `json:"logsBloom"`
	Logs              []byte `json:"logs"`

	TxHash            common.Hash    `json:"transactionHash"`
	ContractAddress   common.Address `json:"contractAddress"`
	GasUsed           uint64         `json:"gasUsed"`
	EffectiveGasPrice *big.Int       `json:"effectiveGasPrice"`
	BlobGasUsed       uint64         `json:"blobGasUsed,omitempty"`
	BlobGasPrice      *big.Int       `json:"blobGasPrice,omitempty"`
	BlockHash         common.Hash    `json:"blockHash,omitempty"`
	BlockNumber       *big.Int       `json:"blockNumber,omitempty"`
	TransactionIndex  uint           `json:"transactionIndex"`
}
