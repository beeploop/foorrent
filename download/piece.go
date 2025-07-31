package download

type piece struct {
	index  int
	hash   [20]byte
	length int
}

type completedPiece struct {
	index int
	buf   []byte
}

func calculatePieceSize(index, pieceLength, totalLength int) int {
	start, end := calculateBoundsOfPieace(index, pieceLength, totalLength)
	return end - start
}

func calculateBoundsOfPieace(index, pieceLength, totalLength int) (int, int) {
	start := index * pieceLength
	end := start + pieceLength
	if end > totalLength {
		end = totalLength
	}

	return start, end
}
