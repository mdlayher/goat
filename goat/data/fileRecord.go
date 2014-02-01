package data

import (
	"encoding/json"
	"log"
	"time"

	"github.com/mdlayher/goat/goat/common"
)

// FileRecord represents a file tracked by tracker
type FileRecord struct {
	ID         int    `json:"id"`
	InfoHash   string `db:"info_hash" json:"infoHash"`
	Verified   bool   `json:"verified"`
	CreateTime int64  `db:"create_time" json:"createTime"`
	UpdateTime int64  `db:"update_time" json:"updateTime"`
}

// jsonFileRecord represents output FileRecord JSON for API
type jsonFileRecord struct {
	ID         int              `json:"id"`
	InfoHash   string           `json:"infoHash"`
	Verified   bool             `json:"verified"`
	CreateTime int64            `json:"createTime"`
	UpdateTime int64            `json:"updateTime"`
	Completed  int              `json:"completed"`
	Seeders    int              `json:"seeders"`
	Leechers   int              `json:"leechers"`
	FileUsers  []FileUserRecord `json:"fileUsers"`
}

// peerInfo represents a peer which will be marked as active or not
type peerInfo struct {
	UserID int `db:"user_id"`
	IP     string
}

// ToJSON converts a FileRecord to a jsonFileRecord struct
func (f FileRecord) ToJSON() []byte {
	// Convert all standard fields to the JSON equivalent struct
	j := jsonFileRecord{}
	j.ID = f.ID
	j.InfoHash = f.InfoHash
	j.Verified = f.Verified
	j.CreateTime = f.CreateTime
	j.UpdateTime = f.UpdateTime

	// Load in FileUserRecords associated with this file
	j.FileUsers = f.Users()

	// Load counts for completions, seeding, leeching
	j.Completed = f.Completed()
	j.Seeders = f.Seeders()
	j.Leechers = f.Leechers()

	// Marshal into JSON
	out, err := json.Marshal(j)
	if err != nil {
		log.Println(err.Error())
		return nil
	}

	return out
}

// FileRecordRepository is used to contain methods to load multiple FileRecord structs
type FileRecordRepository struct {
}

// Delete FileRecord from storage
func (f FileRecord) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete FileRecord
	if err = db.DeleteFileRecord(f.ID, "id"); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Save FileRecord to storage
func (f FileRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save FileRecord
	if err := db.SaveFileRecord(f); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load FileRecord from storage
func (f FileRecord) Load(id interface{}, col string) FileRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return f
	}

	// Load FileRecord by column
	if f, err = db.LoadFileRecord(id, col); err != nil {
		log.Println(err.Error())
		return FileRecord{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return f
}

// Completed returns the number of completions, active or not, on this file
func (f FileRecord) Completed() (completed int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Retrieve number of file completions
	if completed, err = db.CountFileRecordCompleted(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}

// Seeders returns the number of seeders on this file
func (f FileRecord) Seeders() (seeders int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Return number of active seeders
	if seeders, err = db.CountFileRecordSeeders(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}

// Leechers returns the number of leechers on this file
func (f FileRecord) Leechers() (leechers int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}

	// Return number of active leechers
	if leechers, err = db.CountFileRecordLeechers(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}

// PeerList returns the compact peer buffer for tracker announce, excluding self
func (f FileRecord) PeerList(exclude string, numwant int) (peers []byte) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Return compact peer list, excluding this IP
	if peers, err = db.GetFileRecordPeerList(f.InfoHash, exclude, numwant); err != nil {
		log.Println(err.Error())
		return nil
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}

// PeerReaper reaps peers who have not recently announced on this torrent, and mark them inactive
func (f FileRecord) PeerReaper() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Retrieve list of inactive users (have not announced in just above maximum interval)
	users, err := db.GetInactiveUserInfo(f.ID, time.Duration(int64(common.Static.Config.Interval))*time.Second+60)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Mark those users inactive on this file
	if err := db.MarkFileUsersInactive(f.ID, users); err != nil {
		log.Println(err.Error())
		return false
	}

	// Print number of peers reaped
	if count := len(users); count > 0 {
		log.Printf("reaper: reaped %d peer(s) on file %d\n", count, f.ID)
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Users loads all FileUserRecord structs associated with this FileRecord struct
func (f FileRecord) Users() []FileUserRecord {
	return new(FileUserRecordRepository).Select(f.ID, "file_id")
}

// All loads all FileRecord structs from storage
func (f FileRecordRepository) All() (files []FileRecord) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Retrieve all files
	if files, err = db.GetAllFileRecords(); err != nil {
		log.Println(err.Error())
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return files
}
