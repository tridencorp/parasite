package tx

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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
