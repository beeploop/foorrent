package message

import (
	"encoding/binary"
	"io"
)

type messageID uint8

const (
	MsgChoke        messageID = 0
	MsgUnchoke      messageID = 1
	MsgInterested   messageID = 2
	MsgUninterested messageID = 3
	MsgHave         messageID = 4
	MsgBitfield     messageID = 5
	MsgRequest      messageID = 6
	MsgPiece        messageID = 7
	MsgCancel       messageID = 8
)

type Message struct {
	ID      messageID
	Payload []byte
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}

	length := uint32(len(m.Payload) + 1) // 1 for id
	buf := make([]byte, length+4)        // 4 byte for length prefix
	binary.BigEndian.PutUint32(buf[:4], length)
	buf[4] = byte(m.ID)
	copy(buf[5:], m.Payload)

	return buf
}

func Read(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// length of 0 means keep-alive
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	if _, err := io.ReadFull(r, messageBuf); err != nil {
		return nil, err
	}

	m := &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}

	return m, nil
}
