package download

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/peer"
)

type DownloadManager struct {
	mu                   sync.Mutex
	downloaded           int
	inProgressDownloads  map[int]*inProgressPiece
	pieceQueue           chan *piece
	downloadChan         chan *block
	downloadCompleteChan chan struct{}
	completedPieceChan   chan *completedPiece
	activePeers          map[string]peer.Peer

	PieceHashes [][20]byte
	PieceLength int
	TotalPieces int
	TotalLength int
	InfoHash    [20]byte
	PeerID      [20]byte
	Peers       []peer.Peer
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

func (dm *DownloadManager) Start() error {
	if err := dm.setQueue(); err != nil {
		return err
	}

	dm.inProgressDownloads = make(map[int]*inProgressPiece)
	dm.activePeers = make(map[string]peer.Peer)

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

	ticker := time.NewTicker(time.Second * 1)
	for {
		select {
		case <-ticker.C:
			goroutines := runtime.NumGoroutine() - 1
			percent := dm.downloadPercent()
			log.Printf("[ %0.2f%% Downloaded ] [ Active Peers: %d ] [ Goroutines: %d ]\n", percent, len(dm.activePeers), goroutines)

		case piece := <-dm.completedPieceChan:
			dm.mu.Lock()
			downloaded := dm.downloaded
			totalLength := dm.TotalLength
			dm.mu.Unlock()

			if downloaded == totalLength {
				dm.downloadCompleteChan <- struct{}{}
			}

			// TODO: Save the completed piece
			_ = piece
		}
	}
}
