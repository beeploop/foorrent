package download

import (
	"encoding/binary"
	"fmt"
)

type block struct {
	index int
	begin int
	data  []byte
}

func parsePiece(payload []byte) (*block, error) {
	if len(payload) <= 8 {
		err := fmt.Errorf("malformed piece payload")
		return nil, err
	}

	block := &block{
		index: int(binary.BigEndian.Uint32(payload[0:4])),
		begin: int(binary.BigEndian.Uint32(payload[4:8])),
		data:  payload[8:],
	}

	return block, nil
}
