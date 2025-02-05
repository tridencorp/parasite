package p2p

import (
	"bytes"
	"crypto/ecdsa"
	"parasite/block"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

// Parse handshake msg from remote peer and sends our response.
func HandleHandshake(req Msg, peer *Peer, srcPub *ecdsa.PublicKey) error {
	handshake := Handshake{}

	// Remote handshake message
	err := rlp.DecodeBytes(req.Data, &handshake)
	if err != nil {
		return err
	}

	// First byte is only a prefix that indicates if the key is compressed. We can omit it.
	handshake.ID = crypto.FromECDSAPub(srcPub)[1:]

	// We won't handle snap protocol for now, leave only newest eth.
	handshake.Caps = []Capability{{"eth", ETH}}

	// This will disable the snappy compression.
	handshake.Version = 0

	buf := bytes.Buffer{}
	err = rlp.Encode(&buf, handshake)
	if err != nil {
		return err
	}

	res := NewMsg(HandshakeMsg, buf.Bytes())
	_, err = peer.Send(res)
	if err != nil {
		return err
	}

	return nil
}

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