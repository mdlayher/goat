package data

import (
	"encoding/hex"
	"errors"
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
func (a AnnounceLog) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save AnnounceLog
	if err := db.SaveAnnounceLog(a); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(ID interface{}, col string) (AnnounceLog, error) {
	a = AnnounceLog{}

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return a, err
	}

	// Load AnnounceLog
	a, err = db.LoadAnnounceLog(ID, col)
	if err != nil {
		return a, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return a, err
	}

	return a, err
}

// Delete AnnounceLog from storage
func (a AnnounceLog) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete AnnounceLog
	if err = db.DeleteAnnounceLog(a.ID, "id"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
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
