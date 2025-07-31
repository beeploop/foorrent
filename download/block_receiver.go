package download

import (
	"log"
)

func (dm *DownloadManager) handleReceivedBlocks() {
	for dm.downloaded < dm.TotalLength {
		block := <-dm.downloadChan

		dm.mu.Lock()
		inProgressPiece, exists := dm.inProgressDownloads[block.index]
		dm.mu.Unlock()

		if !exists {
			log.Println("Received block for unknown piece")
			continue
		}
		copied := copy(inProgressPiece.buf[block.begin:], block.data)
		inProgressPiece.downloaded += copied

		if !inProgressPiece.isComplete() {
			continue
		}

		if err := inProgressPiece.verifyHash(); err != nil {
			log.Println(err.Error())

			dm.pieceQueue <- &piece{
				index:  block.index,
				hash:   inProgressPiece.expectedHash,
				length: len(inProgressPiece.buf),
			}

			dm.mu.Lock()
			delete(dm.inProgressDownloads, block.index)
			dm.mu.Unlock()

			continue
		}

		// TODO: Let other peers know that we have this piece, emit `have` message

		dm.mu.Lock()
		dm.downloaded += len(inProgressPiece.buf)
		delete(dm.inProgressDownloads, block.index)
		dm.mu.Unlock()

		piece := &completedPiece{
			index: block.index,
		}
		copy(piece.buf[:], inProgressPiece.buf)

		dm.completedPieceChan <- piece
	}
}
