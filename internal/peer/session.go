package peer

import (
	"context"
	"encoding/binary"
	"log"
	"time"

	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/message"
	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/piece"
	"github.com/beeploop/foorrent/internal/tracker"
)

type session struct {
	peer       tracker.Peer
	choked     bool
	interested bool
	bitField   bitfield.BitField
	client     *client
	pm         *piece.Manager
}

func newSession(peerID [20]byte, peer tracker.Peer, torrent metadata.Torrent, pm *piece.Manager) (*session, error) {
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

	session := &session{
		peer:       peer,
		choked:     true,
		interested: false,
		bitField:   make(bitfield.BitField, len(hashes)),
		client:     c,
		pm:         pm,
	}

	return session, nil
}

func (s *session) start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 2)
	defer ticker.Stop()
	defer s.Close()

	if err := s.SendInterested(); err != nil {
		log.Println("failed to send interested message")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.SendKeepAlive(); err != nil {
			}
		default:
			msg, err := message.Read(s.client.conn)
			if err != nil {
				return
			}

			// Received keep-alive message
			if msg == nil {
				continue
			}

			switch msg.ID {
			case message.MsgChoke:
				s.choked = true

			case message.MsgUnchoke:
				s.choked = false

			case message.MsgHave:
				index := binary.BigEndian.Uint32(msg.Payload)
				s.bitField.SetPiece(int(index))

			case message.MsgBitfield:
				s.bitField = msg.Payload

			case message.MsgPiece:
				if len(msg.Payload) <= 8 {
					log.Println("malformed piece payload")
					return
				}

				index := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
				begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
				data := msg.Payload[8:]

				complete := s.pm.AddBlock(index, begin, data)
				if complete {
					if err := s.SendHave(index); err != nil {
						log.Println("failed to send have message")
						continue
					} else {
						log.Println("sent have message")
					}
				}

			case message.MsgInterested:
				log.Println("received interested message")
				// TODO: Implement

			case message.MsgUninterested:
				log.Println("received uninterested message")
				// TODO: Implement

			case message.MsgRequest:
				log.Println("received a request message")
				// TODO: Implement

			case message.MsgCancel:
				log.Println("received a cancel message")
				// TODO: Implement

			default:
				log.Println("received unknown mesage")
			}

			if !s.choked {
				block, ok := s.pm.NextRequest(s.bitField)
				if !ok {
					continue
				}

				if err := s.SendRequest(block.Index, block.Offset, block.Length); err != nil {
					log.Println("error sending request for a block")
					continue
				}
			}
		}
	}
}

func (s *session) Close() {
	s.client.close()
}

func (s *session) SendKeepAlive() error {
	if _, err := s.client.conn.Write([]byte{0, 0, 0, 0}); err != nil {
		return err
	}
	return nil
}

func (s *session) SendRequest(index, begin, length int) error {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(begin))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))

	msg := &message.Message{
		ID:      message.MsgRequest,
		Payload: payload,
	}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *session) SendInterested() error {
	msg := &message.Message{ID: message.MsgInterested}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *session) SendUninterested() error {
	msg := &message.Message{ID: message.MsgUninterested}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *session) SendChoke() error {
	msg := &message.Message{ID: message.MsgChoke}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *session) SendUnchoke() error {
	msg := &message.Message{ID: message.MsgUnchoke}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}

func (s *session) SendHave(index int) error {
	msg := &message.Message{ID: message.MsgHave}
	if _, err := s.client.conn.Write(msg.Serialize()); err != nil {
		return err
	}
	return nil
}
