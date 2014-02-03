package data

import (
	"time"
)

var (
	// DBConnectFunc connects to a database backend
	DBConnectFunc func() (dbModel, error)
	// DBCloseFunc closes connection to a database backend
	DBCloseFunc = func() {}
	// DBNameFunc returns the name of a database backend
	DBNameFunc = func() string { return "" }
	// DBPingFunc checks connectivity to a database backend
	DBPingFunc = func() bool { return true }
)

// MySQLDSN is set via command-line, and can be used to override all MySQL configuration
var MySQLDSN *string

// QLDBPath is set via command-line, and can be used to override ql database location
var QLDBPath *string

// DBConnect connects to a database
func DBConnect() (dbModel, error) {
	return DBConnectFunc()
}

// DBName will retrieve the name of the database backend currently in use
func DBName() string {
	return DBNameFunc()
}

// DBPing will attempt to "ping" the database to verify connectivity
func DBPing() bool {
	return DBPingFunc()
}

// dbModel represents a database interface, and defines functions which act on it
type dbModel interface {
	Close() error

	// --- AnnounceLog.go ---
	DeleteAnnounceLog(interface{}, string) error
	LoadAnnounceLog(interface{}, string) (AnnounceLog, error)
	SaveAnnounceLog(AnnounceLog) error

	// --- APIKey.go ---
	DeleteAPIKey(interface{}, string) error
	LoadAPIKey(interface{}, string) (APIKey, error)
	SaveAPIKey(APIKey) error

	// --- FileRecord.go ---
	DeleteFileRecord(interface{}, string) error
	LoadFileRecord(interface{}, string) (FileRecord, error)
	SaveFileRecord(FileRecord) error
	CountFileRecordCompleted(int) (int, error)
	CountFileRecordSeeders(int) (int, error)
	CountFileRecordLeechers(int) (int, error)
	GetFileRecordPeerList(string, int) ([]Peer, error)
	GetInactiveUserInfo(int, time.Duration) ([]peerInfo, error)
	MarkFileUsersInactive(int, []peerInfo) error
	GetAllFileRecords() ([]FileRecord, error)

	// --- FileUserRecord.go ---
	DeleteFileUserRecord(int, int, string) error
	LoadFileUserRecord(int, int, string) (FileUserRecord, error)
	SaveFileUserRecord(FileUserRecord) error
	LoadFileUserRepository(interface{}, string) ([]FileUserRecord, error)

	// --- ScrapeLog.go ---
	DeleteScrapeLog(interface{}, string) error
	LoadScrapeLog(interface{}, string) (ScrapeLog, error)
	SaveScrapeLog(ScrapeLog) error

	// --- UserRecord.go ---
	DeleteUserRecord(interface{}, string) error
	LoadUserRecord(interface{}, string) (UserRecord, error)
	SaveUserRecord(UserRecord) error
	GetUserUploaded(int) (int64, error)
	GetUserDownloaded(int) (int64, error)
	GetUserSeeding(int) (int, error)
	GetUserLeeching(int) (int, error)

	// --- WhitelistRecord.go ---
	LoadWhitelistRecord(interface{}, string) (WhitelistRecord, error)
	SaveWhitelistRecord(WhitelistRecord) error
}
