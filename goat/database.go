package goat

import (
	"encoding/binary"
	"net"
	"time"
)

var (
	dbConnectFunc func() (dbmodel, error)
	dbCloseFunc   func()      = func() {}
	dbPingFunc    func() bool = func() bool { return true }
)

// dbConnect connects to a database
func dbConnect() (dbmodel, error) {
	return dbConnectFunc()
}

// Attempt to "ping" the database to verify connectivity
func dbPing() bool {
	return dbPingFunc()
}

type dbmodel interface {
	// --- announceLog.go ---
	LoadAnnounceLog(interface{}, string) (announceLog, error)
	SaveAnnounceLog(announceLog) error

	// --- apiKey.go ---
	LoadApiKey(interface{}, string) (apiKey, error)
	SaveApiKey(apiKey) error

	// --- fileRecord.go ---
	LoadFileRecord(interface{}, string) (fileRecord, error)
	SaveFileRecord(fileRecord) error
	CountFileRecordCompleted(int) (int, error)
	CountFileRecordSeeders(int) (int, error)
	CountFileRecordLeechers(int) (int, error)
	GetFileRecordPeerList(string, string, int) ([]byte, error)
	GetInactiveUserInfo(int, time.Duration) ([]userinfo, error)
	MarkFileUsersInactive(int, []userinfo) error
	GetAllFileRecords() ([]fileRecord, error)

	// --- fileUserRecord.go ---
	LoadFileUserRecord(int, int, string) (fileUserRecord, error)
	SaveFileUserRecord(fileUserRecord) error
	LoadFileUserRepository(interface{}, string) ([]fileUserRecord, error)

	// --- scrapeLog.go ---
	LoadScrapeLog(interface{}, string) (scrapeLog, error)
	SaveScrapeLog(scrapeLog) error

	// --- userRecord.go ---
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

type userinfo struct {
	UserID int `db:"user_id"`
	IP     string
}

// Parse IP and port into byte buffer
func ip2b(ip_ string, port_ uint16) []byte {
	ip, port := [4]byte{}, [2]byte{}
	binary.BigEndian.PutUint32(ip[:],
		binary.BigEndian.Uint32(net.ParseIP(ip_).To4()))
	binary.BigEndian.PutUint16(port[:], port_)
	return append(ip[:], port[:]...)
}
