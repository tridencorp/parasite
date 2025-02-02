package node

import (
	"crypto/ecdsa"
	"fmt"
	"net"

	"github.com/ethereum/go-ethereum/p2p/enode"
)

// Return enode string. Address must be in "ip:port" format.
func Enode(pub *ecdsa.PublicKey, address string) (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return "", err
	}

	node := enode.NewV4(pub, addr.IP, addr.Port, 0)
	return node.String(), nil
}

// Get public key and IP address from enode string.
func ParseEnode(rawEnode string) (*ecdsa.PublicKey, string, error) {
	node, err := enode.ParseV4(rawEnode)
	if err != nil {
		return nil, "", err
	}

	address, ok := node.TCPEndpoint()
	if !ok {
		// @CLEAN: At this point, the IP should have already been validated.
		return nil, "",	fmt.Errorf("cannot parse node IP address: %s", node.IPAddr())
	}

	return node.Pubkey(), address.String(), nil
}
