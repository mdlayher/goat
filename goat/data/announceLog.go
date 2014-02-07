package data

import (
	"encoding/hex"
	"errors"
	"log"
	"net/url"
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
		log.Println(err.Error())
		return false
	}

	// Save AnnounceLog
	if err := db.SaveAnnounceLog(a); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(ID interface{}, col string) AnnounceLog {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return AnnounceLog{}
	}

	// Load AnnounceLog
	if a, err = db.LoadAnnounceLog(ID, col); err != nil {
		log.Println(err.Error())
		return AnnounceLog{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return a
}

// Delete AnnounceLog from storage
func (a AnnounceLog) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete AnnounceLog
	if err = db.DeleteAnnounceLog(a.ID, "id"); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// FromValues generates an AnnounceLog struct from a url.Values map
func (a *AnnounceLog) FromValues(query url.Values) error {
	// Required parameters

	// info_hash (20 characters, 40 characters after hex encode)
	a.InfoHash = hex.EncodeToString([]byte(query.Get("info_hash")))
	if len(a.InfoHash) != 40 {
		return errors.New("info_hash must be exactly 20 characters")
	}

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
		return errors.New("invalid integer parameter: port")
	}
	a.Port = port

	// uploaded
	uploaded, err := strconv.ParseInt(query.Get("uploaded"), 10, 64)
	if err != nil {
		return errors.New("invalid integer parameter: uploaded")
	}
	a.Uploaded = uploaded

	// downloaded
	downloaded, err := strconv.ParseInt(query.Get("downloaded"), 10, 64)
	if err != nil {
		return errors.New("invalid integer parameter: downloaded")
	}
	a.Downloaded = downloaded

	// left
	left, err := strconv.ParseInt(query.Get("left"), 10, 64)
	if err != nil {
		return errors.New("invalid integer parameter: left")
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

	return nil
}
