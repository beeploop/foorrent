package metadata

import (
	"os"

	"github.com/jackpal/bencode-go"
)

func Read(path string) (Torrent, error) {
	file, err := os.Open(path)
	if err != nil {
		return Torrent{}, err
	}
	defer file.Close()

	var content Torrent
	if err := bencode.Unmarshal(file, &content); err != nil {
		return Torrent{}, err
	}

	return content, nil
}
