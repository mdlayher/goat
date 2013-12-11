package goat

import (
	"bencode"
)

// Tracker announce request
func TrackerAnnounce(resChan chan []byte) {
	resChan <- bencode.EncString("announce successful")
}

// Report a bencoded []byte response as specified by input string
func TrackerError(resChan chan []byte, err string) {
	resChan <- bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"min interval":   bencode.EncInt(1800),
		"interval":       bencode.EncInt(3600),
	})
}
