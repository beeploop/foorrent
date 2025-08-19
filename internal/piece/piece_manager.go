package piece

import (
	"sync"

	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/storage"
)

type PieceManager struct {
	mu      sync.Mutex
	torrent metadata.Torrent
	pieces  []Piece
	storage storage.Storage
}

func NewManager(torrent metadata.Torrent, storage storage.Storage) (*PieceManager, error) {
	pieces, err := initializePieces(torrent)
	if err != nil {
		return nil, err
	}

	m := &PieceManager{
		torrent: torrent,
		pieces:  pieces,
		storage: storage,
	}
	return m, nil
}

// Returns block, ok (indicating if has block to download)
func (m *PieceManager) NextRequest(peerBitField bitfield.BitField) (Block, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for index, piece := range m.pieces {
		if !peerBitField.HasPiece(index) {
			continue
		}

		if piece.isComplete() {
			continue
		}

		for blockIndex, state := range piece.Blocks {
			if state == Missing {
				offset := blockIndex * MAX_BLOCK_SIZE
				length := MAX_BLOCK_SIZE

				// clamp last block to the length of data
				if offset+length > len(piece.Data) {
					length = len(piece.Data) - offset
				}

				piece.Blocks[blockIndex] = Requested
				block := Block{Index: index, Offset: offset, Length: length}
				return block, true
			}
		}
	}

	return Block{}, false
}

func (m *PieceManager) AddBlock(index, offset int, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	piece := m.pieces[index]
	blockIndex := offset / MAX_BLOCK_SIZE

	// Only copy the received block when it is not already received
	if piece.Blocks[blockIndex] == Missing {
		copy(piece.Data[offset:], data)
		piece.Blocks[blockIndex] = Done
	}

	if piece.isComplete() {
		if err := piece.verify(); err != nil {
			piece.reset()
		}

		m.storage.WritePiece(piece.Index, piece.Length, piece.Data)
	}
}
