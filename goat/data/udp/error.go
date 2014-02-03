package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// ErrorResponse represents a tracker error response in the UDP format
type ErrorResponse struct {
	Action  uint32
	TransID uint32
	Error   string
}

// UnmarshalBinary creates a ErrorResponse from a packed byte array
func (u *ErrorResponse) UnmarshalBinary(buf []byte) (err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to create ErrorResponse from bytes")
		}
	}()

	// Action (uint32, must be 3 for error)
	u.Action = binary.BigEndian.Uint32(buf[0:4])
	if u.Action != uint32(3) {
		return fmt.Errorf("invalid action '%d' for ErrorResponse", u.Action)
	}

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[4:8])

	// Error (string)
	u.Error = string(buf[8:len(buf)])

	return nil
}

// MarshalBinary creates a packed byte array from a ErrorResponse
func (u ErrorResponse) MarshalBinary() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (uint32, must be 3 for error)
	if u.Action != uint32(3) {
		return nil, fmt.Errorf("invalid action '%d' for ErrorResponse", u.Action)
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
