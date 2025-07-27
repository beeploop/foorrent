package torrent

type bencodeInfo struct {
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
}

type bencodeContent struct {
	Announce     string      `bencode:"announce"`
	Comment      string      `bencode:"comment"`
	CreationDate int         `bencode:"creation date"`
	Info         bencodeInfo `bencode:"info"`
}

func (b *bencodeContent) ToTorrentFile() (TorrentFile, error) {
	t := TorrentFile{
		Announce:     b.Announce,
		Comment:      b.Comment,
		CreationDate: b.CreationDate,
		Info: TorrentFileInfo{
			Length:      b.Info.Length,
			Name:        b.Info.Name,
			PieceLength: b.Info.PieceLength,
			Pieces:      b.Info.Pieces,
		},
	}

	return t, nil
}
