package goat

import (
	"encoding/hex"
	"log"
	"net/url"
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

	if err := db.SaveAnnounceLog(a); err != nil {
		log.Println(err.Error())
		return false
	}

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
	if a, err = db.LoadAnnounceLog(ID, col); err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	return a
}

// FromValues generates an announceLog struct from a url.Values map
func (a announceLog) FromValues(query url.Values) announceLog {
	a = announceLog{}

	// Required parameters

	// info_hash
	a.InfoHash = hex.EncodeToString([]byte(query.Get("info_hash")))

	// passkey
	a.Passkey = query.Get("passkey")

	// key
	a.Key = query.Get("key")

	// ip
	a.IP = query.Get("ip")

	// udp
	if query.Get("udp") == "1" {
		a.UDP = true
	} else {
		a.UDP = false
	}

	// port
	port, err := strconv.Atoi(query.Get("port"))
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Port = port

	// uploaded
	uploaded, err := strconv.ParseInt(query.Get("uploaded"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Uploaded = uploaded

	// downloaded
	downloaded, err := strconv.ParseInt(query.Get("downloaded"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Downloaded = downloaded

	// left
	left, err := strconv.ParseInt(query.Get("left"), 10, 64)
	if err != nil {
		log.Println(err.Error())
		return announceLog{}
	}
	a.Left = left

	// Optional parameters

	// event
	if query.Get("event") != "" {
		a.Event = query.Get("event")
	}

	// BitTorrent client, User-Agent header
	a.Client = query.Get("client")

	// Current UNIX timestamp
	a.Time = time.Now().Unix()

	// Return the created announce
	return a
}
