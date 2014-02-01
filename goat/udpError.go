package goat

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// udpErrorResponse represents a tracker error response in the UDP format
type udpErrorResponse struct {
	Action  uint32
	TransID uint32
	Error   string
}

// FromBytes creates a udpErrorResponse from a packed byte array
func (u udpErrorResponse) FromBytes(buf []byte) (p udpErrorResponse, err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			p = udpErrorResponse{}
			err = errors.New("failed to create udpErrorResponse from bytes")
		}
	}()

	// Action (uint32, must be 3 for error)
	u.Action = binary.BigEndian.Uint32(buf[0:4])
	if u.Action != uint32(3) {
		return udpErrorResponse{}, errors.New("invalid action for udpErrorResponse")
	}

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[4:8])

	// Error (string)
	u.Error = string(buf[8:len(buf)])

	return u, nil
}

// ToBytes creates a packed byte array from a udpErrorResponse
func (u udpErrorResponse) ToBytes() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (uint32, must be 3 for error)
	if u.Action != uint32(3) {
		return nil, errors.New("invalid action for udpErrorResponse")
	}

	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// Error (string)
	if err := binary.Write(res, binary.BigEndian, []byte(u.Error)); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
