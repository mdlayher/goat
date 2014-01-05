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

	// Store announce log
	query := "INSERT INTO announce_log " +
		"(`info_hash`, `passkey`, `key`, `ip`, `port`, `udp`, `uploaded`, `downloaded`, `left`, `event`, `client`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.Passkey, a.Key, a.IP, a.Port, a.UDP, a.Uploaded, a.Downloaded, a.Left, a.Event, a.Client)
	tx.Commit()

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(id interface{}, col string) AnnounceLog {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return a
	}

	// Fetch announce log into struct
	a = AnnounceLog{}
	db.Get(&a, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)

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

	// Note: integers previously validated

	// port
	port, _ := strconv.Atoi(query["port"])
	a.Port = port

	// udp
	if query["udp"] == "1" {
		a.UDP = true
	} else {
		a.UDP = false
	}

	// uploaded
	uploaded, _ := strconv.ParseInt(query["uploaded"], 10, 64)
	a.Uploaded = uploaded

	// downloaded
	downloaded, _ := strconv.ParseInt(query["downloaded"], 10, 64)
	a.Downloaded = downloaded

	// left
	left, _ := strconv.ParseInt(query["left"], 10, 64)
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
