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
