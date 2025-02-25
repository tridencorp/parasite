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
	ParentHash  common.Hash    `json:"parentHash"`
	UncleHash   common.Hash    `json:"sha3Uncles"`
	Coinbase    common.Address `json:"miner"`
	Root        common.Hash    `json:"stateRoot"`
	TxHash      common.Hash    `json:"transactionsRoot"`
	ReceiptHash common.Hash    `json:"receiptsRoot"`
	Bloom       Bloom          `json:"logsBloom"`
	Difficulty  *big.Int       `json:"difficulty"`
	Number      *big.Int       `json:"number"`
	GasLimit    uint64         `json:"gasLimit"`
	GasUsed     uint64         `json:"gasUsed"`
	Time        uint64         `json:"timestamp"`
	Extra       []byte         `json:"extraData"`
	MixDigest   common.Hash    `json:"mixHash"`
	Nonce       BlockNonce     `json:"nonce"`

	BaseFee          *big.Int     `json:"baseFeePerGas"         rlp:"optional"`
	WithdrawalsHash  *common.Hash `json:"withdrawalsRoot"       rlp:"optional"`
	BlobGasUsed      *uint64      `json:"blobGasUsed"           rlp:"optional"`
	ExcessBlobGas    *uint64      `json:"excessBlobGas"         rlp:"optional"`
	ParentBeaconRoot *common.Hash `json:"parentBeaconBlockRoot" rlp:"optional"`
	RequestsHash     *common.Hash `json:"requestsHash"          rlp:"optional"`
}

// Compute block header hash.
func (header *BlockHeader) Hash() (common.Hash, error) {
	b, err := rlp.EncodeToBytes(header)
	if err != nil {
		return common.Hash{}, err
	}

	return crypto.Keccak256Hash(b), nil
}

type Block struct {
	BlockHeader
	Transactions []any `json:"transactions"`
	Uncles       []any `json:"uncles"`
}
