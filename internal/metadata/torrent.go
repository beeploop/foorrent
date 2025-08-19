package metadata

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"path/filepath"

	"github.com/jackpal/bencode-go"
)

type Torrent struct {
	Announce     string      `bencode:"announce"`
	Comment      string      `bencode:"comment"`
	CreationDate int         `bencode:"creation date"`
	Info         TorrentInfo `bencode:"info"`
}

type TorrentInfo struct {
	// Size of the file in bytes
	Length int `bencode:"length,omitempty"`
	// Suggested filename
	Name string `bencode:"name"`
	// Number of bytes per piece
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	// Files in a multi-file torrent
	Files []TorrentFile `bencode:"files,omitempty"`
}

type TorrentFile struct {
	Length int      `bencode:"length"`
	Path   []string `bencode:"path"`
}

func (t *Torrent) InfoHash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, t.Info); err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func (t *Torrent) PieceHashes() ([][20]byte, error) {
	data := []byte(t.Info.Pieces)
	hashSize := 20 // 20 byte for each hash | length of Sha-1

	if len(data)%hashSize != 0 {
		err := fmt.Errorf("Malformed pieces of length %d\n", len(data))
		return nil, err
	}

	numOfHashes := len(data) / hashSize
	hashes := make([][20]byte, numOfHashes)

	for i := range numOfHashes {
		start := i * hashSize
		end := start + hashSize
		copy(hashes[i][:], data[start:end])
	}

	return hashes, nil
}

func (t *Torrent) IsSingleFileMode() bool {
	return len(t.Info.Files) == 0
}

func (t *Torrent) FileMap() []FileEntry {
	list := make([]FileEntry, len(t.Info.Files))

	offset := 0
	for i, file := range t.Info.Files {
		list[i] = FileEntry{
			Path:   filepath.Join(file.Path...),
			Length: int64(file.Length),
			Offset: int64(offset),
		}

		offset += file.Length
	}

	return list
}

func (t *Torrent) TotalSize() int {
	if t.IsSingleFileMode() {
		return t.Info.Length
	}

	total := 0
	for _, file := range t.Info.Files {
		total += file.Length
	}
	return total
}
