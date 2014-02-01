package goat

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/url"
	"strconv"
)

// udpTracker generates responses in the UDP datagram format
type udpTracker struct {
	TransID uint32
}

// Announce announces using UDP format
func (u udpTracker) Announce(query url.Values, file fileRecord) []byte {
	// Create UDP announce response
	announce := udpAnnounceResponse{
		Action:   1,
		TransID:  u.TransID,
		Interval: uint32(static.Config.Interval),
		Leechers: uint32(file.Leechers()),
		Seeders:  uint32(file.Seeders()),
	}

	// Convert to UDP byte buffer
	announceBuf, err := announce.ToBytes()
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	// Numwant
	numwant, err := strconv.Atoi(query.Get("numwant"))
	if err != nil {
		numwant = 50
	}

	// Add compact peer list
	res := bytes.NewBuffer(announceBuf)
	err = binary.Write(res, binary.BigEndian, file.PeerList(query.Get("ip"), numwant))
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP announce response")
	}

	return res.Bytes()
}

// Error reports a UDP []byte response packed datagram
func (u udpTracker) Error(msg string) []byte {
	// Create UDP error response
	errRes := udpErrorResponse{
		Action:  3,
		TransID: u.TransID,
		Error:   msg,
	}

	// Convert to UDP byte buffer
	buf, err := errRes.ToBytes()
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP error response")
	}

	return buf
}

// Protocol returns the protocol used by this tracker
func (u udpTracker) Protocol() string {
	return "UDP"
}

// Scrape scrapes using UDP format
func (u udpTracker) Scrape(files []fileRecord) []byte {
	// Iterate all files, grabbing their statistics
	stats := make([]udpScrapeStats, 0)
	for _, file := range files {
		stat := udpScrapeStats{}

		// Seeders
		stat.Seeders = uint32(file.Seeders())

		// Completed
		stat.Completed = uint32(file.Completed())

		// Leechers
		stat.Leechers = uint32(file.Leechers())

		// Append to slice
		stats = append(stats[:], stat)
	}

	// Create UDP scrape response
	scrape := udpScrapeResponse{
		Action:    2,
		TransID:   u.TransID,
		FileStats: stats,
	}

	// Convert to UDP byte buffer
	buf, err := scrape.ToBytes()
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP scrape response")
	}

	return buf
}
