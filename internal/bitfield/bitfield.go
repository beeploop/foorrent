package bitfield

type BitField []byte

func (bf BitField) HasPiece(index int) bool {
	byteIndex := index / 8 // Determine in which byte in []byte to look at
	bitOffset := index % 8 // Which position/index in the byte to look at

	if byteIndex < 0 || byteIndex >= len(bf) {
		return false
	}

	return bf[byteIndex]>>uint(7-bitOffset)&1 != 0
}

func (bf BitField) SetPiece(index int) {
	byteIndex := index / 8
	bitOffset := index % 8

	if byteIndex < 0 || byteIndex >= len(bf) {
		return
	}
	bf[byteIndex] |= 1 << uint(7-bitOffset)
}
