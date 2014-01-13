package goat

import (
	"encoding/json"
	"log"
	"time"
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

// peerInfo represents a peer which will be marked as active or not
type peerInfo struct {
	UserID int `db:"user_id"`
	IP     string
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

	if err := db.SaveFileRecord(f); err != nil {
		log.Println(err.Error())
		return false
	}

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
	if f, err = db.LoadFileRecord(id, col); err != nil {
		log.Println(err.Error())
		return fileRecord{}
	}

	return f
}

// Completed returns the number of completions, active or not, on this file
func (f fileRecord) Completed() (completed int) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	if completed, err = db.CountFileRecordCompleted(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}
	return
}

// Seeders returns the number of seeders on this file
func (f fileRecord) Seeders() (seeders int) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	if seeders, err = db.CountFileRecordSeeders(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}
	return
}

// Leechers returns the number of leechers on this file
func (f fileRecord) Leechers() (leechers int) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return 0
	}
	if leechers, err = db.CountFileRecordLeechers(f.ID); err != nil {
		log.Println(err.Error())
		return -1
	}
	return
}

// PeerList returns the compact peer buffer for tracker announce, excluding self
func (f fileRecord) PeerList(exclude string, numwant int) (peers []byte) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if peers, err = db.GetFileRecordPeerList(f.InfoHash, exclude, numwant); err != nil {
		log.Println(err.Error())
		return nil
	}
	return
}

// PeerReaper reaps peers who have not recently announced on this torrent, and mark them inactive
func (f fileRecord) PeerReaper() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	users, err := db.GetInactiveUserInfo(f.ID, time.Duration(int64(static.Config.Interval))*time.Second+60)
	if err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.MarkFileUsersInactive(f.ID, users); err != nil {
		log.Println(err.Error())
		return false
	}

	if count := len(users); count > 0 {
		log.Printf("reaper: reaped %d peer(s) on file %d\n", count, f.ID)
	}
	return true
}

// Users loads all fileUserRecord structs associated with this fileRecord struct
func (f fileRecord) Users() []fileUserRecord {
	return new(fileUserRecordRepository).Select(f.ID, "file_id")
}

// All loads all fileRecord structs from storage
func (f fileRecordRepository) All() (files []fileRecord) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if files, err = db.GetAllFileRecords(); err != nil {
		log.Println(err.Error())
	}
	return files
}
