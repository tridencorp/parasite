package test

import (
	"crypto/rand"
	"math/big"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
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

// Create p2p headers response message.
func HeadersResponseMsg(code, num int) (*p2p.Msg, []*p2p.BlockHeader) {
	headers := Headers(num)
	data    := []any{uint64(code), headers}
	req, _  := rlp.EncodeToBytes(data)
	msg     := p2p.NewMsg(1, req)
	
	return msg, headers
}
