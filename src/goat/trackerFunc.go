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
func MapToAnnounceLog(query map[string]string) AnnounceLog {
	var announce AnnounceLog

	// Required parameters

	// info_hash
	announce.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// peer_id
	announce.PeerId = hex.EncodeToString([]byte(query["peer_id"]))

	// ip
	announce.Ip = query["ip"]

	// Note: integers previously validated in connHandler.go

	// port
	port, _ := strconv.Atoi(query["port"])
	announce.Port = port

	// uploaded
	uploaded, _ := strconv.Atoi(query["uploaded"])
	announce.Uploaded = uploaded

	// downloaded
	downloaded, _ := strconv.Atoi(query["downloaded"])
	announce.Downloaded = downloaded

	// left
	left, _ := strconv.Atoi(query["left"])
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
