package goat

import (
	"bencode"
	"encoding/binary"
	"encoding/json"
	"net"
	"os"
	"runtime"
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

	// Don't use compact announce by default
	compact := false
	if _, ok := query["compact"]; ok && query["compact"][0] == "1" {
		compact = true
	}

	// Fake tracker announce response
	resChan <- fakeAnnounceResponse(compact)
}

type ServerStatus struct {
	Pid int
	Hostname string
	Platform string
	Architecture string
	NumCpu int
	NumGoroutine int
	/*
	CpuPercent float32,
	MemoryMb float32,
	*/
}

// Tracker status request
func TrackerStatus(resChan chan []byte) {
	// Get system hostname
	hostname, _ := os.Hostname()

	res, err := json.Marshal(ServerStatus{
		os.Getpid(),
		hostname,
		runtime.GOOS,
		runtime.GOARCH,
		runtime.NumCPU(),
		runtime.NumGoroutine(),
	})
	if err != nil {
		resChan <- nil
	}

	resChan <- res
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
func fakeAnnounceResponse(compact bool) []byte {
	res := map[string][]byte{
		// complete, downloaded, incomplete: used to report client download statistics
		// min interval, interval: used to tell clients how often to announce
		"complete":     bencode.EncInt(99),
		"downloaded":   bencode.EncInt(200),
		"incomplete":   bencode.EncInt(1),
		"min interval": bencode.EncInt(1800),
		"interval":     bencode.EncInt(3600),
	}

	// Send a compact response if requested
	if compact {
		res["peers"] = bencode.EncBytes(compactPeerList())
	} else {
		// peers: list of dictionaries
		res["peers"] = bencode.EncList([][]byte{
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
		})
	}

	return bencode.EncDictMap(res)
}

// Save for later: Generate a compact peer list in binary format
func compactPeerList() []byte {
	// Empty byte buffer
	buf := []byte("")

	// Add a bunch of fake peers to list
	for i := 0; i < 5; i++ {

		// Compact peers into binary format: ip ip ip ip port port
		ip := [4]byte{}
		binary.BigEndian.PutUint32(ip[:], binary.BigEndian.Uint32(net.ParseIP("255.255.255.255").To4()))
		port := [2]byte{}
		binary.BigEndian.PutUint16(port[:], 8080)

		// Append to end of buffer
		buf = append(buf[:], append(ip[:], port[:]...)...)
	}

	return buf
}
