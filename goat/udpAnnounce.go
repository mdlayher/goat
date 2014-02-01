package goat

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net/url"
	"strconv"
)

// udpAnnounceRequest represents a tracker announce in the UDP format
type udpAnnounceRequest struct {
	ConnID     uint64
	Action     uint32
	TransID    uint32
	InfoHash   []byte
	PeerID     []byte
	Downloaded uint64
	Left       uint64
	Uploaded   uint64
	Event      uint32
	IP         uint32
	Key        uint32
	Numwant    uint32
	Port       uint16
}

// FromBytes creates a udpAnnounceRequest from a packed byte array
func (u udpAnnounceRequest) FromBytes(buf []byte) (p udpAnnounceRequest, err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			p = udpAnnounceRequest{}
			err = errors.New("failed to create udpAnnounceRequest from bytes")
		}
	}()

	// ConnID (uint64)
	u.ConnID = binary.BigEndian.Uint64(buf[0:8])

	// Action (uint32) (Announce = 1)
	u.Action = binary.BigEndian.Uint32(buf[8:12])
	if u.Action != uint32(1) {
		return udpAnnounceRequest{}, errors.New("invalid action for udpAnnounceRequest")
	}

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[12:16])

	// InfoHash (20 bytes)
	u.InfoHash = buf[16:36]
	if len(u.InfoHash) != 20 {
		return udpAnnounceRequest{}, errors.New("info_hash must be exactly 20 bytes")
	}

	// PeerID (20 bytes)
	u.PeerID = buf[36:56]
	if len(u.PeerID) != 20 {
		return udpAnnounceRequest{}, errors.New("peer_id must be exactly 20 bytes")
	}

	// Downloaded (uint64)
	u.Downloaded = binary.BigEndian.Uint64(buf[56:64])

	// Left (uint64)
	u.Left = binary.BigEndian.Uint64(buf[64:72])

	// Uploaded (uint64)
	u.Uploaded = binary.BigEndian.Uint64(buf[72:80])

	// Event (uint32)
	u.Event = binary.BigEndian.Uint32(buf[80:84])

	// IP (uint32)
	u.IP = binary.BigEndian.Uint32(buf[84:88])

	// Key (uint32)
	u.Key = binary.BigEndian.Uint32(buf[88:92])

	// Numwant (uint32)
	numwant := binary.BigEndian.Uint32(buf[92:96])
	// If numwant is uint32 max, use protocol default of 50
	if numwant == uint32(4294967295) {
		numwant = 50
	}
	u.Numwant = numwant

	// Port (uint16)
	u.Port = binary.BigEndian.Uint16(buf[96:98])

	return u, nil
}

// ToBytes creates a packed byte array from a udpAnnounceRequest
func (u udpAnnounceRequest) ToBytes() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// ConnID (uint64)
	if err := binary.Write(res, binary.BigEndian, u.ConnID); err != nil {
		return nil, err
	}

	// Action (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// InfoHash (20 bytes)
	if len(u.InfoHash) != 20 {
		return nil, errors.New("info_hash must be exactly 20 bytes")
	}

	if err := binary.Write(res, binary.BigEndian, u.InfoHash); err != nil {
		return nil, err
	}

	// PeerID (20 bytes)
	if len(u.PeerID) != 20 {
		return nil, errors.New("peer_id must be exactly 20 bytes")
	}

	if err := binary.Write(res, binary.BigEndian, u.PeerID); err != nil {
		return nil, err
	}

	// Downloaded (uint64)
	if err := binary.Write(res, binary.BigEndian, u.Downloaded); err != nil {
		return nil, err
	}

	// Left (uint64)
	if err := binary.Write(res, binary.BigEndian, u.Left); err != nil {
		return nil, err
	}

	// Uploaded (uint64)
	if err := binary.Write(res, binary.BigEndian, u.Uploaded); err != nil {
		return nil, err
	}

	// Event (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Event); err != nil {
		return nil, err
	}

	// IP (uint32)
	if err := binary.Write(res, binary.BigEndian, u.IP); err != nil {
		return nil, err
	}

	// Key (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Key); err != nil {
		return nil, err
	}

	// Numwant (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Numwant); err != nil {
		return nil, err
	}

	// Port (uint16)
	if err := binary.Write(res, binary.BigEndian, u.Port); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

