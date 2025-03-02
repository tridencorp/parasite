package p2p

import (
	"math/rand/v2"

	"github.com/ethereum/go-ethereum/rlp"
)

// Message sent over the p2p network.
type Msg struct {
  Code  uint64
  Size  uint32
  Data  []byte
  ReqId uint64

	Payload any

	// Each message has a handler to which it can be dispatched.
  Handler chan *Msg
}

func NewMsg(code int, data []byte) *Msg {
	return &Msg{
		Code: uint64(code),
		Size: uint32(len(data)),
		Data: data,
	}
}

func EncodeMsg(code int, data any) (*Msg, error) {
	req := Request{
		ReqID: rand.Uint64(),
		Data: data,
	}

	bytes, err := rlp.EncodeToBytes(req)
	if err != nil {
		return nil, err
	}

	msg := NewMsg(code, bytes)
	msg.ReqId = req.ReqID

	return msg, nil
}

func DecodeMsg(bytes []byte, dst any) error {
	return rlp.DecodeBytes(bytes, dst)
}
