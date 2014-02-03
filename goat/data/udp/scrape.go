package udp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
)

// ScrapeRequest represents a tracker scrape in the UDP format
type ScrapeRequest struct {
	ConnID     uint64
	Action     uint32
	TransID    uint32
	InfoHashes [][]byte
}

// UnmarshalBinary creates a ScrapeRequest from a packed byte array
func (u *ScrapeRequest) UnmarshalBinary(buf []byte) (err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to create ScrapeRequest from bytes")
		}
	}()

	// ConnID (uint64)
	u.ConnID = binary.BigEndian.Uint64(buf[0:8])

	// Action (uint32, must be 2 for scrape)
	u.Action = binary.BigEndian.Uint32(buf[8:12])
	if u.Action != uint32(2) {
		return fmt.Errorf("invalid action '%d' for ScrapeRequest", u.Action)
	}

	// TransID (uint32)
	u.TransID = binary.BigEndian.Uint32(buf[12:16])

	// Begin gathering info hashes
	u.InfoHashes = make([][]byte, 0)

	// Loop and iterate info_hash, up to 70 total (74 is said to be max by BEP15)
	for i := 16; i < 16+(70*20); i += 20 {
		// Validate that we are not appending nil bytes
		if i >= len(buf) || buf[i] == byte(0) {
			break
		}

		u.InfoHashes = append(u.InfoHashes[:], buf[i:i+20])
	}

	return nil
}

// MarshalBinary creates a packed byte array from a ScrapeRequest
func (u ScrapeRequest) MarshalBinary() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// ConnID (uint64)
	if err := binary.Write(res, binary.BigEndian, u.ConnID); err != nil {
		return nil, err
	}

	// Action (uint32, must be 2 for scrape)
	if u.Action != uint32(2) {
		return nil, fmt.Errorf("invalid action '%d' for ScrapeRequest", u.Action)
	}

	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// InfoHashes ([]byte, 20 bytes each, iterated)
	for _, hash := range u.InfoHashes {
		// Ensure all hashes are 20 bytes
		if len(hash) != 20 {
			return nil, errors.New("info_hash must be exactly 20 bytes")
		}

		// Write each hash
		if err := binary.Write(res, binary.BigEndian, hash); err != nil {
			return nil, err
		}
	}

	return res.Bytes(), nil
}

// ToValues creates a url.Values struct from a ScrapeRequest
func (u ScrapeRequest) ToValues() url.Values {
	// Initialize query map
	query := url.Values{}
	query.Set("udp", "1")

	// Iterate info hashes and convert into strings
	hashes := make([]string, 0)
	for _, hash := range u.InfoHashes {
		hashes = append(hashes[:], string(hash))
	}
	query["info_hash"] = hashes

	// Return final query map
	return query
}

// ScrapeResponse represents a tracker scrape response in the UDP format
type ScrapeResponse struct {
	Action    uint32
	TransID   uint32
	FileStats []ScrapeStats
}

// ScrapeStats represents one dictionary of stats about a file from a UDP scrape response
type ScrapeStats struct {
	Seeders   uint32
	Completed uint32
	Leechers  uint32
}

// UnmarshalBinary creates a ScrapeResponse from a packed byte array
func (u *ScrapeResponse) UnmarshalBinary(buf []byte) (err error) {
	// Set up recovery function to catch a panic as an error
	// This will run if we attempt to access an out of bounds index
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("failed to create ScrapeResponse from bytes")
		}
	}()

	// Action (must be 2 for scrape)
	u.Action = binary.BigEndian.Uint32(buf[0:4])
	if u.Action != uint32(2) {
		return fmt.Errorf("invalid action '%d' for ScrapeResponse", u.Action)
	}

	// Transaction ID
	u.TransID = binary.BigEndian.Uint32(buf[4:8])

	// FileStats
	u.FileStats = make([]ScrapeStats, 0)

	// Iterate file stats buffer
	i := 8
	for {
		// Validate that we are not seeking beyond buffer
		if i >= len(buf) {
			break
		}

		// File stats
		stats := ScrapeStats{}

		// Seeders
		stats.Seeders = binary.BigEndian.Uint32(buf[i : i+4])
		i += 4

		// Completed
		stats.Completed = binary.BigEndian.Uint32(buf[i : i+4])
		i += 4

		// Leechers
		stats.Leechers = binary.BigEndian.Uint32(buf[i : i+4])
		i += 4

		// Append stats
		u.FileStats = append(u.FileStats[:], stats)
	}

	return nil
}

// MarshalBinary creates a packed byte array from a ScrapeResponse
func (u ScrapeResponse) MarshalBinary() ([]byte, error) {
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (uint32, must be 2 for scrape)
	if u.Action != uint32(2) {
		return nil, fmt.Errorf("invalid action '%d' for ScrapeResponse", u.Action)
	}

	if err := binary.Write(res, binary.BigEndian, u.Action); err != nil {
		return nil, err
	}

	// TransID (uint32)
	if err := binary.Write(res, binary.BigEndian, u.TransID); err != nil {
		return nil, err
	}

	// FileStats ([]ScrapeStats, iterated)
	for _, stats := range u.FileStats {
		// Seeders
		if err := binary.Write(res, binary.BigEndian, stats.Seeders); err != nil {
			return nil, err
		}

		// Completed
		if err := binary.Write(res, binary.BigEndian, stats.Completed); err != nil {
			return nil, err
		}

		// Leechers
		if err := binary.Write(res, binary.BigEndian, stats.Leechers); err != nil {
			return nil, err
		}
	}

	return res.Bytes(), nil
}
