package key

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

// Create ECDSA (secp256k1) private key from hex string
func FromHex(hex string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(hex)
}

// Generate new ECDSA private key.
func Private() (*ecdsa.PrivateKey, error) {
	return crypto.GenerateKey()
}

// Convert ecdsa.PubliceKey to bytes.
func PubToBytes(pub *ecdsa.PublicKey) []byte {
	// First byte is only a prefix that indicates if the key is compressed. 
	// We can remove it.
	return crypto.FromECDSAPub(pub)[1:]
}
