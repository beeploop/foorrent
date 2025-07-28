package client

import (
	"fmt"
	"net"
	"time"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/message"
)

func readBitField(conn net.Conn) (bitfield.BitField, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{})

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		err := fmt.Errorf("Expected bitfield but got %v\n", msg)
		return nil, err
	}

	if msg.ID != message.MsgBitfield {
		err := fmt.Errorf("Expected bitfield but got message ID: %d", msg.ID)
		return nil, err
	}

	return msg.Payload, nil
}
