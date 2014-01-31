package goat

import (
	"log"
	"net/url"
)

// torrentTracker defines the common interface for trackers to generate their responses
type torrentTracker interface {
	Announce(url.Values, fileRecord) []byte
	Error(string) []byte
	Protocol() string
	Scrape([]fileRecord) []byte
}

// trackerAnnounce announces a tracker request
func trackerAnnounce(tracker torrentTracker, user userRecord, query url.Values) []byte {
	// Store announce information in struct
	announce := new(announceLog).FromValues(query)
	if announce == (announceLog{}) {
		return tracker.Error("Malformed announce")
	}

	// Request to store announce
	go announce.Save()

	// Only report event when needed
	event := ""
	if announce.Event != "" {
		event = announce.Event + " "
	}

	log.Printf("announce: [%s %s:%d] %s%s", tracker.Protocol(), announce.IP, announce.Port, event, announce.InfoHash)

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
		return tracker.Error("Unregistered torrent")
	}

	// Ensure file is verified, meaning we will permit tracking of it
	if !file.Verified {
		return tracker.Error("Unverified torrent")
	}

	// Launch peer reaper to remove old peers from this file
	go file.PeerReaper()

	// If UDP tracker, we cannot reliably detect user, so we announce anonymously
	if _, ok := tracker.(udpTracker); ok {
		return tracker.Announce(query, file)
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
	return tracker.Announce(query, file)
}

// trackerScrape scrapes a tracker request
func trackerScrape(tracker torrentTracker, query url.Values) []byte {
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
			return tracker.Error("Malformed scrape")
		}

		// Request to store scrape
		go scrape.Save()

		log.Printf("scrape: [%s %s] %s", tracker.Protocol(), scrape.IP, scrape.InfoHash)

		// Check for a matching file via info_hash
		file := new(fileRecord).Load(scrape.InfoHash, "info_hash")
		if file == (fileRecord{}) {
			// Torrent is not currently registered
			return tracker.Error("Unregistered torrent")
		}

		// Ensure file is verified, meaning we will permit scraping of it
		if !file.Verified {
			return tracker.Error("Unverified torrent")
		}

		// Launch peer reaper to remove old peers from this file
		go file.PeerReaper()

		// File is valid, add it to list to be scraped
		scrapeFiles = append(scrapeFiles[:], file)
	}

	// Create scrape
	return tracker.Scrape(scrapeFiles)
}
