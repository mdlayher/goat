package tracker

import (
	"bytes"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"

	// Import bencode library
	bencode "code.google.com/p/bencode-go"
)

// HTTPTracker generates responses in the HTTP bencode format
type HTTPTracker struct {
}

// AnnounceResponse defines the response structure of an HTTP tracker announce
type AnnounceResponse struct {
	Complete    int    "complete"
	Incomplete  int    "incomplete"
	Interval    int    "interval"
	MinInterval int    "min interval"
	Peers       string "peers"
}

// Announce announces using HTTP format
func (h HTTPTracker) Announce(query url.Values, file data.FileRecord) []byte {
	// Generate response struct
	announce := AnnounceResponse{
		Interval:    common.Static.Config.Interval,
		MinInterval: common.Static.Config.Interval / 2,
	}

	// Get seeders count on file
	var err error
	announce.Complete, err = file.Seeders()
	if err != nil {
		log.Println(err.Error())
	}

	// Get leechers count on file
	announce.Incomplete, err = file.Leechers()
	if err != nil {
		log.Println(err.Error())
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
		return h.Error(ErrAnnounceFailure.Error())
	}

	// Generate compact peer list of length numwant
	// Note: because we are HTTP, we can mark second parameter as 'true' to get a
	// more accurate peer list
	compactPeers, err := file.CompactPeerList(numwant, true)
	if err != nil {
		log.Println(err.Error())
		return h.Error(ErrPeerListFailure.Error())
	}

	// Because the bencode marshaler does not handle compact, binary peer list conversion,
	// we handle it manually here.

	// Get initial buffer, chop off 3 bytes: "0:e", append the actual list length with new colon
	out := buf.Bytes()
	out = append(out[0:len(out)-3], []byte(strconv.Itoa(len(compactPeers))+":")...)

	// Append peers list, terminate with an "e"
	return append(append(out, compactPeers...), byte('e'))
}

// errorResponse defines the response structure of an HTTP tracker error
type errorResponse struct {
	FailureReason string "failure reason"
	Interval      int    "interval"
	MinInterval   int    "min interval"
}

// Error reports a bencoded []byte response as specified by input string
func (h HTTPTracker) Error(err string) []byte {
	res := errorResponse{
		FailureReason: err,
		Interval:      common.Static.Config.Interval,
		MinInterval:   common.Static.Config.Interval / 2,
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
func (h HTTPTracker) Protocol() string {
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
func (h HTTPTracker) Scrape(files []data.FileRecord) []byte {
	// Response struct
	scrape := scrapeResponse{
		Files: make(map[string]scrapeFile),
	}

	// WaitGroup to wait for all scrape file entries to be generated
	var wg sync.WaitGroup
	wg.Add(len(files))

	// Mutex for safe locking on map writes
	var mutex sync.RWMutex

	// Iterate all files in parallel
	for _, f := range files {
		go func(f data.FileRecord, scrape *scrapeResponse, mutex *sync.RWMutex, wg *sync.WaitGroup) {
			// Generate scrapeFile struct
			fileInfo := scrapeFile{}
			var err error

			// Seeders count
			fileInfo.Complete, err = f.Seeders()
			if err != nil {
				log.Println(err.Error())
			}

			// Completion count
			fileInfo.Downloaded, err = f.Completed()
			if err != nil {
				log.Println(err.Error())
			}

			// Leechers count
			fileInfo.Incomplete, err = f.Leechers()
			if err != nil {
				log.Println(err.Error())
			}

			// Add hash and file info to map
			mutex.Lock()
			scrape.Files[f.InfoHash] = fileInfo
			mutex.Unlock()

			// Inform waitgroup that this file is ready
			wg.Done()
		}(f, &scrape, &mutex, &wg)
	}

	// Wait for all information to be generated
	wg.Wait()

	// Marshal struct into bencode
	buf := bytes.NewBuffer(make([]byte, 0))
	if err := bencode.Marshal(buf, scrape); err != nil {
		log.Println(err.Error())
		return h.Error(ErrScrapeFailure.Error())
	}

	return buf.Bytes()
}
