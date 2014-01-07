package goat

import (
	"encoding/binary"
	"net"
	"time"
)

var (
	dbConnectFunc func() (DbModel, error)
	dbCloseFunc   func() = func() {}
)

func DBConnect() (DbModel, error) {
	return dbConnectFunc()
}

type DbModel interface {
	// --- announceLog.go ---
	LoadAnnounceLog(interface{}, string) (AnnounceLog, error)
	SaveAnnounceLog(AnnounceLog) error

	// --- fileRecord.go ---
	LoadFileRecord(interface{}, string) (FileRecord, error)
	SaveFileRecord(FileRecord) error
	CountFileRecordCompleted(int) (int, error)
	CountFileRecordSeeders(int) (int, error)
	CountFileRecordLeechers(int) (int, error)
	GetFileRecordPeerList(string, string, int) ([]byte, error)
	GetInactiveUserInfo(int, time.Duration) ([]userinfo, error)
	MarkFileUsersInactive(int, []userinfo) error

	// --- fileUserRecord.go ---
	LoadFileUserRecord(int, int, string) (FileUserRecord, error)
	SaveFileUserRecord(FileUserRecord) error

	// --- scrapeLog.go ---
	LoadScrapeLog(interface{}, string) (ScrapeLog, error)
	SaveScrapeLog(ScrapeLog) error

	// --- userRecord.go ---
	LoadUserRecord(interface{}, string) (UserRecord, error)
	SaveUserRecord(UserRecord) error
	GetUserUploaded(int) (int64, error)
	GetUserDownloaded(int) (int64, error)

	// --- whitelistRecord.go ---
	LoadWhitelistRecord(interface{}, string) (WhitelistRecord, error)
	SaveWhitelistRecord(WhitelistRecord) error
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
