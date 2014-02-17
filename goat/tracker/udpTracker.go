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
	}

	// Calculate file seeders and leechers
	seeders, err := file.Seeders()
	if err != nil {
		log.Println(err.Error())
	}
	announce.Seeders = uint32(seeders)

	leechers, err := file.Leechers()
	if err != nil {
		log.Println(err.Error())
	}
	announce.Leechers = uint32(leechers)

	// Convert to UDP byte buffer
	announceBuf, err := announce.MarshalBinary()
	if err != nil {
		log.Println(err.Error())
		return u.Error(ErrAnnounceFailure.Error())
	}

	// Numwant
	numwant, err := strconv.Atoi(query.Get("numwant"))
	if err != nil {
		numwant = 50
	}

	// Retrieve compact peer list
	// Note: because we are UDP, we send the second parameter 'false' to get
	// a "best guess" peer list, due to anonymous announces
	peers, err := file.CompactPeerList(numwant, false)
	if err != nil {
		log.Println(err.Error())
		return u.Error(ErrPeerListFailure.Error())
	}

	// Add compact peer list
	res := bytes.NewBuffer(announceBuf)
	err = binary.Write(res, binary.BigEndian, peers)
	if err != nil {
		log.Println(err.Error())
		return u.Error(ErrPeerListFailure.Error())
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
		return u.Error(ErrErrorFailure.Error())
	}

	return buf
}

// Protocol returns the protocol used by this tracker
func (u UDPTracker) Protocol() string {
	return "UDP"
}

// Scrape scrapes using UDP format
func (u UDPTracker) Scrape(files []data.FileRecord) []byte {
	// Buffered channel to receive UDP scrape stats structs
	resChan := make(chan *udp.ScrapeStats, len(files))

	// Iterate all files in parallel
	for _, f := range files {
		go func(f data.FileRecord, resChan chan *udp.ScrapeStats) {
			stat := udp.ScrapeStats{}

			// Seeders count
			var err error
			seeders, err := f.Seeders()
			if err != nil {
				log.Println(err.Error())
			}
			stat.Seeders = uint32(seeders)

			// Completion count
			completed, err := f.Completed()
			if err != nil {
				log.Println(err.Error())
			}
			stat.Completed = uint32(completed)

			// Leechers count
			leechers, err := f.Leechers()
			if err != nil {
				log.Println(err.Error())
			}
			stat.Leechers = uint32(leechers)

			// Return results on channel
			resChan <- &stat
		}(f, resChan)
	}

	// Fetch all results from channel
	stats := make([]udp.ScrapeStats, 0)
	for stat := range resChan {
		stats = append(stats[:], *stat)

		// Break once all file information has been received
		if len(stats) == len(files) {
			break
		}
	}

	// Close response channel
	close(resChan)

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
		return u.Error(ErrScrapeFailure.Error())
	}

	return buf
}
