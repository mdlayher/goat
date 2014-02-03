package tracker

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/url"
	"strconv"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
	"github.com/mdlayher/goat/goat/data/udp"
)

// UDPTracker generates responses in the UDP datagram format
type UDPTracker struct {
	TransID uint32
}

// Announce announces using UDP format
func (u UDPTracker) Announce(query url.Values, file data.FileRecord) []byte {
	// Create UDP announce response
	announce := udp.AnnounceResponse{
		Action:   1,
		TransID:  u.TransID,
		Interval: uint32(common.Static.Config.Interval),
		Leechers: uint32(file.Leechers()),
		Seeders:  uint32(file.Seeders()),
	}

	// Convert to UDP byte buffer
	announceBuf, err := announce.MarshalBinary()
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
func (u UDPTracker) Error(msg string) []byte {
	// Create UDP error response
	errRes := udp.ErrorResponse{
		Action:  3,
		TransID: u.TransID,
		Error:   msg,
	}

	// Convert to UDP byte buffer
	buf, err := errRes.MarshalBinary()
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP error response")
	}

	return buf
}

// Protocol returns the protocol used by this tracker
func (u UDPTracker) Protocol() string {
	return "UDP"
}

// Scrape scrapes using UDP format
func (u UDPTracker) Scrape(files []data.FileRecord) []byte {
	// Iterate all files, grabbing their statistics
	stats := make([]udp.ScrapeStats, 0)
	for _, file := range files {
		stat := udp.ScrapeStats{}

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
	scrape := udp.ScrapeResponse{
		Action:    2,
		TransID:   u.TransID,
		FileStats: stats,
	}

	// Convert to UDP byte buffer
	buf, err := scrape.MarshalBinary()
	if err != nil {
		log.Println(err.Error())
		return u.Error("Could not create UDP scrape response")
	}

	return buf
}
