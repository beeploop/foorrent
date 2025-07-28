package p2p

import (
	"log"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/peer"
)

type PeerToPeer struct {
	Peers       []peer.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
}

func (p *PeerToPeer) InitiateDownloadProcess() error {
	for _, peer := range p.Peers {
		peer := peer
		go func() {
			c, err := client.New(peer, p.PeerID, p.InfoHash)
			if err != nil {
				log.Printf("Failed Peer: %s, Error: %s\n", peer.String(), err.Error())
				return
			}
			defer c.Conn.Close()

			state := &PeerState{
				Choked:   true,
				BitField: make(bitfield.BitField, len(p.PieceHashes)),
			}
			processMessages(c, state)
		}()
	}

	select {} // TODO: Handle graceful waiting and exit
}
