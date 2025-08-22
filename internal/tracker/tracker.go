package tracker

import (
	"fmt"
	"net/url"
)

const DEFAULT_PORT uint16 = 6881

type Tracker interface {
	Request(TrackerInput) (Result, error)
}

type Result struct {
	Interval int
	Peers    []Peer
	Leechers int
	Seeders  int
}

type TrackerInput struct {
	Announce   string
	InfoHash   [20]byte
	PeerID     [20]byte
	Port       uint16
	Left       int
	Downloaded int
	Uploaded   int
}

func TrackerFactory(input TrackerInput) (Tracker, error) {
	u, err := url.Parse(input.Announce)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http":
		return NewHTTPTracker(), nil
	case "udp":
		return NewUDPTracker(), nil
	default:
		err := fmt.Errorf("unknown announce scheme: %s\n", u.Scheme)
		return nil, err
	}
}
