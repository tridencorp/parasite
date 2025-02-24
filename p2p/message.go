package p2p

import "github.com/ethereum/go-ethereum/rlp"

// Message sent over the p2p network.
type Msg struct {
  Code  uint64
  Size  uint32
  Data  []byte
  ReqId uint64
  
	// Each message has a handler to which it can be dispatched.
  Handler chan Msg
}

func NewMsg(code int, data []byte) *Msg {
	return &Msg{
		Code: uint64(code),
		Size: uint32(len(data)),
		Data: data,
	}
}

func EncodeMsg(code int, data any) (*Msg, error) {
	bytes, err := rlp.EncodeToBytes(data)
	if err != nil {
		return nil, err
	}

	return NewMsg(code, bytes), nil
}
