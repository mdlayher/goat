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
	Static.LogChan <- fmt.Sprintf("announce: [up: %d, down: %d, left: %d]", announce.Uploaded, announce.Downloaded, announce.Left)

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
		file.Completed = 0
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

	// Update records AFTER user gets their response
	// NOTE: not goroutines because these occur after the client gets a response anyway, and because these
	// can create some interesting race conditions if not updated synchronously
	file.Save()
	fileUser.Save()
	user.Save()
}

// Report a bencoded []byte response as specified by input string
func TrackerError(err string) []byte {
	return bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"interval":       bencode.EncInt(RandRange(3200, 4000)),
		"min interval":   bencode.EncInt(1800),
	})
}
