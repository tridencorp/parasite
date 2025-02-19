package rpc

type Node struct {
	address string
}

// Create new node to which we will be sending our
// rpc requests. Address should be in ip:port format.
func NewNode(address string) *Node {
	return &Node{address}
}