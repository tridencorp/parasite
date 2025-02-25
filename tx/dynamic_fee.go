package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type AccessList []AccessTuple

type AccessTuple struct {
	Address     common.Address
	StorageKeys []common.Hash 
}

type DynamicFee struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int
	GasFeeCap  *big.Int
	Gas        uint64
	To         *common.Address `rlp:"nil"`
	Value      *big.Int
	Data       []byte
	AccessList AccessList

	V *big.Int
	R *big.Int
	S *big.Int
}

func (tx *DynamicFee) nonce()    uint64          { return tx.Nonce }
func (tx *DynamicFee) gas()      uint64          { return tx.Gas }
func (tx *DynamicFee) gasPrice() *big.Int        { return tx.GasFeeCap }
func (tx *DynamicFee) to()       *common.Address { return tx.To }
func (tx *DynamicFee) value()    *big.Int        { return tx.Value }
func (tx *DynamicFee) data()     []byte          { return tx.Data }
