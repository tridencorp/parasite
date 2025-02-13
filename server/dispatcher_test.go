package server

import (
	"parasite/p2p"
	"testing"
)

func TestMessageDispatch(t *testing.T) {
	// Messages should be dispatched to correct handlers.
	// Handler should receive the same message that the dispatcher received.

	msg1 := p2p.NewMsg(p2p.BlockHeadersMsg,    []byte("BlockHeadersMsg"))
	msg3 := p2p.NewMsg(p2p.TransactionsMsg,    []byte("TransactionsMsg"))
	msg4 := p2p.NewMsg(p2p.BlockBodiesMsg,     []byte("BlockBodiesMsg"))     
	msg5 := p2p.NewMsg(p2p.ReceiptsMsg,        []byte("ReceiptsMsg"))        
	msg6 := p2p.NewMsg(p2p.GetBlockHeadersMsg, []byte("GetBlockHeadersMsg"))
	msg7 := p2p.NewMsg(p2p.GetBlockBodiesMsg,  []byte("GetBlockBodiesMsg"))
	msg8 := p2p.NewMsg(p2p.GetReceiptsMsg,     []byte("GetReceiptsMsg"))
	msg9 := p2p.NewMsg(p2p.NewPooledTransactionHashesMsg, []byte("NewPooledTransactionHashesMsg"))

	tests := []struct {
		name 	   string
		message  p2p.Msg
		expected p2p.Msg
	}{
		{ name: "BlockHeadersMsg",    message: msg1, expected: msg1},
		{ name: "TransactionsMsg",    message: msg3, expected: msg3},
		{ name: "BlockBodiesMsg",     message: msg4, expected: msg4},
		{ name: "ReceiptsMsg",        message: msg5, expected: msg5},
		{ name: "GetBlockHeadersMsg", message: msg6, expected: msg6},
		{ name: "GetBlockBodiesMsg",  message: msg7, expected: msg7},
		{ name: "GetReceiptsMsg",     message: msg8, expected: msg8},
		{ name: "NewPooledTransactionHashesMsg", message: msg9, expected: msg9},
	}

	dispatcher := NewDispatcher()
	handler, _ := dispatcher.Channels()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			dispatcher.Dispatch(nil, test.message)
			result := <- handler

			if result.Code != test.expected.Code {
				t.Errorf("Expected %d, got %d", test.expected, result)
			}
		})
	}
}
