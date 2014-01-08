package goat

import (
	"database/sql"
	"encoding/hex"
	"log"
	"strconv"
	"time"
)

// announceLog represents an announce, to be logged to storage
type announceLog struct {
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

// Save announceLog to storage
func (a announceLog) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
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

// Load announceLog from storage
func (a announceLog) Load(ID interface{}, col string) announceLog {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}

	// Fetch announce log into struct
	a = announceLog{}
	err = db.Get(&a, "SELECT * FROM announce_log WHERE `"+col+"`=?", ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return announceLog{}
	}

	return a
}

// FromMap generates an announceLog struct from a query map
func (a announceLog) FromMap(query map[string]string) announceLog {
	a = announceLog{}

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
		log.Println(err.Error())
		return announceLog{}
	}
	a.Port = port

	// uploaded
	uploaded, err := strconv.ParseInt(query["uploaded"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Uploaded = uploaded

	// downloaded
	downloaded, err := strconv.ParseInt(query["downloaded"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Downloaded = downloaded

	// left
	left, err := strconv.ParseInt(query["left"], 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
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
