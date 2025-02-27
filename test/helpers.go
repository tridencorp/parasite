package test

import (
	"crypto/rand"
	"math/big"
	"parasite/p2p"
)

func RandInt(max int) *big.Int {
	rnd, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return rnd
}

// Generate given number of headers.
func Headers(n int) []*p2p.BlockHeader {
	headers := []*p2p.BlockHeader{}

	for i:=0; i<n; i++ {
		header := &p2p.BlockHeader{
			Number: RandInt(100_000_000),
		}

		headers = append(headers, header)
	}

	return headers
}