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

	// passkey
	announce.Passkey = query["passkey"]

	// key
	announce.Key = query["key"]

	// ip
	announce.Ip = query["ip"]

	// Note: integers previously validated in connHandler.go

	// port
	port, _ := strconv.Atoi(query["port"])
	announce.Port = port

	// udp
	if query["udp"] == "1" {
		announce.Udp = true
	} else {
		announce.Udp = false
	}

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

	// BitTorrent client, User-Agent header
	announce.Client = query["client"]

	// Current UNIX timestamp
	announce.Time = time.Now().Unix()

	// Return the created announce
	return announce
}

// Generate a ScrapeLog struct from a query map
func MapToScrapeLog(query map[string]string) ScrapeLog {
	var scrape ScrapeLog

	// Required parameters

	// info_hash
	scrape.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// passkey
	scrape.Passkey = query["passkey"]

	// ip
	scrape.Ip = query["ip"]

	// Current UNIX timestamp
	scrape.Time = time.Now().Unix()

	// Return the created scrape
	return scrape
}

// Generate a random announce interval in the specified range
func RandRange(min int, max int) int {
	rand.Seed(time.Now().Unix())
	return min + rand.Intn(max-min)
}
