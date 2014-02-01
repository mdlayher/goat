package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// ConnectResponse represents a connect response for a UDP tracker connection
type ConnectResponse struct {
	Action  uint32
	TransID uint32
	ConnID  uint64
}

// FromBytes creates a ConnectResponse from a packed byte array
func (u ConnectResponse) FromBytes(buf []byte) (p ConnectResponse, err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			p = ConnectResponse{}
			err = errors.New("failed to create ConnectResponse from bytes")
		}
	}()

	// Action (uint32, must be 0 for connect)
	u.Action = binary.BigEndian.Uint32(buf[0:4])
	if u.Action != uint32(0) {
		return ConnectResponse{}, errors.New("invalid action for ConnectResponse")
	}

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[4:8])

	// ConnID (uint64)
	u.ConnID = binary.BigEndian.Uint64(buf[8:16])

	return u, nil
}

// ToBytes creates a packed byte array from a ConnectResponse
func (u ConnectResponse) ToBytes() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (uint32, must be 0 for connect)
	if u.Action != uint32(0) {
		return nil, errors.New("invalid action for ConnectResponse")
	}

	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// ConnID (uint64)
	if err := binary.Write(res, binary.BigEndian, u.ConnID); err != nil {
		return nil, err
	}

	return res.Bytes(), nil
}
