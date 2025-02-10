package server

import (
	"fmt"
	"parasite/log"
	"parasite/p2p"

	"github.com/ethereum/go-ethereum/rlp"
)

// Dispatching received messages to designated handlers.
type Dispatcher struct {
	peer *p2p.Peer

	handler chan p2p.Msg
	failure chan p2p.Msg
}

// Create new Dispatcher,
func NewDispatcher(peer *p2p.Peer, handler chan p2p.Msg, failure chan p2p.Msg) *Dispatcher {
	return &Dispatcher{peer, handler, failure}
}

// Main dispatcher responsible for dispatching all incomming messages.
// It uses 2 channels: one for normal message handling and another one 
// for sending errors.
func (dsp *Dispatcher) Dispatch(msg p2p.Msg) { 
  // (2) PingMsg
  if msg.Code == p2p.PingMsg {
    log.Info("!!! Got Ping !!!")
    dsp.peer.Send(p2p.NewMsg(p2p.PongMsg, []byte{}))
    return
  }

  // (19) GetBlockHeadersMsg
  if msg.Code == p2p.GetBlockHeadersMsg {
    log.Info("%d", msg.Code)
    log.Info("%v", msg.Data)    
    return
  }

  // (20) BlockHeadersMsg
  if msg.Code == p2p.BlockHeadersMsg {
    log.Info("!!! Got headers !!!")

    res, err := p2p.BlockHeadersRes(msg)
    if err != nil {
      log.Error("%v", err)
    }

    // Send msg with requsted headers to our handler.
    reqMsg, exists := dsp.peer.RequestedMsgs[res.ReqId]
    if exists {
      reqMsg.Handler <- msg
    }

    return
  }

  // (1) DiscMsg
  if msg.Code == p2p.DiscMsg {
    log.Error("!!! DISCONECT FROM NODE !!!")

    type DiscReason uint8
    var disc []DiscReason

    rlp.DecodeBytes(msg.Data, &disc)
    log.Error("Disconnect from peer: %s", p2p.DiscReasons[disc[0]])
    return
  }

  if msg.Code == p2p.NewPooledTransactionHashesMsg { 
    log.Info("Request: %d : NewPooledTransactions", msg.Code)

    pooledTx, err := p2p.PooledTransactions(msg)
    if err != nil {
      fmt.Println(err)
    }
    
    fmt.Printf("%v", pooledTx)
    return
  }

  if msg.Code == p2p.TransactionsMsg {
    log.Info("Request: %d : p2p.TransactionsMsg", msg.Code)

    txs, err := p2p.NewTransactions(msg)
    if err != nil {
      fmt.Println(err)
    }

    fmt.Printf("%v", txs)
    return
  }

  // (21) GetBlockBodiesMsg
  // We got request for block bodies.
  if msg.Code == p2p.GetBlockBodiesMsg {
    log.Info("Request: %d : p2p.GetBlockBodiesMsg", msg.Code)

    // Find the requested blocks and send them back to peer.
    // For now we don't have any data so we send empty response.
    blockBodies, _ := p2p.BlockBodiesRes(msg)

    data, err := rlp.EncodeToBytes(blockBodies)
    if err != nil {
      log.Error(err.Error())
    }

    msg := p2p.NewMsg(p2p.GetBlockBodiesMsg, data)
    fmt.Println(msg)
    return
  }

  // @TODO: Needs to be implemented
  if msg.Code == p2p.BlockBodiesMsg { log.Error("Implement %d", p2p.BlockBodiesMsg)    ;return }
  if msg.Code == p2p.GetReceiptsMsg { log.Error("Implement %d", p2p.GetReceiptsMsg)    ;return }
  if msg.Code == p2p.ReceiptsMsg    { log.Error("Implement %d", p2p.ReceiptsMsg)       ;return }

  // If we are here then we have unsupported message. 
  // Just print it for now.
  log.Error("Unknown msg code: %d\n", msg.Code)
}
