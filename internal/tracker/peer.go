package tracker

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p *Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func parsePeersList(peersBytes []byte) ([]Peer, error) {
	peerSize := 6 // 6 bytes per peer | 4 for IP, 2 for port
	numOfPeers := len(peersBytes) / peerSize

	if len(peersBytes)%peerSize != 0 {
		err := fmt.Errorf("Malformed peers")
		return nil, err
	}

	peers := make([]Peer, numOfPeers)
	for i := range numOfPeers {
		offset := i * peerSize

		peer := Peer{
			IP:   net.IP(peersBytes[offset : offset+4]),
			Port: binary.BigEndian.Uint16(peersBytes[offset+4 : offset+6]),
		}
		peers[i] = peer
	}

	return peers, nil
}
