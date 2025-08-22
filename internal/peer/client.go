package peer

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/beeploop/foorrent/internal/handshake"
)

type client struct {
	conn net.Conn
}

func createClient(ip net.IP, port uint16) (*client, error) {
	addr := net.JoinHostPort(ip.String(), strconv.Itoa(int(port)))
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return nil, err
	}

	client := &client{
		conn: conn,
	}

	return client, nil
}

func (c *client) close() {
	c.conn.Close()
}

func (c *client) handshake(peerID, infoHash [20]byte) (*handshake.HandShake, error) {
	c.conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer c.conn.SetDeadline(time.Time{})

	handshakeRequest := handshake.New(infoHash, peerID)
	if _, err := c.conn.Write(handshakeRequest.Serialize()); err != nil {
		return nil, err
	}

	res, err := handshake.Read(c.conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, errors.New("handshake response mismatched data, sever connection")
	}

	return res, nil
}
