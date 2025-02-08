package p2p

import (
	"math/rand/v2"
	"parasite/block"

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
	Start   uint64
	Amount  uint64
	Skip    uint64
	Reverse bool
}

type blockHeadersRes struct {
	ReqId   uint64
	Headers []*block.BlockHeader
}

// Create BlockHeaders request message.
func BlockHeadersReq(start, amount, skip uint64, reverse bool) (Msg, error) {
	reqID := rand.Uint64()
	
	req := blockHeadersReq{
		Start:   start,
		Amount:  amount,
		Skip:    skip,
		Reverse: false,
	}

	data, err := rlp.EncodeToBytes([]any{reqID, req})
	if err != nil {
		return Msg{}, err
	}

	return NewMsg(GetBlockHeadersMsg, data), nil
}
