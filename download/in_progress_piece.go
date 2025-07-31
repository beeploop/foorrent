package download

import (
	"bytes"
	"crypto/sha1"
	"fmt"
)

type inProgressPiece struct {
	buf          []byte
	downloaded   int
	expectedHash [20]byte
}

func (p *inProgressPiece) isComplete() bool {
	return p.downloaded == len(p.buf)
}

func (p *inProgressPiece) verifyHash() error {
	hash := sha1.Sum(p.buf)
	if !bytes.Equal(hash[:], p.expectedHash[:]) {
		err := fmt.Errorf("Piece failed verification check")
		return err
	}
	return nil
}
