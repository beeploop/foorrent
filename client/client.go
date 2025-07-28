package client

import (
	"net"
	"time"

	"github.com/beeploop/foorrent/peer"
)

const (
	client_timeout = time.Second * 5
)

type Client struct {
	Conn     net.Conn
	Peer     peer.Peer
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), client_timeout)
	if err != nil {
		return nil, err
	}

	if _, err := performHandshake(conn, peerID, infoHash); err != nil {
		conn.Close()
		return nil, err
	}

	client := &Client{
		Conn:     conn,
		Peer:     peer,
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return client, nil
}
