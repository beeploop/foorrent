package download

import (
	"log"

	"github.com/beeploop/foorrent/client"
)

const MAX_BLOCK_SIZE = 16 * 1024 // 16KB

func (dm *DownloadManager) requestPieces(c *client.Client, session *peerSession) {
	for piece := range dm.pieceQueue {
		if !session.bf.HasPiece(piece.index) {
			dm.pieceQueue <- piece
			continue
		}

		if session.choked {
			log.Println("could not request for blocks because choked")
			dm.pieceQueue <- piece
			continue
		}

		dm.mu.Lock()
		dm.inProgressDownloads[piece.index] = &inProgressPiece{
			buf:          make([]byte, piece.length),
			downloaded:   0,
			expectedHash: piece.hash,
		}
		dm.mu.Unlock()

		log.Println("Requesting blocks for piece: ", piece.index)
		for begin := 0; begin < piece.length; begin += MAX_BLOCK_SIZE {
			length := min(MAX_BLOCK_SIZE, piece.length-begin)
			if err := c.SendRequest(piece.index, begin, length); err != nil {
				log.Printf("Request for block failed. index: %d, begin: %d, length: %d\n", piece.index, begin, length)
				dm.pieceQueue <- piece
				continue
			}
		}
	}
}
