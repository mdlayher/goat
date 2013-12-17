package goat

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net"
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
	hash.Write([]byte("announce_log" + log.InfoHash))
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

	// Store announce log
	query := "INSERT INTO announce_log " +
	"(`info_hash`, `peer_id`, `ip`, `port`, `uploaded`, `downloaded`, `left`, `event`, `time`) " +
	"VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.PeerId, a.Ip, a.Port, a.Uploaded, a.Downloaded, a.Left, a.Event)
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
	Verified   bool
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

	// Store or update file information
	query := "INSERT INTO files " +
	"(`info_hash`, `verified`, `leechers`, `seeders`, `completed`, `create_time`, `update_time`) " +
	"VALUES (?, ?, ?, ?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()) " +
	"ON DUPLICATE KEY UPDATE " +
	"`verified`=values(`verified`), `leechers`=values(`leechers`), `seeders`=values(`seeders`), " +
	"`completed`=values(`completed`), `update_time`=UNIX_TIMESTAMP();"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, f.InfoHash, f.Verified, f.Leechers, f.Seeders, f.Completed)
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

// Return compact peer buffer for tracker announce, excluding self
func (f FileRecord) PeerList(exclude string, numwant int) []byte {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return nil
	}

	// Anonymous Peer struct
	peer := struct {
		Ip   string
		Port uint16
	}{
		"",
		0,
	}

	// Buffer for compact list
	buf := make([]byte, 0)

	// Get IP and port of all peers who are active and seeding this file
	query := "SELECT DISTINCT announce_log.ip,announce_log.port FROM announce_log " +
	"JOIN files ON announce_log.info_hash = files.info_hash " +
	"JOIN files_users ON files.id = files_users.file_id " +
	"WHERE files_users.active=1 " +
	"AND files.info_hash=? " +
	"AND announce_log.ip != ? " +
	"LIMIT ?;"

	rows, err := db.Queryx(query, f.InfoHash, exclude, numwant)
	if err != nil {
		Static.LogChan <- err.Error()
		return buf
	}

	// Iterate all rows
	for rows.Next() {
		// Scan row results
		rows.StructScan(&peer)

		// Parse IP into byte buffer
		ip := [4]byte{}
		binary.BigEndian.PutUint32(ip[:], binary.BigEndian.Uint32(net.ParseIP(peer.Ip).To4()))

		// Parse port into byte buffer
		port := [2]byte{}
		binary.BigEndian.PutUint16(port[:], peer.Port)

		// Append ip/port to end of list
		buf = append(buf[:], append(ip[:], port[:]...)...)
	}

	return buf
}

// Struct representing a file tracked by tracker
type FileUserRecord struct {
	FileId     int `db:"file_id"`
	UserId     int `db:"user_id"`
	Active     bool
	Completed  bool
	Announced  int
	Uploaded   int64
	Downloaded int64
	Left       int64
	Time       int64
}

// Save FileUserRecord to storage
func (f FileUserRecord) Save() bool {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	// Insert or update a file/user relationship record
	query := "INSERT INTO files_users " +
	"(`file_id`, `user_id`, `active`, `completed`, `announced`, `uploaded`, `downloaded`, `left`, `time`) " +
	"VALUES (?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP()) " +
	"ON DUPLICATE KEY UPDATE " +
	"`active`=values(`active`), `completed`=values(`completed`), `announced`=values(`announced`), " +
	"`uploaded`=values(`uploaded`), `downloaded`=values(`downloaded`), `left`=values(`left`), " +
	"`time`=UNIX_TIMESTAMP();"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, f.FileId, f.UserId, f.Active, f.Completed, f.Announced, f.Uploaded, f.Downloaded, f.Left)
	tx.Commit()

	return true
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileId interface{}, userId interface{}) FileUserRecord {
	// Open database connection
	db, err := DbConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return f
	}

	// Fetch announce log into struct
	f = FileUserRecord{}
	db.Get(&f, "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=?", fileId, userId)
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

	// Calculate total uploaded and downloaded
	if u.Id != 0 {
		totals := struct {
			Uploaded   int64
			Downloaded int64
		}{
			0,
			0,
		}

		// Get sum of upload and download for this user
		totalQuery := "SELECT SUM(uploaded) AS uploaded, " +
		"SUM(downloaded) AS downloaded " +
		"FROM files_users " +
		"WHERE user_id = ?"

		err = db.Get(&totals, totalQuery, u.Id)
		if err == nil {
			// Store in struct
			u.Uploaded = totals.Uploaded
			u.Downloaded = totals.Downloaded
		}
	}

	// Insert or update a user record
	query := "INSERT INTO users " +
	"(`username`, `passkey`, `torrent_limit`, `uploaded`, `downloaded`) " +
	"VALUES (?, ?, ?, ?, ?) " +
	"ON DUPLICATE KEY UPDATE " +
	"`username`=values(`username`), `passkey`=values(`passkey`), `torrent_limit`=values(`torrent_limit`), " +
	"`uploaded`=values(`uploaded`), `downloaded`=values(`downloaded`);"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, u.Username, u.Passkey, u.TorrentLimit, u.Uploaded, u.Downloaded)
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
