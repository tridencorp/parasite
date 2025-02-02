package p2p

// Message sent over the p2p network.
type Msg struct {
	Code uint64
	Size uint32
	Data []byte
}
