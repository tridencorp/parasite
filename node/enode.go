package node

import (
	"crypto"
	"crypto/ecdsa"
	"net"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Return enode string. Address must be in "ip:port" format.
func Enode(pub crypto.PublicKey, address string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return "", err
	}

	node := enode.NewV4(pub.(*ecdsa.PublicKey), addr.IP, addr.Port, 0)
	return node.String(), nil
}
