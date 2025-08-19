package peer

import (
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/message"
	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/piece"
	"github.com/beeploop/foorrent/internal/tracker"
)

type Session struct {
	Choked     bool
	Interested bool
	BitField   bitfield.BitField
	Client     *client
	Manager    piece.Manager
}

func New(peerID [20]byte, peer tracker.Peer, torrent metadata.Torrent, pm piece.Manager) (*Session, error) {
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
		Manager:    pm,
	}

	return session, nil
}

func (s *Session) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute * 2)
	defer ticker.Stop()
	defer s.Close()

	if err := s.SendInterested(); err != nil {
		fmt.Println("failed to send interested message")
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.SendKeepAlive(); err != nil {
				fmt.Println("error sending keep-alive: ", err.Error())
			}
		default:
			msg, err := message.Read(s.Client.conn)
			if err != nil {
				return
			}

			// Received keep-alive message
			if msg == nil {
				fmt.Println("received keep-alive")
				continue
			}

			switch msg.ID {
			case message.MsgChoke:
				fmt.Println("received choke")
				s.Choked = true

			case message.MsgUnchoke:
				fmt.Println("received unchoke")
				s.Choked = false

			case message.MsgHave:
				fmt.Println("received have")
				index := binary.BigEndian.Uint32(msg.Payload)
				s.BitField.SetPiece(int(index))

			case message.MsgBitfield:
				fmt.Println("received bitfield")
				s.BitField = msg.Payload

			case message.MsgPiece:
				fmt.Println("received piece")
				if len(msg.Payload) <= 8 {
					log.Println("malformed piece payload")
					return
				}

				index := int(binary.BigEndian.Uint32(msg.Payload[0:4]))
				begin := int(binary.BigEndian.Uint32(msg.Payload[4:8]))
				data := msg.Payload[8:]

				s.Manager.AddBlock(index, begin, data)

			case message.MsgInterested:
				fmt.Println("received interested")

			case message.MsgUninterested:
				fmt.Println("received uninterested")

			case message.MsgRequest:
				fmt.Println("received request")

			case message.MsgCancel:
				fmt.Println("received cancel")

			default:
				fmt.Println("received unknown mesage")
			}

			if !s.Choked {
				block, ok := s.Manager.NextRequest(s.BitField)
				if !ok {
					continue
				}

				fmt.Printf(">>> Requesting piece=%d, begin=%d\n", block.Index, block.Offset)
				if err := s.SendRequest(block.Index, block.Offset, block.Length); err != nil {
					log.Println("error sending request for a block")
					continue
				}
			}
		}
	}
}

func (s *Session) Close() {
	s.Client.close()
}

func (s *Session) SendKeepAlive() error {
	if _, err := s.Client.conn.Write([]byte{0, 0, 0, 0}); err != nil {
		return err
	}
	return nil
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
