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

	Static.LogChan <- fmt.Sprintf("announce: [ip: %s, port:%d]", announce.Ip, announce.Port)
	Static.LogChan <- fmt.Sprintf("announce: [info_hash: %s]", announce.InfoHash)

	// Only report announce when needed
	if announce.Event != "" {
		Static.LogChan <- fmt.Sprintf("announce: [event: %s]", announce.Event)
	}

	// Check for a matching file via info_hash
	file := new(FileRecord).Load(announce.InfoHash, "info_hash")
	if file == (FileRecord{}) {
		// Torrent is not currently registered
		resChan <- TrackerError("Unregistered torrent")

		// Create an entry in file table for this hash, but mark it as unverified
		file.InfoHash = announce.InfoHash
		file.Verified = false
		file.Leechers = 0
		file.Seeders = 0
		file.Completed = 0
		go file.Save()
		return
	}

	// Ensure file is verified, meaning we will permit tracking of it
	if !file.Verified {
		resChan <- TrackerError("Unverified torrent")
		return
	}

	// Check for existing record for this user with this file
	fileUser := new(FileUserRecord).Load(file.Id, user.Id)
	if fileUser == (FileUserRecord{}) {
		// Create new relationship
		fileUser.FileId = file.Id
		fileUser.UserId = user.Id
		fileUser.Active = true
		fileUser.Completed = false
		fileUser.Announced = 1
		fileUser.Uploaded = announce.Uploaded
		fileUser.Downloaded = announce.Downloaded
		fileUser.Left = announce.Left

		// Add a leecher to this file, UNLESS they have already completed it
		if announce.Left == 0 || announce.Event == "completed" {
			file.Seeders = file.Seeders + 1
		} else {
			file.Leechers = file.Leechers + 1
		}
	} else {
		// Else, pre-existing record, so update
		// Check for stopped status
		if announce.Event != "stopped" {
			fileUser.Active = true
		} else {
			fileUser.Active = false

			// Remove seeder if applicable
			if announce.Left == 0 && file.Seeders > 0 {
				file.Seeders = file.Seeders - 1
			}
		}

		// Check for completion
		if announce.Event == "completed" && announce.Left == 0 {
			fileUser.Completed = true

			// Mark file as completed by another user
			file.Completed = file.Completed + 1

			// Decrement leecher, add seeder
			if (file.Leechers > 0) {
				file.Leechers = file.Leechers - 1
			}
			file.Seeders = file.Seeders + 1
		} else {
			fileUser.Completed = false
		}

		// Add an announce
		fileUser.Announced = fileUser.Announced + 1

		// Store latest statistics, but do so in a sane way (no removing upload/download, no adding left)
		if (announce.Uploaded > fileUser.Uploaded) {
			fileUser.Uploaded = announce.Uploaded
		}
		if (announce.Downloaded > fileUser.Downloaded) {
			fileUser.Downloaded = announce.Downloaded
		}
		if (announce.Left < fileUser.Left) {
			fileUser.Left = announce.Left
		}
	}

	// Update File record
	go file.Save()

	// Update User record
	go user.Save()

	// Insert or update the FileUser record
	go fileUser.Save()

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

	// Tracker announce response, with generated peerlist of length numwant, excluding this user
	resChan <- bencode.EncDictMap(map[string][]byte{
		"interval":     bencode.EncInt(RandRange(3200, 4000)),
		"min interval": bencode.EncInt(1800),
		"peers":        bencode.EncBytes(file.PeerList(query["ip"], numwant)),
	})
}

// Report a bencoded []byte response as specified by input string
func TrackerError(err string) []byte {
	return bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"interval":       bencode.EncInt(RandRange(3200, 4000)),
		"min interval":   bencode.EncInt(1800),
	})
}
