package data

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

// Peer represents an IP and port peer, used as part of the peer list
type Peer struct {
	IP   string
	Port uint16
}

// MarshalBinary creates a packed byte array from a peer
func (p Peer) MarshalBinary() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// IP (uint32)
	if err := binary.Write(res, binary.BigEndian, binary.BigEndian.Uint32((net.ParseIP(p.IP).To4()))); err != nil {
		return nil, err
	}

	// Port (uint16)
	if err := binary.Write(res, binary.BigEndian, uint16(p.Port)); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}

// UnmarshalBinary creates a Peer from a packed byte array
func (p *Peer) UnmarshalBinary(buf []byte) (err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to create Peer from bytes")
		}
	}()

	// IP (uint32 -> string)
	p.IP = net.IPv4(buf[0], buf[1], buf[2], buf[3]).String()

	// Port (uint16)
	p.Port = binary.BigEndian.Uint16(buf[4:6])

	return nil
}
