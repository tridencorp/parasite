package tx

import (
	"fmt"
	"math/big"
	"parasite/db"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type Transaction interface {
	nonce() 	 uint64
	gas() 		 uint64
	gasPrice() *big.Int
	to() 			*common.Address
	value() 	*big.Int
	data() 		 []byte
}

const (
	LegacyTxType     = 0x00
	AccessListTxType = 0x01
	DynamicFeeTxType = 0x02
	BlobTxType       = 0x03
)

// TODO: Tx is temporary workaround for stupid ethereum transaction design.
// Will be simplified.
type Tx struct {
	tx Transaction
	db db.Tx
}

func (tx *Tx) Nonce()    uint64          { return tx.tx.nonce() }
func (tx *Tx) Gas()      uint64          { return tx.tx.gas() }
func (tx *Tx) GasPrice() *big.Int        { return tx.tx.gasPrice() }
func (tx *Tx) To()       *common.Address { return tx.tx.to() }
func (tx *Tx) Value()    *big.Int        { return tx.tx.value() }
func (tx *Tx) Data()     []byte          { return tx.tx.data() }

func (tx *Tx) DecodeRLP(stream *rlp.Stream) error {
	kind, size, err := stream.Kind()
	if err != nil {
		return err 
	}
	
	if kind == rlp.Byte { 
		return fmt.Errorf("another stupid rlp error")
	}

	if kind == rlp.List {
		legacy := new(Legacy)

		err := stream.Decode(legacy)
		if err != nil {
			return err
		}

		dbTx := db.Tx{}
		dbTx.Type     = LegacyTxType
		dbTx.To       = legacy.To
		dbTx.Nonce    = legacy.Nonce
		dbTx.GasPrice = legacy.GasPrice
		dbTx.Gas      = legacy.Gas
		dbTx.Value    = legacy.Value
		dbTx.Data     = legacy.Data
		dbTx.V        = legacy.V
		dbTx.R        = legacy.R
		dbTx.S        = legacy.S
		
		tx.tx = legacy
		tx.db = dbTx
		return nil
	}

	buf := make([]byte, size)
	stream.ReadBytes(buf)

	// Not LegacyTx, check other types.
	if buf[0] == DynamicFeeTxType {
		// Remove type byte.
		buf = buf[1:]

		dynamic := new(DynamicFee)
		err := rlp.DecodeBytes(buf, dynamic)
		if err != nil {
			return nil
		}

		dbTx := db.Tx{}
		dbTx.Type      = DynamicFeeTxType
		dbTx.To        = dynamic.To
		dbTx.Nonce     = dynamic.Nonce
		dbTx.ChainID   = dynamic.ChainID
		dbTx.GasPrice  = dynamic.GasFeeCap
		dbTx.GasTipCap = dynamic.GasTipCap
		dbTx.Gas       = dynamic.Gas
		dbTx.Value     = dynamic.Value
		dbTx.Data      = dynamic.Data
		dbTx.V         = dynamic.V
		dbTx.R         = dynamic.R
		dbTx.S         = dynamic.S

		tx.tx = dynamic
		tx.db = dbTx
		return nil
	}

	return fmt.Errorf("!!! Unknown Tx type !!!")
}
