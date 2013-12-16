package goat

import (
	"bencode"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net"
	"strconv"
	"time"
)

// Tracker announce request
func TrackerAnnounce(passkey string, query map[string]string, resChan chan []byte) {
	// Store announce information in struct
	announce := mapToAnnounceLog(query, resChan)

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
	announceRes := fakeAnnounceResponse(numwant)
	Static.LogChan <- fmt.Sprintf("res: %s", announceRes)
	resChan <- announceRes
}

// Report a bencoded []byte response as specified by input string
func TrackerError(resChan chan []byte, err string) {
	resChan <- bencode.EncDictMap(map[string][]byte{
		"failure reason": bencode.EncString(err),
		"interval":       bencode.EncInt(3600),
		"min interval":   bencode.EncInt(1800),
	})
}

// Generate an AnnounceLog struct from a query map
func mapToAnnounceLog(query map[string]string, resChan chan []byte) AnnounceLog {
	var announce AnnounceLog

	// Required parameters

	// info_hash
	infoHash := make([]byte, 64)
	hex.Encode(infoHash, []byte(query["info_hash"]))
	announce.InfoHash = string(infoHash)

	// peer_id
	peerId := make([]byte, 64)
	hex.Encode(peerId, []byte(query["peer_id"]))
	announce.PeerId = string(peerId)

	// ip
	announce.Ip = query["ip"]

	// port
	port, err := strconv.Atoi(query["port"])
	if err != nil {
		TrackerError(resChan, "parameter port is not a valid integer")
	}
	announce.Port = port

	// uploaded
	uploaded, err := strconv.Atoi(query["uploaded"])
	if err != nil {
		TrackerError(resChan, "parameter uploaded is not a valid integer")
	}
	announce.Uploaded = uploaded

	// downloaded
	downloaded, err := strconv.Atoi(query["downloaded"])
	if err != nil {
		TrackerError(resChan, "parameter downloaded is not a valid integer")
	}
	announce.Downloaded = downloaded

	// left
	left, err := strconv.Atoi(query["left"])
	if err != nil {
		TrackerError(resChan, "parameter left is not a valid integer")
	}
	announce.Left = left

	// Optional parameters

	// event
	if event, ok := query["event"]; ok {
		announce.Event = event
	}

	// Current UNIX timestamp
	announce.Time = time.Now().Unix()

	// Return the created announce
	return announce
}

// Generate a fake announce response
func fakeAnnounceResponse(numwant int) []byte {
	// For now, we completely ignore numwant
	_ = numwant

	res := map[string][]byte{
		// complete, downloaded, incomplete: used to report client download statistics
		// min interval, interval: used to tell clients how often to announce
		// TODO: these may not be necessary, or may cause problems
		// Need to do further research
		/*
		"complete":     bencode.EncInt(0),
		"downloaded":   bencode.EncInt(0),
		"incomplete":   bencode.EncInt(0),
		*/
		"interval":     bencode.EncInt(3600),
		"min interval": bencode.EncInt(1800),
		"peers":        bencode.EncBytes(compactPeerList()),
	}

	return bencode.EncDictMap(res)
}

// Save for later: Generate a compact peer list in binary format
func compactPeerList() []byte {
	// Empty byte buffer
	buf := []byte("")

	// Add a bunch of fake peers to list
	for i := 0; i < 1; i++ {
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
