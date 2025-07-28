package p2p

import (
	"encoding/binary"
	"log"
	"sync"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/message"
)

type PeerState struct {
	mu       sync.Mutex
	Choked   bool
	BitField bitfield.BitField
}

func processMessages(c *client.Client, state *PeerState) {
	for {
		msg, err := message.Read(c.Conn)
		if err != nil {
			log.Println("Error reading message from peer: ", err.Error())
			return
		}
		if msg == nil {
			continue
		}

		switch msg.ID {
		case message.MsgChoke:
			log.Println("received choke from: ", c.Peer.IP)
			state.mu.Lock()
			state.Choked = true
			state.mu.Unlock()

		case message.MsgUnchoke:
			log.Println("received unchoke from: ", c.Peer.IP)
			state.mu.Lock()
			state.Choked = false
			state.mu.Unlock()

		case message.MsgInterested:
			log.Println("received interested from: ", c.Peer.IP)
			// TODO: Handle interested messages

		case message.MsgUninterested:
			log.Println("received uninterested from: ", c.Peer.IP)
			// TODO: Handle uninterested messages

		case message.MsgHave:
			log.Println("received have from: ", c.Peer.IP)
			index := binary.BigEndian.Uint32(msg.Payload)
			state.mu.Lock()
			state.BitField.SetPiece(int(index))
			state.mu.Unlock()

		case message.MsgBitfield:
			log.Println("received bitfield from: ", c.Peer.IP)
			state.mu.Lock()
			state.BitField = msg.Payload
			state.mu.Unlock()

		case message.MsgRequest:
			log.Println("received request from: ", c.Peer.IP)
			// TODO: Handle request messages

		case message.MsgPiece:
			log.Println("received piece from: ", c.Peer.IP)
			// TODO: Handle downloading piece data

		case message.MsgCancel:
			log.Println("received cancel from: ", c.Peer.IP)
			// TODO: Handle cancel messages

		default:
			log.Println("unknown message ID from: ", c.Peer.IP)
		}
	}
}
