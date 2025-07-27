package torrent

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

const port uint16 = 6881

type trackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func buildTrackerURL(peerID [20]byte, torrent TorrentFile) (string, error) {
	baseURL, err := url.Parse(torrent.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("info_hash", string(torrent.InfoHash[:]))
	params.Add("peer_id", string(peerID[:]))
	params.Add("port", strconv.Itoa(int(port)))
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(torrent.Info.Length))

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}

func contactTracker(url string) (trackerResponse, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Get(url)
	if err != nil {
		return trackerResponse{}, err
	}
	defer res.Body.Close()

	var trackerResp trackerResponse
	if err := bencode.Unmarshal(res.Body, &trackerResp); err != nil {
		return trackerResponse{}, err
	}

	return trackerResp, nil
}
