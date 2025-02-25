package db

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// This is our main transaction type. It will be used
// in all internal operations.
// All ethereum transactions will be mapped to it.
type Tx struct {
	Type              uint8
	From              *common.Address
	To                *common.Address
	Value             *big.Int
	Nonce             uint64
	Hash              common.Hash
	ChainID           *big.Int
	Status            uint64
	BlockNumber       *big.Int

	GasUsed           uint64
	GasPrice          *big.Int
	CumulativeGasUsed uint64
	Gas               uint64
	GasTipCap         *big.Int
	V, R, S           *big.Int

	Data []byte
}
