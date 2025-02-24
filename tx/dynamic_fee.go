package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type AccessList []AccessTuple

type AccessTuple struct {
	Address     common.Address `json:"address"`
	StorageKeys []common.Hash  `json:"storageKeys"`
}

type DynamicFee struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList AccessList

	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`
}

func (tx *DynamicFee) nonce()    uint64          { return tx.Nonce }
func (tx *DynamicFee) gas()      uint64          { return tx.Gas }
func (tx *DynamicFee) gasPrice() *big.Int        { return tx.GasFeeCap }
func (tx *DynamicFee) to()       *common.Address { return tx.To }
func (tx *DynamicFee) value()    *big.Int        { return tx.Value }
func (tx *DynamicFee) data()     []byte          { return tx.Data }
