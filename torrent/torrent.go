package torrent

import (
	"fmt"
	"os"

	"github.com/beeploop/foorrent/peer"
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

func Download(torrent TorrentFile) error {
	peerID, err := peer.CreatePeerID()
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

	fmt.Println(response.Interval)
	fmt.Println(response.Peers)
	return nil
}
