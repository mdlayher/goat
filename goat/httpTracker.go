package goat

import (
	"bytes"
	"log"
	"net/url"
	"strconv"

	// Import bencode library
	bencode "code.google.com/p/bencode-go"
)

// httpTracker generates responses in the HTTP bencode format
type httpTracker struct {
}

// announceResponse defines the response structure of an HTTP tracker announce
type announceResponse struct {
	Complete    int    "complete"
	Incomplete  int    "incomplete"
	Interval    int    "interval"
	MinInterval int    "min interval"
	Peers       string "peers"
}

// h.Announce announces using HTTP format
func (h httpTracker) Announce(query url.Values, file fileRecord) []byte {
	// Generate response struct
	announce := announceResponse{
		Complete:    file.Seeders(),
		Incomplete:  file.Leechers(),
		Interval:    static.Config.Interval,
		MinInterval: static.Config.Interval / 2,
	}

	// Check for numwant parameter, return up to that number of peers
	// Default is 50 per protocol
	numwant := 50
	if query.Get("numwant") != "" {
		// Verify numwant is an integer
		num, err := strconv.Atoi(query.Get("numwant"))
		if err == nil {
			numwant = num
		}
	}

	// Marshal struct into bencode
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := bencode.Marshal(buf, announce); err != nil {
		log.Println(err.Error())
		return h.Error("Tracker error: failed to create announce response")
	}

	// Generate compact peer list of length numwant, exclude this user
	peers := file.PeerList(query.Get("ip"), numwant)

	// Because the bencode marshaler does not handle compact, binary peer list conversion,
	// we handle it manually here.

	// Get initial buffer, chop off 3 bytes: "0:e", append the actual list length with new colon
	out := buf.Bytes()
	out = append(out[0:len(out)-3], []byte(strconv.Itoa(len(peers))+":")...)

	// Append peers list, terminate with an "e"
	out = append(append(out, peers...), byte('e'))

	// Return final announce message
	return out
}

// errorResponse defines the response structure of an HTTP tracker error
type errorResponse struct {
	FailureReason string "failure reason"
	Interval      int    "interval"
	MinInterval   int    "min interval"
}

// Error reports a bencoded []byte response as specified by input string
func (h httpTracker) Error(err string) []byte {
	res := errorResponse{
		FailureReason: err,
		Interval:      static.Config.Interval,
		MinInterval:   static.Config.Interval / 2,
	}

	// Marshal struct into bencode
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := bencode.Marshal(buf, res); err != nil {
		log.Println(err.Error())
		return nil
	}

	return buf.Bytes()
}

// Protocol returns the protocol used by this tracker
func (h httpTracker) Protocol() string {
	return "HTTP"
}

// scrapeResponse defines the top-level response structure of an HTTP tracker scrape
type scrapeResponse struct {
	Files map[string]scrapeFile "files"
}

// scrapeFile defines the fields of a scrape response for a single info_hash
type scrapeFile struct {
	Complete   int "complete"
	Downloaded int "downloaded"
	Incomplete int "incomplete"
	// optional field: Name string "name"
}

// Scrape reports scrape for one or more files, using HTTP format
func (h httpTracker) Scrape(files []fileRecord) []byte {
	// Response struct
	scrape := scrapeResponse{
		Files: make(map[string]scrapeFile),
	}

	// Iterate all files
	for _, file := range files {
		// Generate scrapeFile struct
		fileInfo := scrapeFile{
			Complete:   file.Seeders(),
			Downloaded: file.Completed(),
			Incomplete: file.Leechers(),
		}

		// Add hash and file info to map
		scrape.Files[file.InfoHash] = fileInfo
	}

	// Marshal struct into bencode
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := bencode.Marshal(buf, scrape); err != nil {
		log.Println(err.Error())
		return h.Error("Tracker error: failed to create scrape response")
	}

	return buf.Bytes()
}
