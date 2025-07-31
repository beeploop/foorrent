package download

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/beeploop/foorrent/bitfield"
	"github.com/beeploop/foorrent/client"
	"github.com/beeploop/foorrent/peer"
)

type DownloadManager struct {
	mu                  sync.Mutex
	downloaded          int
	inProgressDownloads map[int]*inProgressPiece
	pieceQueue          chan *piece
	downloadChan        chan *block
	completedPieceChan  chan *completedPiece
	activePeers         map[string]peer.Peer

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
	dm.downloadChan = make(chan *block)
	dm.completedPieceChan = make(chan *completedPiece)
	dm.activePeers = make(map[string]peer.Peer)

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
			percent := dm.downloadPercent()
			log.Printf("[ %0.2f%% Downloaded ] [ Active Peers: %d ]\n", percent, len(dm.activePeers))
		case piece := <-dm.completedPieceChan:
			// TODO: Save the completed piece
			_ = piece

			percent := dm.downloadPercent()
			log.Printf("[ %0.2f%% Downloaded ] [ Active Peers: %d ]\n", percent, len(dm.activePeers))
		}
	}
}
