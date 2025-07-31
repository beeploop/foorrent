package download

import (
	"encoding/binary"
	"sync"

	"github.com/beeploop/foorrent/bitfield"
)

type peerSession struct {
	mu     sync.Mutex
	choked bool
	bf     bitfield.BitField
}

func (s *peerSession) choke() {
	s.mu.Lock()
	s.choked = true
	s.mu.Unlock()
}

func (s *peerSession) unchoke() {
	s.mu.Lock()
	s.choked = false
	s.mu.Unlock()
}

func (s *peerSession) setBitfield(bf []byte) {
	s.mu.Lock()
	s.bf = bf
	s.mu.Unlock()
}

func (s *peerSession) putPiece(payload []byte) {
	index := binary.BigEndian.Uint32(payload)
	s.mu.Lock()
	s.bf.SetPiece(int(index))
	s.mu.Unlock()
}
