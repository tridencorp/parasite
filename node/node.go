package node

import (
	"bytes"
	"crypto/ecdsa"
	"net"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
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

// Connect to given node, perform handshake and exchange status msg.
func Connect(enode string, srcPrv *ecdsa.PrivateKey) (*p2p.Peer, error) {
	// Get node's pub key and address.
	dstPub, address, err := ParseEnode(enode)
	if err != nil {
		return nil, err
	}

	// Connect to node.
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// Perform initial handshake.
	dst := rlpx.NewConn(conn, dstPub)
	_, err = dst.Handshake(srcPrv)
	if err != nil {
		return nil, err
	}

	// (0) HandshakeMsg: Perform post init handshake.
	handshake := Handshake{}

	_, data, _, err := dst.Read()
	if err != nil {
		return nil, err
	}

	// Decode remote handshake message.
	err = rlp.DecodeBytes(data, &handshake)
	if err != nil {
		return nil, err
	}

	// First byte is only a prefix that indicates if the key is compressed. We can omit it.
	srcPub := srcPrv.PublicKey
	handshake.ID = crypto.FromECDSAPub(&srcPub)[1:]

	// We won't handle snap protocol for now, leave only newest eth.
	handshake.Caps = []Capability{{"eth", p2p.ETH}}

	// This will disable the snappy compression.
	handshake.Version = 0

	buf := bytes.Buffer{}
	err = rlp.Encode(&buf, handshake)
	if err != nil {
		return nil, err
	}

	_, err = dst.Write(p2p.HandshakeMsg, buf.Bytes())
	if err != nil {
		return nil, err
	}

	// (16) StatusMsg: Exchange status msg.
	// 
	// @HACK: We just resend the same status that we got from remote peer.
	_, data, _, err = dst.Read()
	if err != nil {
		return nil, err
	}

	_, err = dst.Write(p2p.StatusMsg, data)
	if err != nil {
		return nil, err
	}

	return p2p.NewPeer(dst), nil
}
