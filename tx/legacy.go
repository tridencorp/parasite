package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Legacy struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *common.Address `rlp:"nil"`
	Value    *big.Int
	Data     []byte
	V, R, S  *big.Int
}

func (tx *Legacy) nonce()    uint64          { return tx.Nonce }
func (tx *Legacy) gas()      uint64          { return tx.Gas }
func (tx *Legacy) gasPrice() *big.Int        { return tx.GasPrice }
func (tx *Legacy) to()       *common.Address { return tx.To }
func (tx *Legacy) value()    *big.Int        { return tx.Value }
func (tx *Legacy) data()     []byte          { return tx.Data }
