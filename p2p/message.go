package p2p

// Message sent over the p2p network.
type Msg struct {
  Code  uint64
  Size  uint32
  Data  []byte
  ReqId uint64
  
  // Each message has handler to which response 
  // for requested message will be send.
  Handler chan Msg
}

func NewMsg(code int, data []byte) Msg {
	return Msg{
		Code: uint64(code),
		Size: uint32(len(data)),
		Data: data,
	}
}
