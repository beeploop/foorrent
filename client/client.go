package client

import (
	"encoding/binary"
	"log"
	"net"
	"time"

	"github.com/beeploop/foorrent/message"
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

func (c *Client) Close() {
	c.Conn.Close()
	log.Printf("connection close, [peer %s]\n", c.Peer.String())
}

func (c *Client) SendKeepAlive() error {
	if _, err := c.Conn.Write([]byte{0, 0, 0, 0}); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendRequest(index, begin, length int) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	msg := &message.Message{
		ID:      message.MsgRequest,
		Payload: payload,
	}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendInterested() error {
	msg := &message.Message{ID: message.MsgInterested}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendUninterested() error {
	msg := &message.Message{ID: message.MsgUninterested}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendChoke() error {
	msg := &message.Message{ID: message.MsgChoke}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendUnchoke() error {
	msg := &message.Message{ID: message.MsgUnchoke}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (c *Client) SendHave(index int) error {
	msg := &message.Message{ID: message.MsgHave}
	if _, err := c.Conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}
