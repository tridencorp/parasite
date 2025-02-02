package p2p

import (
	"bytes"
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

	// For now let's try to send the same handshake only changing the ID.
	buf := bytes.Buffer{}

	// First byte is only a prefix that indicates if the key is compressed. We can omit it.
	handshake.ID = crypto.FromECDSAPub(srcPub)[1:]

	err = rlp.Encode(&buf, handshake)
	if err != nil {
		return err
	}

	res := Msg{Code: HandshakeMsg, Size: uint32(buf.Len()), Data: buf.Bytes()}
	_, err = peer.Send(res)
	if err != nil {
		return err
	}

	return nil
}
