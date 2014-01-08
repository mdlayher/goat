package goat

import (
	"database/sql"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
)

// fileRecord represents a file tracked by tracker
type fileRecord struct {
	ID         int    `json:"id"`
	InfoHash   string `db:"info_hash" json:"infoHash"`
	Verified   bool   `json:"verified"`
	CreateTime int64  `db:"create_time" json:"createTime"`
	UpdateTime int64  `db:"update_time" json:"updateTime"`
}

// jsonFileRecord represents output fileRecord JSON for API
type jsonFileRecord struct {
	ID         int              `json:"id"`
	InfoHash   string           `json:"infoHash"`
	Verified   bool             `json:"verified"`
	CreateTime int64            `json:"createTime"`
	UpdateTime int64            `json:"updateTime"`
	FileUsers  []fileUserRecord `json:"fileUsers"`
}

// ToJSON converts a fileRecord to a jsonFileRecord struct
func (f fileRecord) ToJSON() []byte {
	// Convert all standard fields to the JSON equivalent struct
	j := jsonFileRecord{}
	j.ID = f.ID
	j.InfoHash = f.InfoHash
	j.Verified = f.Verified
	j.CreateTime = f.CreateTime
	j.UpdateTime = f.UpdateTime

	// Load in fileUserRecords associated with this file
	j.FileUsers = f.Users()

	// Marshal into JSON
	out, err := json.Marshal(j)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return out
}

// fileRecordRepository is used to contain methods to load multiple fileRecord structs
type fileRecordRepository struct {
}

// Save fileRecord to storage
func (f fileRecord) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Store or update file information
	query := "INSERT INTO files " +
		"(`info_hash`, `verified`, `create_time`, `update_time`) " +
		"VALUES (?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()) " +
		"ON DUPLICATE KEY UPDATE " +
		"`verified`=values(`verified`), `update_time`=UNIX_TIMESTAMP();"

	// Create database transaction, do insert, commit
	tx := db.MustBegin()
	tx.Execl(query, f.InfoHash, f.Verified)
	tx.Commit()

	return true
}

// Load fileRecord from storage
func (f fileRecord) Load(id interface{}, col string) fileRecord {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return f
	}

	// Fetch announce log into struct
	f = fileRecord{}
	err = db.Get(&f, "SELECT * FROM files WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return fileRecord{}
	}

	return f
}

// Completed returns the number of completions, active or not, on this file
func (f fileRecord) Completed() int {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Completeds struct
	completed := struct {
		Completed int
	}{
		0,
	}

	// Calculate number of completions on this file, defined as users who are completed, and 0 left
	err = db.Get(&completed, "SELECT COUNT(user_id) AS completed FROM files_users WHERE file_id = ? AND completed = 1 AND `left` = 0;", f.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return completed.Completed
}

// Seeders returns the number of seeders on this file
func (f fileRecord) Seeders() int {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Seeders struct
	seeders := struct {
		Seeders int
	}{
		0,
	}

	// Calculate number of seeders on this file, defined as users who are active, completed, and 0 left
	err = db.Get(&seeders, "SELECT COUNT(user_id) AS seeders FROM files_users WHERE file_id = ? AND active = 1 AND completed = 1 AND `left` = 0;", f.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return seeders.Seeders
}

// Leechers returns the number of leechers on this file
func (f fileRecord) Leechers() int {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Anonymous Leechers struct
	leechers := struct {
		Leechers int
	}{
		0,
	}

	// Calculate number of leechers on this file, defined as users who are active, completed, and positive bytes left
	db.Get(&leechers, "SELECT COUNT(user_id) AS leechers FROM files_users WHERE file_id = ? AND active = 1 AND completed = 0 AND `left` > 0;", f.ID)
	if err != nil {
		log.Println(err.Error())
		return -1
	}

	return leechers.Leechers
}

// PeerList returns the compact peer buffer for tracker announce, excluding self
func (f fileRecord) PeerList(exclude string, numwant int) []byte {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Anonymous Peer struct
	peer := struct {
		IP   string
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
		"AND announce_log.ip = files_users.ip " +
		"WHERE files_users.active=1 " +
		"AND files.info_hash=? " +
		"AND announce_log.ip != ? " +
		"LIMIT ?;"

	rows, err := db.Queryx(query, f.InfoHash, exclude, numwant)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	// Iterate all rows
	for rows.Next() {
		// Scan row results
		rows.StructScan(&peer)

		// Parse IP into byte buffer
		ip := [4]byte{}
		binary.BigEndian.PutUint32(ip[:], binary.BigEndian.Uint32(net.ParseIP(peer.IP).To4()))

		// Parse port into byte buffer
		port := [2]byte{}
		binary.BigEndian.PutUint16(port[:], peer.Port)

		// Append ip/port to end of list
		buf = append(buf[:], append(ip[:], port[:]...)...)
	}

	return buf
}

// PeerReaper reaps peers who have not recently announced on this torrent, and mark them inactive
func (f fileRecord) PeerReaper() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Anonymous result struct
	result := struct {
		UserID int `db:"user_id"`
		IP     string
	}{
		0,
		"",
	}

	// Query for user IDs associated with this file, who are marked active but have not announced recently
	query := "SELECT user_id, ip FROM files_users " +
		"WHERE time < (UNIX_TIMESTAMP() - ?) " +
		"AND active = 1 " +
		"AND file_id = ?;"

	rows, err := db.Queryx(query, static.Config.Interval+60, f.ID)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Count of reaped peers
	count := 0

	// Iterate all rows
	tx := db.MustBegin()
	for rows.Next() {
		// Scan row results
		rows.StructScan(&result)

		// Mark peer as inactive
		reapQuery := "UPDATE files_users SET active = 0 WHERE file_id = ? AND user_id = ? AND ip = ?;"
		tx.Execl(reapQuery, f.ID, result.UserID, result.IP)

		count++
	}

	// Commit transaction
	tx.Commit()

	if count > 0 {
		log.Printf("reaper: reaped %d peer(s) on file %d\n", count, f.ID)
	}

	return true
}

// Users loads all fileUserRecord structs associated with this fileRecord struct
func (f fileRecord) Users() []fileUserRecord {
	return new(fileUserRecordRepository).Select(f.ID, "file_id")
}

// All loads all fileRecord structs from storage
func (f fileRecordRepository) All() []fileRecord {
	files := make([]fileRecord, 0)

	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return files
	}

	// Load all files
	rows, err := db.Queryx("SELECT * FROM files")
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return files
	}

	// Iterate all rows and build array
	file := fileRecord{}
	for rows.Next() {
		// Scan row results
		rows.StructScan(&file)
		files = append(files[:], file)
	}

	return files
}
