package p2p

import (
	"parasite/block"

	"github.com/ethereum/go-ethereum/rlp"
)

// Handle status msg.
func HandleStatus(req Msg, peer *Peer) error {
	// @HACK: We just resend the same status.
	_, err := peer.Send(req)
	if err != nil {
		return err
	}

	return nil
}

// Decode block headers that we got from peer.
func HandleBlockHeaders(msg Msg) ([]*block.BlockHeader, error) {
	headers := new(BlockHeaders)

	err := rlp.DecodeBytes(msg.Data, headers)
	if err != nil {
		return nil, err
	}

	return headers.Headers, nil
}