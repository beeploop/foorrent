package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"

	"github.com/jackpal/bencode-go"
)

type TorrentFile struct {
	Announce     string // URL of the tracker
	Comment      string
	CreationDate int // Unix timestamp
	Info         TorrentFileInfo
	InfoHash     [20]byte
	PieceHashes  [][20]byte // List of split pieces of the pieces string
}

type TorrentFileInfo struct {
	Length      int    // Size of the file in bytes
	Name        string // Suggested filename
	PieceLength int    // Number of bytes per piece
	Pieces      string
}

func torrentFileFromBencode(data bencodeContent) (TorrentFile, error) {
	info := TorrentFileInfo{
		Length:      data.Info.Length,
		Name:        data.Info.Name,
		PieceLength: data.Info.PieceLength,
		Pieces:      data.Info.Pieces,
	}

	infoHash, err := createHashInfo(data.Info)
	if err != nil {
		return TorrentFile{}, err
	}

	pieceHashes, err := splitPiecesString(info.Pieces)
	if err != nil {
		return TorrentFile{}, err
	}

	t := TorrentFile{
		Announce:     data.Announce,
		Comment:      data.Comment,
		CreationDate: data.CreationDate,
		Info:         info,
		InfoHash:     infoHash,
		PieceHashes:  pieceHashes,
	}

	return t, nil
}

func createHashInfo(info bencodeInfo) ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, info); err != nil {
		return [20]byte{}, err
	}

	hash := sha1.Sum(buf.Bytes())
	return hash, nil
}

func splitPiecesString(pieces string) ([][20]byte, error) {
	data := []byte(pieces)
	hashSize := 20 // 20 byte for each hash | length of Sha-1

	if len(data)%hashSize != 0 {
		err := fmt.Errorf("Malformed pieces of length %d\n", len(data))
		return nil, err
	}

	numOfHashes := len(data) / hashSize
	hashes := make([][20]byte, numOfHashes)

	for i := 0; i < numOfHashes; i++ {
		start := i * hashSize
		end := start + hashSize
		copy(hashes[i][:], data[start:end])
	}

	return hashes, nil
}
