package data

import (
	"log"
)

// FileUserRecord represents a file tracked by tracker
type FileUserRecord struct {
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

// FileUserRecordRepository is used to contain methods to load multiple FileRecord structs
type FileUserRecordRepository struct {
}

// Delete FileUserRecord from storage
func (f FileUserRecord) Delete() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Delete FileUserRecord
	if err = db.DeleteFileUserRecord(f.FileID, f.UserID, f.IP); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Save FileUserRecord to storage
func (f FileUserRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return false
	}

	// Save FileUserRecord
	if err := db.SaveFileUserRecord(f); err != nil {
		log.Println(err.Error())
		return false
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return true
}

// Load FileUserRecord from storage
func (f FileUserRecord) Load(fileID int, userID int, ip string) FileUserRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return f
	}

	// Load FileUserRecord using file ID, user ID, IP triple
	if f, err = db.LoadFileUserRecord(fileID, userID, ip); err != nil {
		log.Println(err.Error())
		return FileUserRecord{}
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return f
}

// Select loads selected FileUserRecord structs from storage
func (f FileUserRecordRepository) Select(id interface{}, col string) (files []FileUserRecord) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Load FileUserRecords matching specified conditions
	if files, err = db.LoadFileUserRepository(id, col); err != nil {
		log.Println(err.Error())
	}

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return
}
