package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"os"
	"strings"
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
var Rounds 		 = 2

func ConnectNode(prv *ecdsa.PrivateKey, addr string) (*Node, error) {
	enode,    _ := enode.ParseV4(addr)
	addrport, _ := enode.UDPEndpoint()
	pubkey      := v4wire.EncodePubkey(enode.Pubkey())

	remote, err := net.ResolveUDPAddr("udp", addrport.String())
	if err != nil {
		return nil, err
	}
	
	fmt.Println("REMOTE: ", remote)
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
	
	return node.Send(req)
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

	return node.Send(req)
}

func (node *Node) Pong(hash []byte) (v4wire.Packet, []byte, error) {
	remote := v4wire.Endpoint{IP: net.ParseIP(node.Addr), UDP: node.Port, TCP: node.Port}
	query  := &v4wire.Pong{To: remote, ReplyTok: hash, Expiration: Expiration}

	req, _, err := v4wire.Encode(node.LocalPrv, query)
	if err != nil {
		return nil, nil, err
	}

	return node.Send(req)
}

func (node *Node) Send(req []byte) (v4wire.Packet, []byte, error) {
	_, err := node.Conn.Write(req)
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

func ParseNeighbors(nodes *[]string, neighbors *v4wire.Neighbors) error {
	for _, node := range neighbors.Nodes {
		pub, _  := v4wire.DecodePubkey(crypto.S256(), node.ID)
		address := fmt.Sprintf("%s:%d", node.IP.String(), node.TCP)

		enode, err := Enode(pub, address, int(node.UDP))
		if err != nil {
			return err
		}

		*nodes = append(*nodes, enode)
	}

	return nil
}

func LoadNodes(filename string) ([]string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return strings.Split(string(data), "\n"), nil
}

func SaveNodes(nodes []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	
	defer file.Close()

	for i, line := range nodes {
		if i == len(nodes)-1 {
			_, err := fmt.Fprint(file, line)
			if err != nil {
				return err
			}
			break
		}

		_, err := fmt.Fprintln(file, line)
		if err != nil {
			return err
		}
	}

	return nil
}

func Discover() {
	prv, _ := crypto.LoadECDSA("./pub")

	newNodes      := []string{}
	knownNodes, _ := LoadNodes("bootnodes.txt")

	for i:=0; i < Rounds; i++ {
		for _, enode := range knownNodes {
			node, err := ConnectNode(prv, enode)
			if err != nil {
				fmt.Println(err)
				continue
			}

			res, hash, err := node.Ping()
			if err != nil {
				fmt.Println(err)
				continue
			}

			if res.Name() == "PING/v4" {
				res, hash, err = node.Pong(hash)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}

			if res.Name() == "NEIGHBORS/v4" {
				ParseNeighbors(&newNodes, res.(*v4wire.Neighbors))
				continue
			}
			
			res, hash, err = node.Findnode()
			if err != nil {
				fmt.Println(err)
				continue
			}
			
			if res.Name() == "NEIGHBORS/v4" {
				ParseNeighbors(&newNodes, res.(*v4wire.Neighbors))
				continue
			}
		}
		
		knownNodes = newNodes
	}
	
	SaveNodes(newNodes, "nodes.txt")
	fmt.Println(len(newNodes))
}
