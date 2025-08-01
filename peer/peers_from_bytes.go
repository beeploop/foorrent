package peers

import (
	"encoding/binary"
	"fmt"
	"net"
)

func PeersFromBytes(peersBytes []byte) ([]Peer, error) {
	peerSize := 6 // 6 bytes per peer | 4 for IP, 2 for port
	numOfPeers := len(peersBytes) / peerSize

	if len(peersBytes)%peerSize != 0 {
		err := fmt.Errorf("Malformed peers")
		return nil, err
	}

	peers := make([]Peer, numOfPeers)
	for i := 0; i < numOfPeers; i++ {
		offset := i * peerSize

		peer := Peer{
			IP:   net.IP(peersBytes[offset : offset+4]),
			Port: binary.BigEndian.Uint16(peersBytes[offset+4 : offset+6]),
		}
		peers[i] = peer
	}

	return peers, nil
}
