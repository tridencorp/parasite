package p2p

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// All ETH68 message codes.
// We will be only supporting the newest ETH protocol.

const ETH = 68

const BaseProtocolSize = 16

const (
	// Base protocol msg codes: 0x00...0x10 (0-16)
	// Only 4 are used for now.
	HandshakeMsg = 0x00
	DiscMsg      = 0x01
	PingMsg      = 0x02
	PongMsg      = 0x03

	// Extended protocol msg codes: 0x10...0x1C (16-28)
	StatusMsg                     = 0x10
	NewBlockHashesMsg             = 0x11 // Not supported anymore
	TransactionsMsg               = 0x12
	GetBlockHeadersMsg            = 0x13
	BlockHeadersMsg               = 0x14
	GetBlockBodiesMsg             = 0x15
	BlockBodiesMsg                = 0x16
	NewBlockMsg                   = 0x17 // Not supported anymore
	NewPooledTranasctionHashesMsg = 0x18
	GetPooledTransactionsMsg      = 0x19
	PooledTransactionsMsg         = 0x1A
	GetReceiptsMsg                = 0x1B
	ReceiptsMsg                   = 0x1C
)

type Capability struct {
	Name    string
	Version uint
}

type Handshake struct {
	Version    uint64
	Name       string
	Caps       []Capability
	ListenPort uint64
	ID         []byte

	// Currently unused, but required for compatibility with ETH.
	Rest []rlp.RawValue `rlp:"tail"`
}

// type GetBlocksPacket struct {
// 	RequestId uint64
// 	Headers []common.Hash
// }

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

type BlockHeaders struct {
	RequestId uint64
	Headers []*BlockHeader
}