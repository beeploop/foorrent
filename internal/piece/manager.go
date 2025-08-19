package piece

import "github.com/beeploop/foorrent/internal/bitfield"

type Manager interface {
	NextRequest(peerBitField bitfield.BitField) (Block, bool)
	AddBlock(index, offset int, data []byte)
}
