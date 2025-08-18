package piece

import (
	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/metadata"
)

type Manager struct {
	torrent metadata.Torrent
	pieces  []*Piece
}

func NewManager(torrent metadata.Torrent) (*Manager, error) {
	pieces, err := initializePieces(torrent)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		torrent: torrent,
		pieces:  pieces,
	}
	return m, nil
}

// Returns index, offset/begin, length, ok (indicating if has piece to download)
func (m *Manager) NextRequest(peerBitField bitfield.BitField) (int, int, int, bool) {
	for i, piece := range m.pieces {
		piece.mu.Lock()

		if piece.isComplete() {
			piece.mu.Unlock()
			continue
		}

		if !peerBitField.HasPiece(i) {
			piece.mu.Unlock()
			continue
		}

		for blockIndex, have := range piece.Blocks {
			if !have {
				offset := blockIndex * MAX_BLOCK_SIZE
				length := MAX_BLOCK_SIZE

				// clamp last block to the length of data
				if offset+length > len(piece.Data) {
					length = len(piece.Data) - offset
				}

				piece.mu.Unlock()
				return i, offset, length, true
			}
		}

		piece.mu.Unlock()
	}

	return 0, 0, 0, false
}

func (m *Manager) AddBlock(index, offset int, data []byte) {
	piece := m.pieces[index]

	piece.mu.Lock()
	defer piece.mu.Unlock()

	blockIndex := offset / MAX_BLOCK_SIZE

	// Only copy the received block when it is not already received
	if piece.Blocks[blockIndex] == false {
		copy(piece.Data[offset:], data)
		piece.Blocks[blockIndex] = true
	}

	if piece.isComplete() {
		if err := piece.verify(); err != nil {
			piece.reset(m.torrent)
		}
	}
}
