package goat

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/url"
	"strconv"

	// Import bencode library
	bencode "code.google.com/p/bencode-go"
)

// trackerAnnounce nnounces a tracker request
func trackerAnnounce(user userRecord, query url.Values, transID []byte) []byte {
	// Store announce information in struct
	announce := new(announceLog).FromValues(query)
	if announce == (announceLog{}) {
		// If a transaction ID is available, this must be a UDP announce
		if transID == nil {
			return httpTrackerError("Malformed announce")
		}

		return udpTrackerError("Malformed announce", transID)
	}

	// Request to store announce
	go announce.Save()

	// Only report event when needed
	event := ""
	if announce.Event != "" {
		event = announce.Event + " "
	}

	// Report protocol
	proto := ""
	if announce.UDP {
		proto = " udp"
	} else {
		proto = "http"
	}

	log.Printf("announce: [%s %s:%d] %s%s", proto, announce.IP, announce.Port, event, announce.InfoHash)

	// Check for a matching file via info_hash
	file := new(fileRecord).Load(announce.InfoHash, "info_hash")
	if file == (fileRecord{}) {
		// Torrent is not currently registered
		log.Printf("tracker: detected new file, awaiting manual approval [hash: %s]", announce.InfoHash)

		// Create an entry in file table for this hash, but mark it as unverified
		file.InfoHash = announce.InfoHash
		file.Verified = false
		go file.Save()

		// Report error
		if !announce.UDP {
			return httpTrackerError("Unregistered torrent")
		}

		return udpTrackerError("Unregistered torrent", transID)
	}

	// Ensure file is verified, meaning we will permit tracking of it
	if !file.Verified {
		if !announce.UDP {
			return httpTrackerError("Unverified torrent")
		}

		return udpTrackerError("Unverified torrent", transID)
	}

	// Launch peer reaper to remove old peers from this file
	go file.PeerReaper()

	// If UDP tracker, we cannot reliably detect user, so we announce anonymously
	if announce.UDP {
		return udpTrackerAnnounce(query, file, transID)
	}

	// Check existing record for this user with this file and this IP
	fileUser := new(fileUserRecord).Load(file.ID, user.ID, query.Get("ip"))

	// New user, starting torrent
	if fileUser == (fileUserRecord{}) {
		// Create new relationship
		fileUser.FileID = file.ID
		fileUser.UserID = user.ID
		fileUser.IP = query.Get("ip")
		fileUser.Active = true
		fileUser.Announced = 1

		// If announce reports 0 left, but no existing record, user is probably the initial seeder
		if announce.Left == 0 {
			fileUser.Completed = true
		} else {
			fileUser.Completed = false
		}

		// Track the initial uploaded, download, and left values
		// NOTE: clients report absolute values, so delta should NEVER be calculated for these
		fileUser.Uploaded = announce.Uploaded
		fileUser.Downloaded = announce.Downloaded
		fileUser.Left = announce.Left
	} else {
		// Else, pre-existing record, so update
		// Event "stopped", mark as inactive
		// NOTE: likely only reported by clients which are actively seeding, NOT when stopped during leeching
		if announce.Event == "stopped" {
			fileUser.Active = false
		} else {
			// Else, "started", "completed", or no status, mark as active
			fileUser.Active = true
		}

		// Check for completion
		// Could be from a peer stating completed, or a seed reporting 0 left
		if announce.Event == "completed" || announce.Left == 0 {
			fileUser.Completed = true
		} else {
			fileUser.Completed = false
		}

		// Add an announce
		fileUser.Announced = fileUser.Announced + 1

		// Store latest statistics, but do so in a sane way (no removing upload/download, no adding left)
		// NOTE: clients report absolute values, so delta should NEVER be calculated for these
		// NOTE: It is also worth noting that if a client re-downloads a file they have previously downloaded,
		// but the fileUserRecord relationship is not cleared, they will essentially get a "free" download, with
		// no extra download penalty to their share ratio
		// For the time being, this behavior will be expected and acceptable
		if announce.Uploaded > fileUser.Uploaded {
			fileUser.Uploaded = announce.Uploaded
		}
		if announce.Downloaded > fileUser.Downloaded {
			fileUser.Downloaded = announce.Downloaded
		}
		if announce.Left < fileUser.Left {
			fileUser.Left = announce.Left
		}
	}

	// Update file/user relationship record
	go fileUser.Save()

	// Create announce
	return httpTrackerAnnounce(query, file, fileUser)
}

