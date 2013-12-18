package goat

import (
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

// Generate an AnnounceLog struct from a query map
func MapToAnnounceLog(query map[string]string) AnnounceLog {
	var announce AnnounceLog

	// Required parameters

	// info_hash
	announce.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// ip
	announce.Ip = query["ip"]

	// Note: integers previously validated in connHandler.go

	// port
	port, _ := strconv.Atoi(query["port"])
	announce.Port = port

	// uploaded
	uploaded, _ := strconv.ParseInt(query["uploaded"], 10, 64)
	announce.Uploaded = uploaded

	// downloaded
	downloaded, _ := strconv.ParseInt(query["downloaded"], 10, 64)
	announce.Downloaded = downloaded

	// left
	left, _ := strconv.ParseInt(query["left"], 10, 64)
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
	return min + rand.Intn(max-min)
}
