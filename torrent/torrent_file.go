package torrent

type TorrentFile struct {
	Announce     string // URL of the tracker
	Comment      string
	CreationDate int // Unix timestamp
	Info         TorrentFileInfo
}

type TorrentFileInfo struct {
	Length      int    // Size of the file in bytes
	Name        string // Suggested filename
	PieceLength int    // Number of bytes per piece
	Pieces      string
}
