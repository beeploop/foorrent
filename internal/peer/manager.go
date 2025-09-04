package peer

import (
	"context"
	"sync"

	"github.com/beeploop/foorrent/internal/metadata"
	"github.com/beeploop/foorrent/internal/piece"
	"github.com/beeploop/foorrent/internal/tracker"
)

type Manager struct {
	mu       sync.Mutex
	torrent  metadata.Torrent
	peerID   [20]byte
	peers    []tracker.Peer
	sessions []*session
	pm       *piece.Manager
}

func NewManager(torrent metadata.Torrent, peerID [20]byte, peers []tracker.Peer, pm *piece.Manager) *Manager {
	return &Manager{
		torrent:  torrent,
		peerID:   peerID,
		peers:    peers,
		sessions: make([]*session, 0),
		pm:       pm,
	}
}

func (m *Manager) Start(ctx context.Context) {
	for _, peer := range m.peers {
		go func(peer tracker.Peer) {
			session, err := newSession(m.peerID, peer, m.torrent, m.pm)
			if err != nil {
				return
			}

			m.mu.Lock()
			m.sessions = append(m.sessions, session)
			m.sessions[len(m.sessions)-1].start(ctx)
			m.mu.Unlock()
		}(peer)
	}
}

func (m *Manager) ActivePeers() (int, int) {
	return len(m.sessions), len(m.peers)
}
