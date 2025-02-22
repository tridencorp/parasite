package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type Node struct {
	LocalPrv  *ecdsa.PrivateKey
	RemotePub  v4wire.Pubkey
	AddrPort   string
	Addr 			 string
	Port       uint16
	Conn  	  *net.UDPConn
} 

var Expiration = uint64(time.Now().Add(20 * time.Second).Unix())
var Deadline   = time.Now().Add(2 * time.Second)
var BuffSize   = 20_000

func ConnectNode(prv *ecdsa.PrivateKey, addr string) (*Node, error) {
	enode,    _ := enode.ParseV4(addr)
	addrport, _ := enode.UDPEndpoint()
	pubkey      := v4wire.EncodePubkey(enode.Pubkey())

	remote, err := net.ResolveUDPAddr("udp", addrport.String())
	if err != nil {
		return nil, err
	}
	
	conn, err := net.DialUDP("udp", nil, remote)
	if err != nil {
		return nil, err
	}
	
	node := &Node{
		LocalPrv:  prv, 
		RemotePub: pubkey,
		AddrPort:  addrport.String(),
		Addr:      addrport.Addr().String(),
		Port:      addrport.Port(),
		Conn: 		 conn,
	}

	return node, nil
}

func (node *Node) Findnode() (packet v4wire.Packet, hash []byte, err error) {
	query := &v4wire.Findnode{
    Target:     node.RemotePub,
    Expiration: Expiration,
  }

	req, hash, err := v4wire.Encode(node.LocalPrv, query)
	if err != nil {
		return nil, nil, err
	}

	_, err = node.Conn.Write(req)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Clean this.
	buff := make([]byte, BuffSize)
	node.Conn.SetDeadline(Deadline)

	size, _, err := node.Conn.ReadFromUDP(buff)
	if err != nil {
		return nil, nil, err
	}

	packet, _, hash, err = v4wire.Decode(buff[:size])
	if err != nil {
		return nil, nil, err
	}

	return packet, hash, nil
}

func (node *Node) Ping() (packet v4wire.Packet, hash []byte, err error) {
	// TODO: clean this.
	local  := v4wire.Endpoint{IP: net.ParseIP("192.168.1.167"), UDP: 60559, TCP: 30304}
	remote := v4wire.Endpoint{IP: net.ParseIP(node.Addr), UDP: node.Port, TCP: node.Port}

	query := &v4wire.Ping{
		Version: 4,
		From: local,
		To:   remote,
		Expiration: Expiration,
  }
	
	req, hash, err := v4wire.Encode(node.LocalPrv, query)
	if err != nil {
		return nil, nil, err
	}

	_, err = node.Conn.Write(req)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Clean this.
	buff := make([]byte, BuffSize)
	node.Conn.SetDeadline(Deadline)

	size, _, err := node.Conn.ReadFromUDP(buff)
	if err != nil {
		return nil, nil, err
	}

	packet, _, hash, err = v4wire.Decode(buff[:size])
	if err != nil {
		return nil, nil, err
	}

	return packet, hash, nil
}

func (node *Node) Pong(hash []byte) (v4wire.Packet, []byte, error) {
	remote := v4wire.Endpoint{IP: net.ParseIP(node.Addr), UDP: node.Port, TCP: node.Port}
	query  := &v4wire.Pong{To: remote, ReplyTok: hash, Expiration: Expiration}

	req, _, err := v4wire.Encode(node.LocalPrv, query)
	if err != nil {
		return nil, nil, err
	}

	_, err = node.Conn.Write(req)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Clean this.
	buff := make([]byte, BuffSize)
	node.Conn.SetDeadline(Deadline)

	size, _, err := node.Conn.ReadFromUDP(buff)
	if err != nil {
		return nil, nil, err
	}

	packet, _, hash, _ := v4wire.Decode(buff[:size])
	return packet, hash, nil
}

func ParseNeighbors(nodes []string, neighbors *v4wire.Neighbors) error {
	for _, node := range neighbors.Nodes {
		pub, _  := v4wire.DecodePubkey(crypto.S256(), node.ID)
		address := fmt.Sprintf("%s:%d", node.IP.String(), node.TCP)

		enode, err := Enode(pub, address, int(node.UDP))
		if err != nil {
			return err
		}

		nodes = append(nodes, enode)
	}

	return nil
}

func Discover() {
	prv, _ := crypto.LoadECDSA("./pub")

	newNodes   := []string{}
	knownNodes := []string{
		"enode://4aeb4ab6c14b23e2c4cfdce879c04b0748a20d8e9b59e25ded2a08143e265c6c25936e74cbc8e641e3312ca288673d91f2f93f8e277de3cfa444ecdaaf982052@157.90.35.166:30303",
		"enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
		"enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",
		"enode://2b252ab6a1d0f971d9722cb839a42cb81db019ba44c08754628ab4a823487071b5695317c8ccd085219c3a03af063495b2f1da8d18218da2d6a82981b45e6ffc@65.108.70.101:30303",
		"enode://6c6d584572133f40fd5f31d4cb46434a426b19a07e01d12aa99a4a058d0fc9c797a1c2e7bac62483f3ff74d98b4d6fe1f7dd5da32bee6d12d237c2be67b777da@65.21.149.19:30303",
		"enode://61bcc1dd9d0ea84be7455f5126c033413a9ec81f630a94695a99513e6246bbe50cc68aa4183999989bdb213efc3e5ff08582b8d1940330fd2e780b0f36519558@125.132.1.160:30303",
		"enode://b74b01b1a5c207d574396c5954a5d771e1395278118c23a54ee1df18e9f163746fea8c3a9fe763408869bcc0543ef75438e4ab52c550978a1e54100d64b692dc@195.201.242.243:30303",
	}

	node, err := ConnectNode(prv, knownNodes[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	res, hash, err := node.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("RES: ", res.Name(), "NODE: ", node.AddrPort)

	if res.Name() == "PING/v4" {
		res, hash, err = node.Pong(hash)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if res.Name() == "NEIGHBORS/v4" {
		fmt.Println("RES 1: ", res.Name(), "NODE: ", node.AddrPort)
		ParseNeighbors(newNodes, res.(*v4wire.Neighbors))
		return
	}

	fmt.Println("RES: ", res.Name(), "NODE: ", node.AddrPort)

	res, hash, err = node.Findnode()
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.Name() == "NEIGHBORS/v4" {
		fmt.Println("RES 2: ", res.Name(), "NODE: ", node.AddrPort)
		ParseNeighbors(newNodes, res.(*v4wire.Neighbors))
		
		fmt.Println("New Nodes")
		fmt.Println(newNodes)
		return
	}

	fmt.Println("NO RESPONSE")
}
