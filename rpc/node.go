package rpc

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand/v2"
	"net/http"
)

type Node struct {
  Address string
}

type Request struct {
	ID      int32  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

type Response struct {
	ID      int32  `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  any    `json:"result"`
}

// Create new node to which we will be sending our
// rpc requests. Address should be in [http|https]://[url|ip]:port format.
func NewNode(address string) *Node {
	return &Node{address}
}

// Send request to RPC node and unmarshal the response.
func (node *Node) Send(method string, response, params any) error {
	request := Request{rand.Int32(), "2.0", method, params}

	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	raw, err := http.Post(node.Address, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	defer raw.Body.Close()

	body, err := io.ReadAll(raw.Body)
	if err != nil {
		return err
	}

	// Unmarshal data to response.
	res := &Response{Result: response}
	err = json.Unmarshal(body, res)
	if err != nil {
		return err
	}

	return nil
}