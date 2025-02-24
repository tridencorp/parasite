package tx

import (
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

func Decode(raw []byte) (*Tx, error) {
	tx := new(Tx)
	
	// Check if we have LegacyTx.
	if IsList(raw) {
		legacy := new(Legacy)
		err := rlp.DecodeBytes(raw, legacy)
		if err != nil {
			return nil, nil
		}
		tx.tx = legacy	
		return tx, nil
	}

	// No LegacyTx, other types.
	data := []byte{}
	rlp.DecodeBytes(raw, &data)
	
	if data[0] == 2 {
		// Remove type byte.
		data = data[1:] 
	
		dynamic := new(DynamicFee)
		err := rlp.DecodeBytes(data, dynamic)
		if err != nil {
			return nil, nil
		}

		tx.tx = dynamic
		return tx, nil
	}
	return nil, nil
}

// Because of crappy ethereum transaction encoding/decoding,
// to detect Legact Transactions we must check if rlp encoding
// is a list. For other types we just check type.
func IsList(raw []byte) bool {
	// TODO: double check this condition. RLP list encoding
	// should begin with 0xc0-0xff prefix.
	if raw[0] >= 192 && raw[0] <= 255 { return true }
	return false
}
