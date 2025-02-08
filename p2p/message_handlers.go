package p2p

import (
	"parasite/block"

	"github.com/ethereum/go-ethereum/rlp"
)

// Decode block headers that we got from peer.
func HandleBlockHeaders(msg Msg) ([]*block.BlockHeader, error) {
	headers := new(blockHeadersRes)

	err := rlp.DecodeBytes(msg.Data, headers)
	if err != nil {
		return nil, err
	}

	return headers.Headers, nil
}