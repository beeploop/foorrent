package piece

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPiece(t *testing.T) {
	t.Run("Test if piece is complete", func(t *testing.T) {
		p := &Piece{
			Blocks: make([]bool, 10),
		}

		assert.Equal(t, false, p.isComplete())
	})
}
