package goat

import (
	"time"
)

var (
	dbConnectFunc func() (dbModel, error)
	dbCloseFunc   func()        = func() {}
	dbNameFunc    func() string = func() string { return "" }
	dbPingFunc    func() bool   = func() bool { return true }
)

// dbConnect connects to a database
func dbConnect() (dbModel, error) {
	return dbConnectFunc()
}

// dbName will retrieve the name of the database backend currently in use
func dbName() string {
	return dbNameFunc()
}

// dbPing will attempt to "ping" the database to verify connectivity
func dbPing() bool {
	return dbPingFunc()
}

// dbModel represents a database interface, and defines functions which act on it
type dbModel interface {
	// --- announceLog.go ---
	DeleteAnnounceLog(interface{}, string) error
	LoadAnnounceLog(interface{}, string) (announceLog, error)
	SaveAnnounceLog(announceLog) error

	// --- apiKey.go ---
	DeleteApiKey(interface{}, string) error
	LoadApiKey(interface{}, string) (apiKey, error)
	SaveApiKey(apiKey) error

	// --- fileRecord.go ---
	DeleteFileRecord(interface{}, string) error
	LoadFileRecord(interface{}, string) (fileRecord, error)
	SaveFileRecord(fileRecord) error
	CountFileRecordCompleted(int) (int, error)
	CountFileRecordSeeders(int) (int, error)
	CountFileRecordLeechers(int) (int, error)
	GetFileRecordPeerList(string, string, int) ([]byte, error)
	GetInactiveUserInfo(int, time.Duration) ([]peerInfo, error)
	MarkFileUsersInactive(int, []peerInfo) error
	GetAllFileRecords() ([]fileRecord, error)

	// --- fileUserRecord.go ---
	DeleteFileUserRecord(int, int, string) error
	LoadFileUserRecord(int, int, string) (fileUserRecord, error)
	SaveFileUserRecord(fileUserRecord) error
	LoadFileUserRepository(interface{}, string) ([]fileUserRecord, error)

	// --- scrapeLog.go ---
	DeleteScrapeLog(interface{}, string) error
	LoadScrapeLog(interface{}, string) (scrapeLog, error)
	SaveScrapeLog(scrapeLog) error

	// --- userRecord.go ---
	DeleteUserRecord(interface{}, string) error
	LoadUserRecord(interface{}, string) (userRecord, error)
	SaveUserRecord(userRecord) error
	GetUserUploaded(int) (int64, error)
	GetUserDownloaded(int) (int64, error)
	GetUserSeeding(int) (int, error)
	GetUserLeeching(int) (int, error)

	// --- whitelistRecord.go ---
	LoadWhitelistRecord(interface{}, string) (whitelistRecord, error)
	SaveWhitelistRecord(whitelistRecord) error
}
