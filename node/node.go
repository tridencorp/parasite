package node

import (
	"bytes"
	"crypto/ecdsa"
	"net"
	"parasite/key"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
)

type Cap struct {
	Name    string
	Version uint
}

type Handshake struct {
	Version    uint64
	Name       string
	Caps       []Cap
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

	setHandshakeFields(&handshake, *srcPrv)

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
	err = exchangeStatus(dst)
	if err != nil {
		return nil, err
	}

	return p2p.NewPeer(dst), nil
}

// Modifying our handshake.
func setHandshakeFields(handshake *Handshake, srcPrv ecdsa.PrivateKey) {
	// ID is basically our servers public key.
	pub := srcPrv.PublicKey
	handshake.ID = key.PubToBytes(&pub)

	// We support only eth protocol for now.
	handshake.Caps = []Cap{{"eth", p2p.ETH}}

	// This will disable the snappy compression.
	handshake.Version = 0
}

// @HACK: We just resend the same status that we got from remote peer.
func exchangeStatus(dst *rlpx.Conn) error {
	_, data, _, err := dst.Read()
	if err != nil {
		return err
	}

	_, err = dst.Write(p2p.StatusMsg, data)
	if err != nil {
		return err
	}

	return nil
}