// ToValues creates a url.Values struct from a udpAnnounceRequest
func (u udpAnnounceRequest) ToValues() url.Values {
	// Initialize query map
	query := url.Values{}
	query.Set("udp", "1")

	// Copy all fields into query map
	query.Set("info_hash", string(u.InfoHash))

	// Integer fields
	query.Set("downloaded", strconv.FormatUint(u.Downloaded, 10))
	query.Set("left", strconv.FormatUint(u.Left, 10))
	query.Set("uploaded", strconv.FormatUint(u.Uploaded, 10))

	// Event, converted to actual string
	switch u.Event {
	case 0:
		query.Set("event", "")
	case 1:
		query.Set("event", "completed")
	case 2:
		query.Set("event", "started")
	case 3:
		query.Set("event", "stopped")
	}

	// IP
	query.Set("ip", strconv.FormatUint(uint64(u.IP), 10))

	// Key
	query.Set("key", strconv.FormatUint(uint64(u.Key), 10))

	// Numwant
	query.Set("numwant", strconv.FormatUint(uint64(u.Numwant), 10))

	// Port
	query.Set("port", strconv.FormatUint(uint64(u.Port), 10))

	// Return final query map
	return query
}

// udpAnnounceResponse represents a tracker announce response in the UDP format
type udpAnnounceResponse struct {
	Action   uint32
	TransID  uint32
	Interval uint32
	Leechers uint32
	Seeders  uint32
	PeerList []compactPeer
}

// FromBytes creates a udpAnnounceResponse from a packed byte array
func (u udpAnnounceResponse) FromBytes(buf []byte) (p udpAnnounceResponse, err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			p = udpAnnounceResponse{}
			err = errors.New("failed to create udpAnnounceResponse from bytes")
		}
	}()

	// Action (uint32) (Announce = 1)
	u.Action = binary.BigEndian.Uint32(buf[0:4])
	if u.Action != uint32(1) {
		return udpAnnounceResponse{}, errors.New("invalid action for udpAnnounceResponse")
	}

	// Transaction ID
	u.TransID = binary.BigEndian.Uint32(buf[4:8])

	// Interval
	u.Interval = binary.BigEndian.Uint32(buf[8:12])

	// Leechers
	u.Leechers = binary.BigEndian.Uint32(buf[12:16])

	// Seeders
	u.Seeders = binary.BigEndian.Uint32(buf[16:20])

	// Peer List
	u.PeerList = make([]compactPeer, 0)

	// Iterate peers buffer
	i := 20
	for {
		// Validate that we are not seeking beyond buffer
		if i >= len(buf) {
			break
		}

		// Append peer
		u.PeerList = append(u.PeerList[:], b2ip(buf[i:i+6]))
		i += 6
	}

	return u, nil
}

// ToBytes creates a packed byte array from a udpAnnounceResponse
func (u udpAnnounceResponse) ToBytes() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (uint32, must be 1 for announce)
	if u.Action != uint32(1) {
		return nil, errors.New("invalid action for udpAnnounceResponse")
	}

	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// Interval (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Interval); err != nil {
		return nil, err
	}

	// Leechers (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Leechers); err != nil {
		return nil, err
	}

	// Seeders (uint32)
	if err := binary.Write(res, binary.BigEndian, u.Seeders); err != nil {
		return nil, err
	}

	// PeerList, []compactPeer, iterated and compressed to compact format
	for _, peer := range u.PeerList{
		// Compact and write
		if err := binary.Write(res, binary.BigEndian, ip2b(peer.IP, peer.Port)); err != nil {
			return nil, err
		}
	}

	return res.Bytes(), nil
}
