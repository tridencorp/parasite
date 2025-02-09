package node

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"parasite/key"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/p2p/rlpx"
	"github.com/ethereum/go-ethereum/rlp"
)

type Handshake struct {
	Version uint64
	Name string

	Caps []struct{ 
		Name string
		Version uint
	}

	ListenPort uint64
	ID []byte

	// Currently unused, but required for compatibility with ETH.
	Rest []rlp.RawValue `rlp:"tail"`
}

// Connect to given node, perform handshake and exchange status msg.
func Connect(enode string, srcPrv *ecdsa.PrivateKey) (*p2p.Peer, error) {
	conn, err := initialHandshake(enode, srcPrv)
	if err != nil {
		return nil, err
	}

	err = handleHandshakeMessage(conn, srcPrv.PublicKey)
	if err != nil {
		return nil, err
	}

	err = exchangeStatus(conn)
	if err != nil {
		return nil, err
	}

	return p2p.NewPeer(conn), nil
}

func initialHandshake(enode string, prv *ecdsa.PrivateKey) (*rlpx.Conn, error)  {
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
	_, err = dst.Handshake(prv)
	if err != nil {
		return nil, err
	}

	return dst, nil
}

func handleHandshakeMessage(conn *rlpx.Conn, pub ecdsa.PublicKey) error {
	handshake := Handshake{}

	_, data, _, err := conn.Read()
	if err != nil {
		return fmt.Errorf("Handshake | error while reading from peer")
	}

	// Decode remote handshake message.
	err = rlp.DecodeBytes(data, &handshake)
	if err != nil {
		return fmt.Errorf("Handshake | cannot decode bytes from peer")
	}

	// ID is basically our servers public key.
	handshake.ID = key.PubToBytes(&pub)

	// This will disable the snappy compression.
	handshake.Version = 0

	buf, err := rlp.EncodeToBytes(handshake)
	if err != nil {
		return err
	}

	_, err = conn.Write(p2p.HandshakeMsg, buf)
	if err != nil {
		return err
	}

	return nil
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
