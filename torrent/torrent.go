package torrent

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/beeploop/foorrent/download"
	"github.com/beeploop/foorrent/peer"
	"github.com/beeploop/foorrent/utils"
	"github.com/jackpal/bencode-go"
)

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, err
	}
	defer file.Close()

	var content bencodeContent
	if err := bencode.Unmarshal(file, &content); err != nil {
		return TorrentFile{}, err
	}

	return torrentFileFromBencode(content)
}

func Download(torrent TorrentFile, outputPath string) error {
	if len(torrent.Info.Files) > 0 {
		err := fmt.Errorf("Does not yet support multi-file torrent")
		return err
	}

	savepath, err := filepath.Abs(outputPath)
	if err != nil {
		return err
	}

	peerID, err := peers.CreatePeerID()
	if err != nil {
		return err
	}

	trackerURL, err := buildTrackerURL(peerID, torrent)
	if err != nil {
		return err
	}

	response, err := contactTracker(trackerURL)
	if err != nil {
		return err
	}

	// TODO: Handle re-contacting the tracker based on the Interval from response

	peers, err := peers.PeersFromBytes([]byte(response.Peers))
	if err != nil {
		return err
	}
	log.Printf("NUMBER OF SEEDERS: %d\n", len(peers))
	log.Printf("TOTAL FILE SIZE: %0.0fMB\n", utils.ToMegabytes(float64(torrent.Info.Length)))

	mgr := &download.DownloadManager{
		OutputPath:  savepath,
		FileName:    torrent.Info.Name,
		PieceHashes: torrent.PieceHashes,
		PieceLength: torrent.Info.PieceLength,
		TotalPieces: len(torrent.PieceHashes),
		TotalLength: torrent.Info.Length,
		InfoHash:    torrent.InfoHash,
		PeerID:      peerID,
		Peers:       peers,
	}
	if err := mgr.Start(); err != nil {
		return err
	}

	return nil
}
