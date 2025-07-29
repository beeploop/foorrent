package p2p

import (
	"log"
	"time"

	"github.com/beeploop/foorrent/client"
)

type DownloadWorker struct {
	c                *client.Client
	jobQueue         chan *pieceJob
	completedJobChan chan *completedJob
	peerState        *PeerState
}

func (d *DownloadWorker) start() {
	d.c.Conn.SetDeadline(time.Now().Add(time.Second * 5))
	defer d.c.Conn.SetDeadline(time.Time{})

	if err := d.c.SendInterested(); err != nil {
		log.Println("Failed to send interested message")
		return
	}

	for job := range d.jobQueue {
		if !d.peerState.BitField.HasPiece(job.index) {
			d.jobQueue <- job // Queue job again if peer doesn't have the piece
			continue
		}

		buf, err := d.attemptDownload()
		if err != nil {
			log.Printf("Piece #%d download failed, requeueing piece", job.index)
			d.jobQueue <- job // Queue job again if download fails
			continue
		}

		if err := d.verifyPiece(buf); err != nil {
			log.Printf("Piece #%d failed verification check\n", job.index)
			d.jobQueue <- job // Queue job again if hash verification failed
			continue
		}

		d.c.SendHave(job.index)

		d.completedJobChan <- &completedJob{
			index: job.index,
			buf:   buf,
		}
	}
}

func (d *DownloadWorker) attemptDownload() ([]byte, error) {
	// TODO: Handle download attempt
	return make([]byte, 0), nil
}

func (d *DownloadWorker) verifyPiece(buf []byte) error {
	// TODO: Handle piece hash verification
	_ = buf
	return nil
}
