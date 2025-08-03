package download

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/peer"
	"github.com/beeploop/foorrent/utils"
)

const (
	OUTPUT_FILE_PERM = 0644
)

type DownloadManager struct {
	mu                   sync.Mutex
	downloaded           int
	inProgressDownloads  map[int]*inProgressPiece
	pieceQueue           chan *piece
	downloadChan         chan *block
	downloadCompleteChan chan struct{}
	completedPieceChan   chan *completedPiece
	activePeers          map[string]peers.Peer

	OutputPath  string
	FileName    string
	PieceHashes [][20]byte
	PieceLength int
	TotalPieces int
	TotalLength int
	InfoHash    [20]byte
	PeerID      [20]byte
	Peers       []peers.Peer
}

func (dm *DownloadManager) setQueue() error {
	if dm.PieceHashes == nil {
		err := fmt.Errorf("piece hashes not passed to download manager")
		return err
	}

	if dm.pieceQueue == nil {
		dm.pieceQueue = make(chan *piece, len(dm.PieceHashes))
	}

	for index, hash := range dm.PieceHashes {
		length := calculatePieceSize(index, dm.PieceLength, dm.TotalLength)
		dm.pieceQueue <- &piece{index, hash, length}
	}

	return nil
}

func (dm *DownloadManager) downloadPercent() float64 {
	if dm.TotalLength == 0 {
		return 0.0
	}
	return float64(dm.downloaded) / float64(dm.TotalLength) * 100
}

func (dm *DownloadManager) getOutputPath() string {
	return filepath.Join(dm.OutputPath, dm.FileName)
}

func (dm *DownloadManager) Start() error {
	if err := dm.setQueue(); err != nil {
		return err
	}

	dm.inProgressDownloads = make(map[int]*inProgressPiece)
	dm.activePeers = make(map[string]peers.Peer)

	dm.downloadChan = make(chan *block)
	dm.completedPieceChan = make(chan *completedPiece)
	dm.downloadCompleteChan = make(chan struct{})

	for _, peer := range dm.Peers {
		go func() {
			c, err := client.New(peer, dm.PeerID, dm.InfoHash)
			if err != nil {
				log.Println("Client for peer failed: ", err.Error())
				return
			}
			dm.mu.Lock()
			dm.activePeers[peer.String()] = peer
			dm.mu.Unlock()

			c.SendInterested()

			session := &peerSession{
				choked: true,
				bf:     make(bitfield.BitField, len(dm.PieceHashes)),
			}
			go dm.handleInboundMessages(c, session)
			go dm.requestPieces(c, session)
			go dm.peerResponder(c)
		}()
	}

	go dm.handleReceivedBlocks()

	file, err := os.OpenFile(dm.getOutputPath(), os.O_CREATE|os.O_WRONLY, OUTPUT_FILE_PERM)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := file.Truncate(int64(dm.TotalLength)); err != nil {
		return err
	}

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			goroutines := runtime.NumGoroutine() - 1
			percent := dm.downloadPercent()
			downloaded := utils.ToMegabytes(float64(dm.downloaded))
			total := utils.ToMegabytes(float64(dm.TotalLength))
			log.Printf(
				"[ %0.2fMB/%0.2fMB ] [ %0.2f%% Downloaded ] [ Active Peers: %d ] [ Goroutines: %d ]\n",
				downloaded, total, percent, len(dm.activePeers), goroutines,
			)

		case piece := <-dm.completedPieceChan:
			dm.mu.Lock()
			downloaded := dm.downloaded
			totalLength := dm.TotalLength
			dm.mu.Unlock()

			if downloaded == totalLength {
				dm.downloadCompleteChan <- struct{}{}
			}

			offset := int64(piece.index * dm.PieceLength)
			bytesWritten, err := file.WriteAt(piece.buf, offset)
			if err != nil {
				err := fmt.Errorf("Error occurred while writing to file, aborted")
				return err
			}

			if bytesWritten != len(piece.buf) {
				err := fmt.Errorf("Bytes written to file doesn't match piece buf length, aborted")
				return err
			}
		}
	}
}
