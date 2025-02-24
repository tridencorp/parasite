package p2p

import (
	"encoding/hex"
	"fmt"
	"parasite/key"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

var local  = "enode://e806157dfc5e11365210e09ad4af4fb129024de4a8c97c7a6c834daf9567200f9e8d03a769c1d53f13286a643e295b6f38073e90a1833c1e06ef23cc402cfecb@127.0.0.1:30303"
var prv, _ = key.Private()

var BlockHash = "c2fa57aa06932cb8ee7ab0e2aec04dcbdd471902a7bc0c9b9777bbcdffdd0687"

func TestGetBlockHeadersMsg(t *testing.T) {
	p, _ := Connect(local, prv)
	p.Start()
	defer p.Close()

	p.GetBlockHeadersMsg(14_678_700, 1)

	msg := <-p.Response
	if len(msg.Data) != 553 {
		t.Errorf("Expected len to be %d, got %d", 553, len(msg.Data))
	}	
	
	headers, _ := DecodeBlockHeadersMsg(msg)
	if len(headers) != 1 {
		t.Errorf("Expected len to be %d, got %d", 1, len(headers))
	}	

	fmt.Println(headers[0].Hash())
}

func TestGetBlockBodiesMsg(t *testing.T) {
	p, _ := Connect(local, prv)
	p.Start()
	defer p.Close()

	hash, _ := hex.DecodeString(BlockHash)
	p.GetBlockBodiesMsg([]common.Hash{common.Hash(hash)})

	msg := <-p.Response

	if len(msg.Data) != 50543 {
		t.Errorf("Expected len to be %d, got %d", 50543, len(msg.Data))
	}	
}
