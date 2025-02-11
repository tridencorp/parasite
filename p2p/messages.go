package p2p

import (
	"math/rand/v2"
	"parasite/block"

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

	// Extended protocol msg codes: 0x10...0x1C (16-28)
	StatusMsg                     = 0x10

	// Not supported anymore
	NewBlockHashesMsg             = 0x11

	TransactionsMsg               = 0x12
	GetBlockHeadersMsg            = 0x13
	BlockHeadersMsg               = 0x14
	GetBlockBodiesMsg             = 0x15
	BlockBodiesMsg                = 0x16

	// Not supported anymore
	NewBlockMsg                   = 0x17

	NewPooledTransactionHashesMsg = 0x18
	GetPooledTransactionsMsg      = 0x19
	PooledTransactionsMsg         = 0x1A
	GetReceiptsMsg                = 0x1B
	ReceiptsMsg                   = 0x1C
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

// Block headers

type blockHeadersReq struct {
	ReqId 	uint64
	Request struct {
		Number  uint64
		Amount  uint64
		Skip    uint64
		Reverse bool
	}
}

type blockHeadersRes struct {
	ReqId   uint64
	Headers []*block.BlockHeader
}

// Create BlockHeaders request message.
func BlockHeadersReq(number, amount, skip uint64, reverse bool) (Msg, error) {
	req := blockHeadersReq{
		ReqId: rand.Uint64(),
		Request: struct {
			Number  uint64
			Amount  uint64
			Skip    uint64
			Reverse bool
		}{
			Number:  number,
			Amount:  amount,
			Skip:    skip,
			Reverse: false,
		},
	}

	data, err := rlp.EncodeToBytes(req)
	if err != nil {
		return Msg{}, err
	}

	msg := NewMsg(GetBlockHeadersMsg, data)
	msg.ReqId = req.ReqId

	return msg, nil
}

// Parse BlockHeaders response.
func BlockHeadersRes(msg Msg) (*blockHeadersRes, error) {
	headers := new(blockHeadersRes)

	err := rlp.DecodeBytes(msg.Data, headers)
	if err != nil {
		return nil, err
	}

	return headers, nil
}

// Create GetBlockBodies request message.
func BlocksReq(headerHashes []common.Hash) (Msg, error) {
	reqId := rand.Uint64()

	// TODO: replace []any with struct.
	data, err := rlp.EncodeToBytes([]any{reqId, headerHashes})
	if err != nil {
		return Msg{}, nil
	}

	return NewMsg(GetBlockBodiesMsg, data), nil
}

type pooledTransactions struct {
	Types  []byte
	Sizes  []uint32
	Hashes []common.Hash
}

// Parse pooled transactions message that we received from peer.
func PooledTransactions(msg Msg) (*pooledTransactions, error) {
	pooledTxs := new(pooledTransactions)

	err := rlp.DecodeBytes(msg.Data, pooledTxs)
	if err != nil {
		return nil, nil
	}

	return pooledTxs, nil
}

// Parse the transaction message that was sent to us during the broadcast.
func NewTransactions(msg Msg) (*[]types.Transaction, error) {
	txs := new([]types.Transaction)

	err := rlp.DecodeBytes(msg.Data, txs)
	if err != nil {
		return nil, nil
	}

	return txs, err
}






