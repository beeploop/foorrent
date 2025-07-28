package torrent

import (
	"os"

	"github.com/beeploop/foorrent/p2p"
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

func Download(torrent TorrentFile, outputPath string) error {
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

	// TODO: Handle re-contacting the tracker based on the Interval from response

	peers, err := peer.PeersFromBytes([]byte(response.Peers))
	if err != nil {
		return err
	}

	peer2peer := &p2p.PeerToPeer{
		Peers:       peers,
		PeerID:      peerID,
		InfoHash:    torrent.InfoHash,
		PieceHashes: torrent.PieceHashes,
	}
	if err := peer2peer.InitiateDownloadProcess(); err != nil {
		return err
	}

	return nil
}
