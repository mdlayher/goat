package goat

import (
	"log"
)

// fileUserRecord represents a file tracked by tracker
type fileUserRecord struct {
	FileID     int    `db:"file_id" json:"fileId"`
	UserID     int    `db:"user_id" json:"userId"`
	IP         string `json:"ip"`
	Active     bool   `json:"active"`
	Completed  bool   `json:"completed"`
	Announced  int    `json:"announced"`
	Uploaded   int64  `json:"uploaded"`
	Downloaded int64  `json:"downloaded"`
	Left       int64  `json:"left"`
	Time       int64  `json:"time"`
}

// fileUserRecordRepository is used to contain methods to load multiple fileRecord structs
type fileUserRecordRepository struct {
}

// Save fileUserRecord to storage
func (f fileUserRecord) Save() bool {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save fileUserRecord
	if err := db.SaveFileUserRecord(f); err != nil {
		log.Println(err.Error())
		return false
	}

	return true
}

// Load fileUserRecord from storage
func (f fileUserRecord) Load(fileID int, userID int, ip string) fileUserRecord {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return f
	}

	// Load fileUserRecord using file ID, user ID, IP triple
	if f, err = db.LoadFileUserRecord(fileID, userID, ip); err != nil {
		log.Println(err.Error())
		return fileUserRecord{}
	}

	return f
}

// Select loads selected fileUserRecord structs from storage
func (f fileUserRecordRepository) Select(id interface{}, col string) (files []fileUserRecord) {
	// Open database connection
	db, err := dbConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Load fileUserRecords matching specified conditions
	if files, err = db.LoadFileUserRepository(id, col); err != nil {
		log.Println(err.Error())
	}

	return
}
