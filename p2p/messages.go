package p2p

import (
	"parasite/block"
	"parasite/tx"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	StatusMsg                     = 0x10
	NewBlockHashesMsg             = 0x11 // Not supported anymore
	TransactionsMsg               = 0x12
	GetBlockHeadersMsg            = 0x13
	BlockHeadersMsg               = 0x14
	GetBlockBodiesMsg             = 0x15
	BlockBodiesMsg                = 0x16
	NewBlockMsg                   = 0x17 // Not supported anymore
	NewPooledTransactionHashesMsg = 0x18
	GetPooledTransactionsMsg      = 0x19
	PooledTransactionsMsg         = 0x1A
	GetReceiptsMsg                = 0x1F
	ReceiptsMsg                   = 0x1A
)

var DiscReasons = []string{
	"disconnect requested",                // 0x00
	"network error",                       // 0x01
	"breach of protocol",                  // 0x02
	"useless peer",                        // 0x03
	"too many peers",                      // 0x04
	"already connected",                   // ...
	"incompatible p2p protocol version",
	"invalid node identity",
	"client quitting",
	"unexpected identity",
	"connected to self",
	"read timeout",
	"subprotocol error",
	"invalid disconnect reason",
}

type Request struct {
	ReqID uint64
	Data  any
}

type getBlockHeadersMsg struct {
	Number  uint64
	Amount  uint64
	Skip    uint64
	Reverse bool
}

type blockHeadersMsg struct {
	ReqID uint64
	Headers []*block.BlockHeader
}

type blockBodiesMsg struct {
	ReqID uint64
	Bodies []*BlockBody
}

type receiptsMsg struct {
	ReqID uint64
	Receipts [][]*tx.Receipt
}

type pooledTransactions struct {
	Types  []byte
	Sizes  []uint32
	Hashes []common.Hash
}

type BlockBody struct {
	Transactions []*tx.Tx
	Uncles       []*block.BlockHeader
	Withdrawals  []*types.Withdrawal `rlp:"optional"`
}

// Decode BlockHeadersMsg response.
func DecodeBlockHeaders(msg *Msg) ([]*block.BlockHeader, error) {
	res := new(blockHeadersMsg)
	return res.Headers, rlp.DecodeBytes(msg.Data, &res)
}

// Decode BlockBodiesMsg.
func DecodeBlockBodies(msg *Msg) ([]*BlockBody, error) {
	res := blockBodiesMsg{}
	return res.Bodies, rlp.DecodeBytes(msg.Data, &res)
}

// Decode ReceiptsMsg.
func DecodeReceipts(msg *Msg) ([][]*tx.Receipt, error) {
	res := receiptsMsg{}
	return res.Receipts, rlp.DecodeBytes(msg.Data, &res)
}
