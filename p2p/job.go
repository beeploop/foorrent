package p2p

type pieceJob struct {
	index  int
	length int
	hash   [20]byte
}

type completedJob struct {
	index int
	buf   []byte
}

func calculatePieceSize(index int, pieceLength int) int {
	start, end := calculateBoundsOfPieace(index, pieceLength)
	return end - start
}

func calculateBoundsOfPieace(index int, pieceLength int) (int, int) {
	start := index * pieceLength
	end := start + pieceLength
	if end > pieceLength {
		end = pieceLength
	}

	return start, end
}
