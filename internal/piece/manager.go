package piece

import (
	"sync"

	"github.com/beeploop/foorrent/internal/bitfield"
	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/storage"
	"github.com/beeploop/foorrent/internal/utils"
)

type Manager struct {
	mu      sync.Mutex
	torrent metadata.Torrent
	pieces  []Piece
	storage storage.Storage
}

func NewManager(torrent metadata.Torrent, storage storage.Storage) (*Manager, error) {
	pieces, err := initializePieces(torrent)
	if err != nil {
		return nil, err
	}

	m := &Manager{
		torrent: torrent,
		pieces:  pieces,
		storage: storage,
	}
	return m, nil
}

// Returns block, ok (indicating if has block to download)
func (m *Manager) NextRequest(peerBitField bitfield.BitField) (Block, bool) {
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

func (m *Manager) AddBlock(index, offset int, data []byte) {
	m.mu.Lock()
	defer m.mu.Unlock()

	piece := m.pieces[index]
	blockIndex := offset / MAX_BLOCK_SIZE

	// Only copy the received block when it is not already received
	if piece.Blocks[blockIndex] != Done {
		copy(piece.Data[offset:], data)
		piece.Blocks[blockIndex] = Done
	}

	if piece.isComplete() {
		if err := piece.verify(); err != nil {
			piece.resetData()
		}

		if err := m.storage.WritePiece(piece.Index, piece.Length, piece.Data); err != nil {
			piece.resetData()
		}
		piece.finalizeAndFree()
	}
}

// Returns number of downloaded pieces and total pieces
func (m *Manager) Downloaded() (int, int) {
	counter := 0
	for _, p := range m.pieces {
		if p.done() {
			counter++
		}
	}

	return counter, len(m.pieces)
}

func (m *Manager) Done() bool {
	return utils.Every(m.pieces, func(p Piece) bool {
		return p.done()
	})
}
