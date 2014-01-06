package goat

import (
	"fmt"
	"time"
)

// FileRecord represents a file tracked by tracker
type FileRecord struct {
	ID         int
	InfoHash   string `db:"info_hash"`
	Verified   bool
	CreateTime int64 `db:"create_time"`
	UpdateTime int64 `db:"update_time"`
}

// Save FileRecord to storage
func (f FileRecord) Save() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	if err := db.SaveFileRecord(f); nil != err {
		Static.LogChan <- err.Error()
		return false
	}

	return true
}

// Load FileRecord from storage
func (f FileRecord) Load(id interface{}, col string) FileRecord {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return f
	}

	if f, err = db.LoadFileRecord(id, col); nil != err {
		Static.LogChan <- err.Error()
		return FileRecord{}
	}

	return f
}

// Completed returns the number of completions, active or not, on this file
func (f FileRecord) Completed() (completed int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return 0
	}

	if completed, err = db.CountFileRecordCompleted(f.ID); err != nil {
		Static.LogChan <- err.Error()
	}
	return
}

// Seeders returns the number of seeders on this file
func (f FileRecord) Seeders() (seeders int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return 0
	}

	if seeders, err = db.CountFileRecordSeeders(f.ID); nil != err {
		Static.LogChan <- err.Error()
	}
	return
}

// Leechers returns the number of leechers on this file
func (f FileRecord) Leechers() (leechers int) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return 0
	}

	if leechers, err = db.CountFileRecordLeechers(f.ID); nil != err {
		Static.LogChan <- err.Error()
	}
	return
}

// PeerList returns the compact peer buffer for tracker announce, excluding self
func (f FileRecord) PeerList(exclude string, numwant int) (peers []byte) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return
	}

	if peers, err = db.GetFileRecordPeerList(f.InfoHash, exclude, numwant); err != nil {
		Static.LogChan <- err.Error()
	}
	return
}

// PeerReaper reaps peers who have not recently announced on this torrent, and mark them inactive
func (f FileRecord) PeerReaper() bool {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	users, err := db.GetInactiveUserInfo(f.ID, time.Duration(int64(Static.Config.Interval))*time.Second+60)
	if err != nil {
		Static.LogChan <- err.Error()
		return false
	}

	if err := db.MarkFileUsersInactive(f.ID, users); nil != err {
		Static.LogChan <- err.Error()
	}

	if count := len(users); count > 0 {
		Static.LogChan <- fmt.Sprintf("reaper: reaped %d peer(s) on file %d", count, f.ID)
	}

	return true
}
