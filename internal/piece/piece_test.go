package piece

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPiece(t *testing.T) {
	t.Run("Test if piece is complete", func(t *testing.T) {
		p := &Piece{
			Blocks: make([]BlockState, 10),
		}

		assert.Equal(t, false, p.isComplete())
	})

	t.Run("Test reseting piece", func(t *testing.T) {
		length := 100
		numOfBlocks := int(math.Ceil(float64(length) / float64(MAX_BLOCK_SIZE)))

		initialData := make([]byte, length)
		initialBlocks := make([]BlockState, numOfBlocks)

		p := &Piece{
			Index:  0,
			Length: 100,
			Data:   initialData,
			Blocks: initialBlocks,
		}

		p.reset()

		assert.EqualValues(t, length, len(p.Data))
		assert.EqualValues(t, numOfBlocks, len(p.Blocks))
		assert.EqualValues(t, p.Data, initialData)
		assert.EqualValues(t, p.Blocks, initialBlocks)
	})
}
