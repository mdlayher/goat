package goat

import (
	"crypto/sha1"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Generates a unique hash for a given struct with its ID, to be used as ID for storage
type Hashable interface {
	Hash() string
}

// Defines methods to save, load, and delete data for an object
type Persistent interface {
	Save() bool
	Load(interface{}, string) interface{}
	Delete() bool
}

type ErrorRes struct {
	ErrLocation string
	Time        int64
	Error       string
}

// Struct representing an announce, to be logged to storage
type AnnounceLog struct {
	Id         int
	InfoHash   string `db:"info_hash"`
	PeerId     string `db:"peer_id"`
	Ip         string
	Port       int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Event      string
	Time       int64
}

// Generate a SHA1 hash of the form: announce_log_InfoHash
func (log AnnounceLog) Hash() string {
	hash := sha1.New()
	hash.Write([]byte("announce_log"+log.InfoHash))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Save AnnounceLog to storage
func (a AnnounceLog) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl("INSERT INTO announce_log (`info_hash`, `peer_id`, `ip`, `port`, `uploaded`, `downloaded`, `left`, `event`, `time`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());",
		a.InfoHash, a.PeerId, a.Ip, a.Port, a.Uploaded, a.Downloaded, a.Left, a.Event)
	tx.Commit()

	return true
}

// Load AnnounceLog from storage
func (a AnnounceLog) Load(id interface{}, col string) AnnounceLog {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return a
	}

	// Fetch announce log into struct
	a = AnnounceLog{}
	db.Get(&a, "SELECT * FROM announce_log WHERE `"+col+"`=?", id)

	return a
}

// Struct representing a file tracked by tracker
type FileRecord struct {
	Id         int
	InfoHash   string `db:"info_hash"`
	Leechers   int
	Seeders    int
	Completed  int
	CreateTime int64 `db:"create_time"`
	UpdateTime int64 `db:"update_time"`
}

// Save FileRecord to storage
func (f FileRecord) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl("INSERT INTO files (`info_hash`, `leechers`, `seeders`, `completed`, `create_time`, `update_time`) VALUES (?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()) ON DUPLICATE KEY UPDATE `leechers`=values(`leechers`), `seeders`=values(`seeders`), `completed`=values(`completed`), `update_time`=UNIX_TIMESTAMP();",
		f.InfoHash, f.Leechers, f.Seeders, f.Completed)
	tx.Commit()

	return true
}

// Load FileRecord from storage
func (f FileRecord) Load(id interface{}, col string) FileRecord {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return f
	}

	// Fetch announce log into struct
	f = FileRecord{}
	db.Get(&f, "SELECT * FROM files WHERE `"+col+"`=?", id)
	return f
}

// Struct representing a user on the tracker
type UserRecord struct {
	Id           int
	Username     string
	Passkey      string
	TorrentLimit int `db:"torrent_limit"`
	Uploaded     int64
	Downloaded   int64
}

// Save UserRecord to storage
func (u UserRecord) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl("INSERT INTO users (`username`, `passkey`, `torrent_limit`, `uploaded`, `downloaded`) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE `username`=values(`username`), `passkey`=values(`passkey`), `torrent_limit`=values(`torrent_limit`), `uploaded`=values(`uploaded`), `downloaded`=values(`downloaded)`;",
		u.Username, u.Passkey, u.TorrentLimit, u.Uploaded, u.Downloaded)
	tx.Commit()

	return true
}

// Load UserRecord from storage
func (u UserRecord) Load(id interface{}, col string) UserRecord {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return u
	}

	// Fetch announce log into struct
	u = UserRecord{}
	db.Get(&u, "SELECT * FROM users WHERE `"+col+"`=?", id)
	return u
}

// Struct representing a scrapelog, to be logged to storage
type ScrapeLog struct {
	InfoHash string
	UserId   string
	Ip       string
	Time     int64
}
