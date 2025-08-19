package piece

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/utils"
)

type BlockState int

const (
	MAX_BLOCK_SIZE = 16 * 1024 // 16KB

	Missing   BlockState = 0
	Requested BlockState = 1
	Done      BlockState = 2
)

type Block struct {
	Index  int
	Offset int
	Length int
}

type Piece struct {
	Index  int
	Length int
	Hash   [20]byte
	Data   []byte
	Blocks []BlockState
}

func initializePieces(torrent metadata.Torrent) ([]Piece, error) {
	hashes, err := torrent.PieceHashes()
	if err != nil {
		return nil, err
	}

	pieces := make([]Piece, 0)
	for i, hash := range hashes {
		length := torrent.Info.PieceLength

		// Last piece could be shorter than piece length, perform a double check
		if i == len(hashes)-1 {
			remainder := torrent.TotalSize() % torrent.Info.PieceLength
			if remainder != 0 {
				length = remainder
			}
		}

		numOfBlocks := int(math.Ceil(float64(length) / float64(MAX_BLOCK_SIZE)))

		pieces = append(pieces, Piece{
			Index:  i,
			Length: length,
			Hash:   hash,
			Data:   make([]byte, length),
			Blocks: make([]BlockState, numOfBlocks),
		})
	}

	return pieces, nil
}

func (p *Piece) verify() error {
	hash := sha1.Sum(p.Data)
	if !bytes.Equal(hash[:], p.Hash[:]) {
		err := fmt.Errorf("Piece failed verification check")
		return err
	}
	return nil
}

func (p *Piece) isComplete() bool {
	return utils.Every(p.Blocks, func(state BlockState) bool {
		return state == Done
	})
}

func (p *Piece) reset() {
	p.Data = make([]byte, len(p.Data))
	p.Blocks = make([]BlockState, len(p.Blocks))
}
