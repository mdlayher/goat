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

	// Store announce information

	// Fetch peer information

	// Fake tracker announce response
	resChan <- fakeAnnounceResponse()
}

// Report a bencoded []byte response as specified by input string
func TrackerError(resChan chan []byte, err string) {
	resChan <- bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"min interval":   bencode.EncInt(1800),
		"interval":       bencode.EncInt(3600),
	})
}

// Generate a fake announce response
func fakeAnnounceResponse() []byte {
	return bencode.EncDictMap(map[string][]byte{
		// min interval, interval: used to tell clients how often to announce
		"min interval": bencode.EncInt(1800),
		"interval":     bencode.EncInt(3600),
		// peers: list of dictionaries
		"peers": bencode.EncList([][]byte{
			// peer dictionary: contains peer id, ip, and port of a peer
			bencode.EncDictMap(map[string][]byte{
				"peer id": bencode.EncString("ABCDEF0123456789ABCD"),
				"ip":      bencode.EncString("127.0.0.1"),
				"port":    bencode.EncInt(6881),
			}),
			bencode.EncDictMap(map[string][]byte{
				"peer id": bencode.EncString("0123456789ABCDEF0123"),
				"ip":      bencode.EncString("127.0.0.1"),
				"port":    bencode.EncInt(6882),
			}),
		}),
	})
}
