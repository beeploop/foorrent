package client

import (
	"bytes"
	"errors"
	"net"
	"time"

	"github.com/beeploop/foorrent/handshake"
)

func performHandshake(conn net.Conn, peerID, infoHash [20]byte) (*handshake.HandShake, error) {
	conn.SetDeadline(time.Now().Add(client_timeout))
	defer conn.SetDeadline(time.Time{})

	handshakeRequest := handshake.New(infoHash, peerID)
	if _, err := conn.Write(handshakeRequest.Serialize()); err != nil {
		return nil, err
	}

	res, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}

	if !bytes.Equal(res.InfoHash[:], infoHash[:]) {
		return nil, errors.New("handshake response mismatched data, sever connection")
	}

	return res, nil
}
