package handshake

import (
	"fmt"
	"io"
)

type HandShake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func New(infoHash, peerID [20]byte) *HandShake {
	return &HandShake{
		Pstr:     "BitTorrent protocol",
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func (h *HandShake) Serialize() []byte {
	// 49 because 20 for infohash, 20 for peerID, 8 reserved, 1 for pstr len
	buf := make([]byte, len(h.Pstr)+49)

	buf[0] = byte(len(h.Pstr)) // length of protocol identifier
	curr := 1
	curr += copy(buf[curr:], []byte(h.Pstr))  // protocol identifier
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], h.InfoHash[:])   // info hash
	curr += copy(buf[curr:], h.PeerID[:])     // peer id

	return buf
}

func Read(r io.Reader) (*HandShake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	pstrLen := int(lengthBuf[0])
	if pstrLen == 0 {
		err := fmt.Errorf("pstr length cannot be 0")
		return nil, err
	}

	handshakeBuffer := make([]byte, pstrLen+48)
	if _, err := io.ReadFull(r, handshakeBuffer); err != nil {
		return nil, err
	}

	var infoHash [20]byte
	var peerID [20]byte

	copy(infoHash[:], handshakeBuffer[pstrLen+8:pstrLen+8+20])
	copy(peerID[:], handshakeBuffer[pstrLen+8+20:])

	h := &HandShake{
		Pstr:     string(handshakeBuffer[:pstrLen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return h, nil
}
