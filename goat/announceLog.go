package goat

import (
	"encoding/hex"
	"strconv"
	"time"
)

// AnnounceLog represents an announce, to be logged to storage
type AnnounceLog struct {
	ID         int
	InfoHash   string `db:"info_hash"`
	Passkey    string
	Key        string
	IP         string
	Port       int
	UDP        bool
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	Client     string
	Time       int64
}

// Save AnnounceLog to storage
func (a AnnounceLog) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	if err := db.SaveAnnounceLog(a); nil != err {
		Static.LogChan <- err.Error()
		return false
	}

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(id interface{}, col string) AnnounceLog {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return AnnounceLog{}
	}

	if a, err = db.LoadAnnounceLog(id, col); nil != err {
		Static.LogChan <- err.Error()
	}

	return a
}

// FromMap generates an AnnounceLog struct from a query map
func (a AnnounceLog) FromMap(query map[string]string) AnnounceLog {
	a = AnnounceLog{}

	// Required parameters

	// info_hash
	a.InfoHash = hex.EncodeToString([]byte(query["info_hash"]))

	// passkey
	a.Passkey = query["passkey"]

	// key
	a.Key = query["key"]

	// ip
	a.IP = query["ip"]

	// udp
	if query["udp"] == "1" {
		a.UDP = true
	} else {
		a.UDP = false
	}

	// port
	port, err := strconv.Atoi(query["port"])
	if err != nil {
		Static.LogChan <- err.Error()
		return AnnounceLog{}
	}
	a.Port = port

	// uploaded
	uploaded, err := strconv.ParseInt(query["uploaded"], 10, 64)
	if err != nil {
		Static.LogChan <- err.Error()
		return AnnounceLog{}
	}
	a.Uploaded = uploaded

	// downloaded
	downloaded, err := strconv.ParseInt(query["downloaded"], 10, 64)
	if err != nil {
		Static.LogChan <- err.Error()
		return AnnounceLog{}
	}
	a.Downloaded = downloaded

	// left
	left, err := strconv.ParseInt(query["left"], 10, 64)
	if err != nil {
		Static.LogChan <- err.Error()
		return AnnounceLog{}
	}
	a.Left = left

	// Optional parameters

	// event
	if event, ok := query["event"]; ok {
		a.Event = event
	}

	// BitTorrent client, User-Agent header
	a.Client = query["client"]

	// Current UNIX timestamp
	a.Time = time.Now().Unix()

	// Return the created announce
	return a
}
