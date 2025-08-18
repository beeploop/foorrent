package peer

import (
	"encoding/binary"

	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/message"
	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/tracker"
)

type Session struct {
	Choked     bool
	Interested bool
	BitField   bitfield.BitField
	Client     *client
}

func New(peerID [20]byte, peer tracker.Peer, torrent metadata.Torrent) (*Session, error) {
	hashes, err := torrent.PieceHashes()
	if err != nil {
		return nil, err
	}

	infoHash, err := torrent.InfoHash()
	if err != nil {
		return nil, err
	}

	c, err := createClient(peer.IP, peer.Port)
	if err != nil {
		return nil, err
	}

	if _, err := c.handshake(peerID, infoHash); err != nil {
		c.close()
		return nil, err
	}

	session := &Session{
		Choked:     true,
		Interested: false,
		BitField:   make(bitfield.BitField, len(hashes)),
		Client:     c,
	}

	return session, nil
}

func (s *Session) Close() {
	s.Client.close()
}

func (s *Session) SendRequest(index, begin, length int) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	msg := &message.Message{
		ID:      message.MsgRequest,
		Payload: payload,
	}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *Session) SendInterested() error {
	msg := &message.Message{ID: message.MsgInterested}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *Session) SendUninterested() error {
	msg := &message.Message{ID: message.MsgUninterested}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *Session) SendChoke() error {
	msg := &message.Message{ID: message.MsgChoke}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *Session) SendUnchoke() error {
	msg := &message.Message{ID: message.MsgUnchoke}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *Session) SendHave(index int) error {
	msg := &message.Message{ID: message.MsgHave}
	if _, err := s.Client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}
