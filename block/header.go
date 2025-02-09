package block

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type Bloom [256]byte
type BlockNonce [8]byte

type BlockHeader struct {
	ParentHash  common.Hash
	UncleHash   common.Hash
	Coinbase    common.Address
	Root        common.Hash
	TxHash      common.Hash
	ReceiptHash common.Hash
	Bloom       Bloom
	Difficulty  *big.Int
	Number      *big.Int
	GasLimit    uint64
	GasUsed     uint64
	Time        uint64
	Extra       []byte
	MixDigest   common.Hash
	Nonce       BlockNonce

	// optional
	BaseFee           *big.Int      `rlp:"optional"`
	WithdrawalsHash   *common.Hash  `rlp:"optional"`
	BlobGasUsed       *uint64       `rlp:"optional"`
	ExcessBlobGas     *uint64       `rlp:"optional"`
	ParentBeaconRoot  *common.Hash  `rlp:"optional"`
	RequestsHash      *common.Hash  `rlp:"optional"`
}

// Compute block header hash.
func (header *BlockHeader) Hash() (common.Hash, error) {
	b, err := rlp.EncodeToBytes(header)
	if err != nil {
		return common.Hash{}, err
	}

	return crypto.Keccak256Hash(b), nil
}