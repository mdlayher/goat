package goat

import (
	"bencode"
	"fmt"
	"strconv"
)

// Tracker announce request
func TrackerAnnounce(user UserRecord, query map[string]string, resChan chan []byte) {
	// Store announce information in struct
	announce := MapToAnnounceLog(query)

	// Request to store announce
	go announce.Save()

	// Only report event when needed
	event := ""
	if announce.Event != "" {
		event = announce.Event + " "
	}

	Static.LogChan <- fmt.Sprintf("announce: [%s:%d] %s%s", announce.Ip, announce.Port, event, announce.InfoHash)

	// Check for a matching file via info_hash
	file := new(FileRecord).Load(announce.InfoHash, "info_hash")
	if file == (FileRecord{}) {
		// Torrent is not currently registered
		resChan <- TrackerError("Unregistered torrent")

		// Create an entry in file table for this hash, but mark it as unverified
		file.InfoHash = announce.InfoHash
		file.Verified = false
		file.Completed = 0

		Static.LogChan <- fmt.Sprintf("tracker: detected new file, awaiting manual approval [hash: %s]", announce.InfoHash)

		go file.Save()
		return
	}

	// Ensure file is verified, meaning we will permit tracking of it
	if !file.Verified {
		resChan <- TrackerError("Unverified torrent")
		return
	}

	// Check existing record for this user with this file
	fileUser := new(FileUserRecord).Load(file.Id, user.Id)

	// New user, starting torrent
	if announce.Event == "started" && fileUser == (FileUserRecord{}) {
		// Create new relationship
		fileUser.FileId = file.Id
		fileUser.UserId = user.Id
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

			// If status completed, mark file as completed by another user
			if announce.Event == "completed" {
				file.Completed = file.Completed + 1
				go file.Save()
			}
		} else {
			fileUser.Completed = false
		}

		// Add an announce
		fileUser.Announced = fileUser.Announced + 1

		// Store latest statistics, but do so in a sane way (no removing upload/download, no adding left)
		// NOTE: clients report absolute values, so delta should NEVER be calculated for these
		// NOTE: It is also worth noting that if a client re-downloads a file they have previously downloaded,
		// but the FileUserRecord relationship is not cleared, they will essentially get a "free" download, with
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

	// Begin generating response map, with current number of known seeders/leechers
	res := map[string][]byte{
		"complete":   bencode.EncInt(file.Seeders()),
		"incomplete": bencode.EncInt(file.Leechers()),
	}

	// If client has not yet completed torrent, ask them to announce more frequently, so they can gather
	// more peers and quickly report their statistics
	if fileUser.Completed == false {
		res["interval"] = bencode.EncInt(RandRange(300, 600))
		res["min interval"] = bencode.EncInt(300)
	} else {
		// Once a torrent has been completed, report statistics less frequently
		res["interval"] = bencode.EncInt(RandRange(3200, 4000))
		res["min interval"] = bencode.EncInt(1800)
	}

	// Check for numwant parameter, return up to that number of peers
	// Default is 50 per protocol
	numwant := 50
	if _, ok := query["numwant"]; ok {
		// Verify numwant is an integer
		num, err := strconv.Atoi(query["numwant"])
		if err == nil {
			numwant = num
		}
	}

	// Generaate compact peer list of length numwant, exclude this user
	res["peers"] = bencode.EncBytes(file.PeerList(query["ip"], numwant))

	// Bencode entire map and return
	resChan <- bencode.EncDictMap(res)
}

// Report a bencoded []byte response as specified by input string
func TrackerError(err string) []byte {
	return bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"interval":       bencode.EncInt(RandRange(3200, 4000)),
		"min interval":   bencode.EncInt(1800),
	})
}
