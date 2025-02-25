package tx

import (
	"fmt"
	"math/big"

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
		stream.Decode(legacy)
		tx.tx = legacy
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
		tx.tx = dynamic
		return nil
	}

	return fmt.Errorf("!!! Unknown Tx type !!!")
}
