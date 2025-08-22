package tracker

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"net"
	"net/url"
)

type UDP_PROTOCOL_EVENT uint32
type UDP_PROTOCOL_ACTION uint32

const (
	UDP_PROTOCOL_ID = 0x41727101980

	EVENT_NONE      UDP_PROTOCOL_EVENT = 0
	EVENT_COMPLETED UDP_PROTOCOL_EVENT = 1
	EVENT_STARTED   UDP_PROTOCOL_EVENT = 2
	EVENT_STOPPED   UDP_PROTOCOL_EVENT = 3

	ACTION_CONNECT  UDP_PROTOCOL_ACTION = 0
	ACTION_ANNOUNCE UDP_PROTOCOL_ACTION = 1
	ACTION_SCRAPE   UDP_PROTOCOL_ACTION = 2
	ACTION_ERROR    UDP_PROTOCOL_ACTION = 3
)

type UDPTracker struct{}

func NewUDPTracker() *UDPTracker {
	return &UDPTracker{}
}

func (t *UDPTracker) Request(input TrackerInput) (Result, error) {
	uri, err := url.Parse(input.Announce)
	if err != nil {
		return Result{}, err
	}

	addr, err := net.ResolveUDPAddr(uri.Scheme, uri.Host)
	if err != nil {
		return Result{}, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return Result{}, err
	}
	defer conn.Close()

	connectReq, connectTxID := t.buildConnectRequest()
	if _, err := conn.Write(connectReq); err != nil {
		return Result{}, err
	}

	connectRes := make([]byte, 1024)
	connectionID, err := t.readConnectResponse(conn, connectRes, connectTxID)
	if err != nil {
		return Result{}, err
	}

	announceReq, announceTxID := t.buildAnnounceRequest(connectionID, input)
	if _, err := conn.Write(announceReq); err != nil {
		return Result{}, err
	}

	announceRes := make([]byte, 1024)
	return t.readAnnounceResponse(conn, announceRes, announceTxID)
}

func (t *UDPTracker) buildConnectRequest() ([]byte, uint32) {
	txID := rand.Uint32()

	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[0:8], UDP_PROTOCOL_ID)         // protocol ID
	binary.BigEndian.PutUint32(buf[8:12], uint32(ACTION_CONNECT)) // action | 0 == connect
	binary.BigEndian.PutUint32(buf[12:16], txID)                  // transaction ID

	return buf, txID
}

// Returns the connection ID and error
func (t *UDPTracker) readConnectResponse(conn *net.UDPConn, resp []byte, txID uint32) (uint64, error) {
	if _, err := conn.Read(resp); err != nil {
		return 0, err
	}

	action := binary.BigEndian.Uint32(resp[0:4])
	connectionID := binary.BigEndian.Uint64(resp[8:16])

	if action != uint32(ACTION_CONNECT) {
		err := fmt.Errorf("Invalid action in response")
		return 0, err
	}

	if binary.BigEndian.Uint32(resp[4:8]) != txID {
		err := fmt.Errorf("Transaction ID mismatch")
		return 0, err
	}

	return connectionID, nil
}

func (t *UDPTracker) buildAnnounceRequest(connectionID uint64, input TrackerInput) ([]byte, uint32) {
	txID := rand.Uint32()

	buf := make([]byte, 98)
	binary.BigEndian.PutUint64(buf[0:8], uint64(connectionID))       // connection ID
	binary.BigEndian.PutUint32(buf[8:12], uint32(ACTION_ANNOUNCE))   // action | 1 == announce
	binary.BigEndian.PutUint32(buf[12:16], txID)                     // transaction ID
	copy(buf[16:36], input.InfoHash[:])                              // torrent info hash
	copy(buf[36:56], input.PeerID[:])                                // peer ID
	binary.BigEndian.PutUint64(buf[56:64], uint64(input.Downloaded)) // downloaded
	binary.BigEndian.PutUint64(buf[64:72], uint64(input.Left))       // left
	binary.BigEndian.PutUint64(buf[72:80], uint64(input.Uploaded))   // uploaded
	binary.BigEndian.PutUint32(buf[80:84], uint32(EVENT_STARTED))    // event | 2 == started
	binary.BigEndian.PutUint32(buf[84:88], 0)                        // IP
	binary.BigEndian.PutUint32(buf[88:92], rand.Uint32())            // randomy generated key
	binary.BigEndian.PutUint32(buf[92:96], 0xFFFFFFFF)               // num want = -1 (tracker decides)
	binary.BigEndian.PutUint16(buf[96:98], input.Port)               // port

	return buf, txID
}

func (t *UDPTracker) readAnnounceResponse(conn *net.UDPConn, resp []byte, txID uint32) (Result, error) {
	n, err := conn.Read(resp)
	if err != nil {
		return Result{}, err
	}

	action := binary.BigEndian.Uint32(resp[0:4])
	interval := binary.BigEndian.Uint32(resp[8:12])
	leechers := binary.BigEndian.Uint32(resp[12:16])
	seeders := binary.BigEndian.Uint32(resp[16:20])

	peers := make([]Peer, 0)
	for i := 0; i < len(resp[20:n]); i += 6 { // 4 bytes IP, 2 bytes PORT
		ip := net.IP(resp[20:n][i : i+4])
		port := binary.BigEndian.Uint16(resp[20:n][i+4 : i+6])

		peers = append(peers, Peer{
			IP:   ip,
			Port: port,
		})
	}

	if action != uint32(ACTION_ANNOUNCE) {
		errMessage := ""
		if action == uint32(ACTION_ERROR) {
			errMessage = string(resp[8:])
		}

		err := fmt.Errorf("Invalid action received from announce response: %d, error: %s\n", action, errMessage)
		return Result{}, err
	}

	if binary.BigEndian.Uint32(resp[4:8]) != txID {
		err := fmt.Errorf("Transaction ID mismatch")
		return Result{}, err
	}

	return Result{
		Interval: int(interval),
		Peers:    peers,
		Leechers: int(leechers),
		Seeders:  int(seeders),
	}, nil
}
