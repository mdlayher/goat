package goat

import (
	"bencode"
	"encoding/binary"
	"encoding/hex"
	"math/rand"
	"net"
	"strconv"
	"time"
)

// Generate an AnnounceLog struct from a query map
func MapToAnnounceLog(query map[string]string, resChan chan []byte) AnnounceLog {
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

// Generate a random announce interval in the specified range
func RandRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max - min)
}

// Generate a fake announce response
func FakeAnnounceResponse(numwant int) []byte {
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
		"interval":     bencode.EncInt(RandRange(3200, 4000)),
		"min interval": bencode.EncInt(1800),
		"peers":        bencode.EncBytes(CompactPeerList()),
	}

	return bencode.EncDictMap(res)
}

// Save for later: Generate a compact peer list in binary format
func CompactPeerList() []byte {
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
