package download

import (
	"log"

	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/message"
)

func (dm *DownloadManager) handleInboundMessages(c *client.Client, session *peerSession) {
	for {
		msg, err := message.Read(c.Conn)
		if err != nil {
			log.Printf("message error: %s, peer: %s\n", err.Error(), c.Peer.String())
			dm.mu.Lock()
			delete(dm.activePeers, c.Peer.String())
			dm.mu.Unlock()
			c.Close()
			return
		}
		// Received keep-alive message
		if msg == nil {
			continue
		}

		switch msg.ID {
		case message.MsgChoke:
			session.choke()

		case message.MsgUnchoke:
			session.unchoke()

		case message.MsgInterested:
			log.Println("received interested from: ", c.Peer.IP)
			// TODO: Handle interested messages

		case message.MsgUninterested:
			log.Println("received uninterested from: ", c.Peer.IP)
			// TODO: Handle uninterested messages

		case message.MsgHave:
			session.putPiece(msg.Payload)

		case message.MsgBitfield:
			session.setBitfield(msg.Payload)

		case message.MsgRequest:
			log.Println("received request from: ", c.Peer.IP)
			// TODO: Handle request messages

		case message.MsgPiece:
			block, err := parsePiece(msg.Payload)
			if err != nil {
				log.Println(err.Error())
				continue
			}
			dm.downloadChan <- block

		case message.MsgCancel:
			log.Println("received cancel from: ", c.Peer.IP)
			// TODO: Handle cancel messages

		default:
			log.Println("unknown message ID from: ", c.Peer.IP)
		}
	}
}