// trackerScrape scrapes a tracker request
func trackerScrape(user userRecord, query url.Values) []byte {
	// List of files to be scraped
	scrapeFiles := make([]fileRecord, 0)

	// Iterate all info_hash values in query
	for _, infoHash := range query["info_hash"] {
		// Make a copy of query, set the info hash as current in loop
		localQuery := query
		localQuery.Set("info_hash", infoHash)

		// Store scrape information in struct
		scrape := new(scrapeLog).FromValues(localQuery)
		if scrape == (scrapeLog{}) {
			return httpTrackerError("Malformed scrape")
		}

		// Request to store scrape
		go scrape.Save()

		log.Printf("scrape: [%s] %s", scrape.IP, scrape.InfoHash)

		// Check for a matching file via info_hash
		file := new(fileRecord).Load(scrape.InfoHash, "info_hash")
		if file == (fileRecord{}) {
			// Torrent is not currently registered
			return httpTrackerError("Unregistered torrent")
		}

		// Ensure file is verified, meaning we will permit scraping of it
		if !file.Verified {
			return httpTrackerError("Unverified torrent")
		}

		// Launch peer reaper to remove old peers from this file
		go file.PeerReaper()

		// File is valid, add it to list to be scraped
		scrapeFiles = append(scrapeFiles[:], file)
	}

	// Create scrape
	return httpTrackerScrape(scrapeFiles)
}

// announceResponse defines the response structure of an HTTP tracker announce
type announceResponse struct {
	Complete    int    "complete"
	Incomplete  int    "incomplete"
	Interval    int    "interval"
	MinInterval int    "min interval"
	Peers       string "peers"
}

// httpTrackerAnnounce announces using HTTP format
func httpTrackerAnnounce(query url.Values, file fileRecord, fileUser fileUserRecord) []byte {
	// Generate response struct
	announce := announceResponse{
		Complete:   file.Seeders(),
		Incomplete: file.Leechers(),
	}

	// If client has not yet completed torrent, ask them to announce more frequently, so they can gather
	// more peers and quickly report their statistics
	if fileUser.Completed == false {
		announce.Interval = 600
		announce.MinInterval = 300
	} else {
		// Once a torrent has been completed, report statistics less frequently
		announce.Interval = randRange(static.Config.Interval-600, static.Config.Interval)
		announce.MinInterval = static.Config.Interval / 2
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
		return httpTrackerError("Tracker error: failed to create announce response")
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

// httpTrackerScrape reports scrape for one or more files, using HTTP format
func httpTrackerScrape(files []fileRecord) []byte {
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
		return httpTrackerError("Tracker error: failed to create scrape response")
	}

	return buf.Bytes()
}

// errorResponse defines the response structure of an HTTP tracker error
type errorResponse struct {
	FailureReason string "failure reason"
	Interval      int    "interval"
	MinInterval   int    "min interval"
}

// httpTrackerError reports a bencoded []byte response as specified by input string
func httpTrackerError(err string) []byte {
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

// udpTrackerAnnounce announces using UDP format
func udpTrackerAnnounce(query url.Values, file fileRecord, transID []byte) []byte {
	// Response buffer
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (1 for announce)
	err := binary.Write(res, binary.BigEndian, uint32(1))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	// Transaction ID
	err = binary.Write(res, binary.BigEndian, transID)
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	// Interval
	err = binary.Write(res, binary.BigEndian, uint32(randRange(static.Config.Interval-600, static.Config.Interval)))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	// Leechers
	err = binary.Write(res, binary.BigEndian, uint32(file.Leechers()))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	// Seeders
	err = binary.Write(res, binary.BigEndian, uint32(file.Seeders()))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	// Peer list
	numwant, err := strconv.Atoi(query.Get("numwant"))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	err = binary.Write(res, binary.BigEndian, file.PeerList(query.Get("ip"), numwant))
	if err != nil {
		log.Println(err.Error())
		return udpTrackerError("Could not create UDP announce response", transID)
	}

	return res.Bytes()
}

// udpTrackerError reports a []byte response packed datagram
func udpTrackerError(msg string, transID []byte) []byte {
	// Response buffer
	res := bytes.NewBuffer(make([]byte, 0))

	// Action (3 for error)
	err := binary.Write(res, binary.BigEndian, uint32(3))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Transaction ID
	err = binary.Write(res, binary.BigEndian, transID)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Error message
	err = binary.Write(res, binary.BigEndian, []byte(msg))
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return res.Bytes()
}
