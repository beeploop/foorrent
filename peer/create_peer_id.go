package peers

import "crypto/rand"

func CreatePeerID() ([20]byte, error) {
	var peerID [20]byte
	_, err := rand.Read(peerID[:])
	if err != nil {
		return peerID, err
	}

	return peerID, nil
}
