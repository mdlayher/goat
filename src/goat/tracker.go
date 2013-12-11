package goat

import (
	"bencode"
)

// Tracker announce request
func TrackerAnnounce(query map[string][]string, resChan chan []byte) {
	// Validate required parameter input
	required := []string{"info_hash", "peer_id", "ip", "port", "uploaded", "downloaded", "left"}
	for _, r := range required {
		if _, ok := query[r]; !ok {
			TrackerError(resChan, "missing required parameter: "+r)
		}
	}

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
