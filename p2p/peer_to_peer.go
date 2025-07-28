package p2p

import (
	"fmt"

	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/peer"
)

type PeerToPeer struct {
	Peers    []peer.Peer
	PeerID   [20]byte
	InfoHash [20]byte
}

func (p *PeerToPeer) InitiateDownloadProcess() error {
	for _, peer := range p.Peers {
		c, err := client.New(peer, p.PeerID, p.InfoHash)
		if err != nil {
			fmt.Printf("Failed Peer: %s, Error: %s\n", peer.String(), err.Error())
			continue
		}
		defer c.Conn.Close()

		fmt.Println(c.BitField)
	}

	return nil
}
