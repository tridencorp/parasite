package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"net"
	"parasite/key"
	"time"

	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	PingPacket = iota + 1
	PongPacket
	FindnodePacket
	NeighborsPacket
	ENRRequestPacket
	ENRResponsePacket
)

type Node struct {
	Addr     string
	Enode    string
	PubKey  *ecdsa.PublicKey
	V4PubKey v4wire.Pubkey
}

func NewNode(address string) *Node {
	node, _ := enode.ParseV4(address)

	return &Node{ Enode: address, PubKey:  node.Pubkey(), V4PubKey: v4wire.EncodePubkey(node.Pubkey())}
}

type Findnode struct {
	Target     v4wire.Pubkey
	Expiration uint64

	// Not used for now.
	Rest []rlp.RawValue `rlp:"tail"`
}

func (req *Findnode) Name() string { return "FINDNODE/v4" }
func (req *Findnode) Kind() byte   { return FindnodePacket }

func Discover() {
	prv, _ := key.Private()

	nodes := []Node{
		// Bootnodes.
		*NewNode("enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303"),
		*NewNode("enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303"),
		*NewNode("enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303"),
		*NewNode("enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303"),
	}

	node := nodes[0]

	req := &v4wire.Findnode{
    Target:     node.V4PubKey,
    Expiration: uint64(time.Now().Add(100000).Unix()),
  }

	packet, _, _ := v4wire.Encode(prv, req)

	addr, err := net.ResolveUDPAddr("udp", "65.21.149.19:30303")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()


	fmt.Println("Sending request")
	_, err = conn.Write(packet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Reading response")
	buff := make([]byte, 20000)
	n, _, err := conn.ReadFromUDP(buff)
	if err != nil {
		log.Fatal(err)
	}

	pack, _, _, _ := v4wire.Decode(buff[:n])

	fmt.Println("Response:")
	fmt.Println(pack)
}