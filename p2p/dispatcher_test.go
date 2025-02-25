package p2p

import (
	"testing"
)

func TestMessageDispatch(t *testing.T) {
	// Messages should be dispatched to correct handlers.
	// Handler should receive the same message that the dispatcher received.

	msg1 := NewMsg(BlockHeadersMsg,    []byte("BlockHeadersMsg"))
	msg3 := NewMsg(TransactionsMsg,    []byte("TransactionsMsg"))
	msg4 := NewMsg(BlockBodiesMsg,     []byte("BlockBodiesMsg"))     
	msg5 := NewMsg(ReceiptsMsg,        []byte("ReceiptsMsg"))        
	msg6 := NewMsg(GetBlockHeadersMsg, []byte("GetBlockHeadersMsg"))
	msg7 := NewMsg(GetBlockBodiesMsg,  []byte("GetBlockBodiesMsg"))
	msg8 := NewMsg(GetReceiptsMsg,     []byte("GetReceiptsMsg"))
	msg9 := NewMsg(NewPooledTransactionHashesMsg, []byte("NewPooledTransactionHashesMsg"))

	tests := []struct {
		name 	   string
		message  *Msg
		expected *Msg
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

	response, failure := make(chan *Msg), make(chan *Msg)
	dispatcher := NewDispatcher(response, failure)
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
