package torrent

import (
	"os"

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
