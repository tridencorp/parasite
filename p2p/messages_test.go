package p2p

import (
	"parasite/key"
	"testing"
)

var local  = "enode://e806157dfc5e11365210e09ad4af4fb129024de4a8c97c7a6c834daf9567200f9e8d03a769c1d53f13286a643e295b6f38073e90a1833c1e06ef23cc402cfecb@127.0.0.1:30303"
var prv, _ = key.Private()

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
}
