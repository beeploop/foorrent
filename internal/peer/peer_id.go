package peer

import "crypto/rand"

func GeneratePeerID() ([20]byte, error) {
	var peerID [20]byte
	if _, err := rand.Read(peerID[:]); err != nil {
		return peerID, err
	}

	return peerID, nil
}
