package tx

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const BloomByteLength = 256

type Bloom [BloomByteLength]byte

type Receipt struct {
	Type              uint8 
	PostState         []byte
	Status            uint64
	CumulativeGasUsed uint64
	Bloom             Bloom 
	Logs              []byte

	TxHash            common.Hash
	ContractAddress   common.Address
	GasUsed           uint64
	EffectiveGasPrice *big.Int
	BlobGasUsed       uint64
	BlobGasPrice      *big.Int
	BlockHash         common.Hash
	BlockNumber       *big.Int   
	TransactionIndex  uint      
}

type Withdrawal struct {
	Index     uint64        
	Validator uint64        
	Address   common.Address
	Amount    uint64        
}

func (r *Receipt) DecodeRLP(stream *rlp.Stream) error {
	kind, size, err := stream.Kind()
	if err != nil {
		return err
	}

	if kind == rlp.Byte {
		return fmt.Errorf("another stupid RLP error")
	}

	if kind == rlp.List {
		data := []any{[]byte{}, []byte{}, []byte{}, []byte{}}
		err := stream.Decode(&data)
		if err != nil {
			return err
		}		
		
		gas := new(big.Int).SetBytes(data[1].([]byte))
		r.CumulativeGasUsed = gas.Uint64()
		return nil
	}

	buf := make([]byte, size)
	err = stream.ReadBytes(buf)
	if err != nil {
		return err
	}

	if buf[0] > 0x00 && buf[0] <= 0x03 {
		if len(buf) <= 1 {
			return fmt.Errorf("Decoding receipts rlp error")
		}

		data := []any{[]byte{}, uint64(0), []byte{}, []byte{}}
		err = rlp.DecodeBytes(buf[1:], &data)
		if err != nil {
			return err
		}

		gas := new(big.Int).SetBytes(data[1].([]byte))
		r.CumulativeGasUsed = gas.Uint64()
		return nil
	}

	return fmt.Errorf("Decoding receipts rlp error")
}