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

// Because of crappy ethereum transaction encoding/decoding,
// to detect Legact Transactions we must check if rlp encoding
// is a list. For other types we just check type.
func IsList(raw []byte) bool {
	// TODO: double check this condition. RLP list encoding
	// should begin with 0xc0-0xff prefix.
	if raw[0] <= 255 {
		return true
	}

	return false
}

// Types:
// LegacyTxType     = 0x00
// AccessListTxType = 0x01
// DynamicFeeTxType = 0x02
// BlobTxType       = 0x03

// Blob
// ----
// ChainID    *uint256.Int
// Nonce      uint64
// GasTipCap  *uint256.Int // a.k.a. maxPriorityFeePerGas
// GasFeeCap  *uint256.Int // a.k.a. maxFeePerGas
// Gas        uint64
// To         common.Address
// Value      *uint256.Int
// Data       []byte
// AccessList AccessList
// BlobFeeCap *uint256.Int // a.k.a. maxFeePerBlobGas
// BlobHashes []common.Hash
// Sidecar *BlobTxSidecar `rlp:"-"`
// V *uint256.Int `json:"v" gencodec:"required"`
// R *uint256.Int `json:"r" gencodec:"required"`
// S *uint256.Int `json:"s" gencodec:"required"`


// AccessList
// ----------
// ChainID    *big.Int        // destination chain ID
// Nonce      uint64          // nonce of sender account
// GasPrice   *big.Int        // wei per gas
// Gas        uint64          // gas limit
// To         *common.Address `rlp:"nil"` // nil means contract creation
// Value      *big.Int        // wei amount
// Data       []byte          // contract invocation input data
// AccessList AccessList      // EIP-2930 access list
// V, R, S    *big.Int        // signature values


// DynaicFee
// ---------
// ChainID    *big.Int
// Nonce      uint64
// GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
// GasFeeCap  *big.Int // a.k.a. maxFeePerGas
// Gas        uint64
// To         *common.Address `rlp:"nil"` // nil means contract creation
// Value      *big.Int
// Data       []byte
// AccessList AccessList
// V *big.Int `json:"v" gencodec:"required"`
// R *big.Int `json:"r" gencodec:"required"`
// S *big.Int `json:"s" gencodec:"required"`
