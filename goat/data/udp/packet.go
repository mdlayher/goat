package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Packet represents the base parameters for a UDP tracker connection
type Packet struct {
	ConnID  uint64
	Action  uint32
	TransID uint32
}

// UnmarshalBinary creates a Packet from a packed byte array
func (u *Packet) UnmarshalBinary(buf []byte) (err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to create Packet from bytes")
		}
	}()

	// ConnID (uint64)
	u.ConnID = binary.BigEndian.Uint64(buf[0:8])

	// Action (uint32, connect: 0, announce: 1, scrape: 2, error: 3)
	u.Action = binary.BigEndian.Uint32(buf[8:12])

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[12:16])

	return nil
}

// MarshalBinary creates a packed byte array from a Packet
func (u Packet) MarshalBinary() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// ConnID (uint64)
	if err := binary.Write(res, binary.BigEndian, u.ConnID); err != nil {
		return nil, err
	}

	// Action (uint32, connect: 0, announce: 1, scrape: 2, error: 3)
	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
