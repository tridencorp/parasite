package p2p

import (
	"crypto/ecdsa"

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
	handshake.ID   = crypto.FromECDSAPub(srcPub)[1:]
	handshake.Caps = []Capability{{"eth", 68}} 

	// This will disable the snappy compression.
	handshake.Version = 0

	buf, err := rlp.EncodeToBytes(handshake)
	if err != nil {
		return err
	}

	res := NewMsg(HandshakeMsg, buf)
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
