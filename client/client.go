package client

import (
	"net"
	"time"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/peer"
)

type Client struct {
	Conn     net.Conn
	Choked   bool
	BitField bitfield.BitField
	Peer     peer.Peer
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), time.Second*5)
	if err != nil {
		return nil, err
	}

	if _, err := performHandshake(conn, peerID, infoHash); err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := readBitField(conn)
	if err != nil {
		return nil, err
	}

	client := &Client{
		Conn:     conn,
		Choked:   true,
		BitField: bf,
		Peer:     peer,
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return client, nil
}
