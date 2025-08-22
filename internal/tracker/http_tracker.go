package tracker

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type response struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
	Leechers int    `bencode:"leechers"`
	Seeder   int    `bencode:"seeders"`
}

type HTTPTracker struct{}

func NewHTTPTracker() *HTTPTracker {
	return &HTTPTracker{}
}

func (t *HTTPTracker) Request(input TrackerInput) (Result, error) {
	uri, err := t.constructURL(input)
	if err != nil {
		return Result{}, err
	}

	resp, err := t.contactTracker(uri)
	if err != nil {
		return Result{}, err
	}

	peers, err := parsePeersList([]byte(resp.Peers))
	if err != nil {
		return Result{}, err
	}

	return Result{
		Interval: resp.Interval,
		Peers:    peers,
	}, nil
}

func (t *HTTPTracker) contactTracker(url string) (response, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Get(url)
	if err != nil {
		return response{}, err
	}
	defer resp.Body.Close()

	var content response
	if err := bencode.Unmarshal(resp.Body, &content); err != nil {
		return response{}, err
	}

	return content, nil
}

func (t *HTTPTracker) constructURL(input TrackerInput) (string, error) {
	baseURL, err := url.Parse(input.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("info_hash", string(input.InfoHash[:]))
	params.Add("peer_id", string(input.PeerID[:]))
	params.Add("port", strconv.Itoa(int(DEFAULT_PORT)))
	params.Add("uploaded", strconv.Itoa(input.Uploaded))
	params.Add("downloaded", strconv.Itoa(input.Downloaded))
	params.Add("compact", "1")
	params.Add("left", strconv.Itoa(input.Left))

	baseURL.RawQuery = params.Encode()

	return baseURL.String(), nil
}
