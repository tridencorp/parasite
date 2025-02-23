package p2p

import (
	"crypto/ecdsa"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/p2p/discover/v4wire"
	"github.com/ethereum/go-ethereum/p2p/enode"
)

type Node struct {
	LocalPrv  *ecdsa.PrivateKey
	LocalPort  uint16
	LocalAddr  string

	RemotePub  v4wire.Pubkey
	AddrPort   string
	Addr       string
	Port       uint16
	Conn  	  *net.UDPConn
} 

var Expiration = uint64(time.Now().Add(20 * time.Second).Unix())
var Deadline   = time.Now().Add(2 * time.Second)
var BuffSize   = 20_000
var Rounds     = 2

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

	addres   := strings.Split(conn.LocalAddr().String(), ":")
	port     := addres[1]
	value, _ := strconv.ParseUint(port, 10, 16)

	node := &Node{
		LocalPrv:  prv,
		LocalPort: uint16(value),
		LocalAddr: addres[0], 
		RemotePub: pubkey,
		AddrPort:  addrport.String(),
		Addr:      addrport.Addr().String(),
		Port:      addrport.Port(),
		Conn:      conn,
	}

	return node, nil
}

func (node *Node) Findnode() (packet v4wire.Packet, hash []byte, err error) {
	query := &v4wire.Findnode{
		Target:     node.RemotePub,
		Expiration: Expiration,
	}

	return node.Send(query)
}

func (node *Node) Ping() (packet v4wire.Packet, hash []byte, err error) {
	local  := v4wire.Endpoint{IP: net.ParseIP(node.LocalAddr), UDP: node.LocalPort, TCP: node.LocalPort}
	remote := v4wire.Endpoint{IP: net.ParseIP(node.Addr),      UDP: node.Port,      TCP: node.Port}

	query := &v4wire.Ping{
		Version: 4,
		From: local,
		To: remote,
		Expiration: Expiration,
	}

	return node.Send(query)
}

func (node *Node) Pong(hash []byte) (v4wire.Packet, []byte, error) {
	remote := v4wire.Endpoint{IP: net.ParseIP(node.Addr), UDP: node.Port, TCP: node.Port}
	query  := &v4wire.Pong{To: remote, ReplyTok: hash, Expiration: Expiration}

	return node.Send(query)
}

func (node *Node) Send(query v4wire.Packet) (v4wire.Packet, []byte, error) {
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

			res, hash, err := node.Findnode()
			if err != nil {
				fmt.Println(err)
				node.Conn.Close()
				continue
			}

			for {
				if res.Name() == "PING/v4" {
					res, hash, err = node.Pong(hash)
					if err != nil {
						fmt.Println(err)
						node.Conn.Close()
						break
					}
				}	

				if res.Name() == "NEIGHBORS/v4" {
					ParseNeighbors(&newNodes, res.(*v4wire.Neighbors))
					node.Conn.Close()
					break
				}
	
				if res.Name() == "PONG/v4" {
					res, hash, err = node.Findnode()
					if err != nil {
						fmt.Println(err)
						node.Conn.Close()
						continue
					}							
				}	
			}
		}

		fmt.Println("Number of new nodes: ", len(newNodes))
		knownNodes = newNodes
	}

	SaveNodes(newNodes, "nodes.txt")
}
