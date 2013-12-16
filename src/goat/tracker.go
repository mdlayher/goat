package goat

import (
	"bencode"
	"fmt"
	"strconv"
)

// Tracker announce request
func TrackerAnnounce(passkey string, query map[string]string, resChan chan []byte) {
	// Store announce information in struct
	announce := MapToAnnounceLog(query, resChan)

	// Request to store announce
	_, ok := DbWrite(announce.InfoHash, announce)
	if !ok {
		Static.LogChan <- "could not write struct from storage"
	}

	Static.LogChan <- fmt.Sprintf("db: [ip: %s, port:%d]", announce.Ip, announce.Port)
	Static.LogChan <- fmt.Sprintf("db: [info_hash: %s]", announce.InfoHash)
	Static.LogChan <- fmt.Sprintf("db: [peer_id: %s]", announce.PeerId)

	// Fetch peer information

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

	// Fake tracker announce response
	announceRes := FakeAnnounceResponse(numwant)
	Static.LogChan <- fmt.Sprintf("res: %s", announceRes)
	resChan <- announceRes
}

// Report a bencoded []byte response as specified by input string
func TrackerError(resChan chan []byte, err string) {
	resChan <- bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"interval":       bencode.EncInt(RandRange(3200, 4000)),
		"min interval":   bencode.EncInt(1800),
	})
}
