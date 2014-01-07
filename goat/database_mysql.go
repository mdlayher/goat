// +build !ql

package goat

import (
	"fmt"
	"time"

	// Bring in the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func init() {
	// dbConnectFunc connects to MySQL database
	dbConnectFunc = func() (DbModel, error) {
		// Generate connection string using configuration
		conn := fmt.Sprintf("%s:%s@/%s", Static.Config.DB.Database, Static.Config.DB.Username, Static.Config.DB.Password)
		// Return connection and associated errors
		db, err := sqlx.Connect("mysql", conn)
		return &dbw{db}, err
	}
}

type dbw struct {
	*sqlx.DB
}

// --- announceLog.go ---

func (db *dbw) LoadAnnounceLog(id interface{}, col string) (AnnounceLog, error) {
	data := AnnounceLog{}
	if err := db.Get(&data, "SELECT * FROM announce_log WHERE `"+col+"`=?", id); err != nil {
		return AnnounceLog{}, err
	}
	return data, nil
}

func (db *dbw) SaveAnnounceLog(a AnnounceLog) error {
	query := "INSERT INTO announce_log " +
		"(`info_hash`, `passkey`, `key`, `ip`, `port`, `udp`, `uploaded`, `downloaded`, `left`, `event`, `client`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.Passkey, a.Key, a.IP, a.Port, a.UDP, a.Uploaded, a.Downloaded, a.Left, a.Event, a.Client)

	return tx.Commit()
}

// --- fileRecord.go ---

func (db *dbw) LoadFileRecord(id interface{}, col string) (FileRecord, error) {
	data := FileRecord{}
	if err := db.Get(&data, "SELECT * FROM files WHERE `"+col+"`=?", id); err != nil {
		return FileRecord{}, err
	}
	return data, nil
}

func (db *dbw) SaveFileRecord(f FileRecord) error {
	query := "INSERT INTO files " +
		"(`info_hash`, `verified`, `create_time`, `update_time`) " +
		"VALUES (?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()) " +
		"ON DUPLICATE KEY UPDATE " +
		"`verified`=values(`verified`), `update_time`=UNIX_TIMESTAMP();"
	tx := db.MustBegin()
	tx.Execl(query, f.InfoHash, f.Verified)
	return tx.Commit()
}

func (db *dbw) CountFileRecordCompleted(id int) (int, error) {
	// Calculate number of completions on this file, defined as users who are completed, and 0 left
	query := "SELECT COUNT(user_id) AS completed FROM files_users WHERE file_id = ? AND completed = 1 AND `left` = 0;"
	result := struct{ Completed int }{0}

	if err := db.Get(&result, query, id); err != nil {
		return -1, err
	}

	return result.Completed, nil
}

func (db *dbw) CountFileRecordSeeders(id int) (int, error) {
	// Calculate number of seeders on this file, defined as users who are active, completed, and 0 left
	query := "SELECT COUNT(user_id) AS seeders FROM files_users WHERE file_id = ? AND active = 1 AND completed = 1 AND `left` = 0;"
	result := struct{ Seeders int }{0}

	if err := db.Get(&result, query, id); err != nil {
		return -1, err
	}

	return result.Seeders, nil
}

func (db *dbw) CountFileRecordLeechers(id int) (int, error) {
	// Calculate number of leechers on this file, defined as users who are active, completed, and 0 left
	query := "SELECT COUNT(user_id) AS leechers FROM files_users WHERE file_id = ? AND active = 1 AND completed = 0 AND `left` > 0;"
	result := struct{ Leechers int }{0}

	if err := db.Get(&result, query, id); err != nil {
		return -1, err
	}

	return result.Leechers, nil
}

func (db *dbw) GetFileRecordPeerList(infohash, exclude string, limit int) ([]byte, error) {
	// Get IP and port of all peers who are active and seeding this file
	query := `SELECT DISTINCT announce_log.ip,announce_log.port FROM announce_log
		JOIN files ON announce_log.info_hash = files.info_hash
		JOIN files_users ON files.id = files_users.file_id
		AND announce_log.ip = files_users.ip
		WHERE files_users.active=1
		AND files.info_hash=?
		AND announce_log.ip != ?
		LIMIT ?;`
	result := struct {
		IP   string
		Port uint16
	}{"", 0}

	rows, err := db.Queryx(query, infohash, exclude, limit)
	if err != nil {
		return nil, err
	}

	// Buffer for compact list
	buf := make([]byte, 0)

	for rows.Next() {
		rows.StructScan(&result)
		Static.LogChan <- fmt.Sprintf("peer: %s:%d", result.IP, result.Port)
		buf = append(buf, ip2b(result.IP, result.Port)...)
	}

	return buf, nil
}

func (db *dbw) GetInactiveUserInfo(fid int, interval_ time.Duration) (users []userinfo, err error) {
	query := `SELECT user_id, ip FROM files_users
		WHERE time < (UNIX_TIMESTAMP() - ?)
		AND active = 1
		AND file_id = ?;`
	result := userinfo{}
	interval := int(interval_ / time.Second)
	var rows *sqlx.Rows

	if rows, err = db.Queryx(query, interval, fid); err == nil {
		for rows.Next() {
			if err = rows.StructScan(&result); nil == err {
				users = append(users, result)
			}
		}
	}
	return
}

func (db *dbw) MarkFileUsersInactive(fid int, users []userinfo) error {
	query := "UPDATE files_users SET active = 0 WHERE file_id = ? AND user_id = ? AND ip = ?;"
	tx := db.MustBegin()
	for _, u := range users {
		tx.Execl(query, fid, u.UserID, u.IP)
	}
	return tx.Commit()
}

// --- fileUserRecord.go ---

func (db *dbw) LoadFileUserRecord(fid, uid int, ip string) (FileUserRecord, error) {
	query := "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?"
	data := FileUserRecord{}
	if err := db.Get(&data, query, fid, uid, ip); err != nil {
		return FileUserRecord{}, err
	}
	return data, nil
}

func (db *dbw) SaveFileUserRecord(f FileUserRecord) error {
	// Insert or update a file/user relationship record
	query := "INSERT INTO files_users " +
		"(`file_id`, `user_id`, `ip`, `active`, `completed`, `announced`, `uploaded`, `downloaded`, `left`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP()) " +
		"ON DUPLICATE KEY UPDATE " +
		"`active`=values(`active`), `completed`=values(`completed`), `announced`=values(`announced`), " +
		"`uploaded`=values(`uploaded`), `downloaded`=values(`downloaded`), `left`=values(`left`), " +
		"`time`=UNIX_TIMESTAMP();"
	tx := db.MustBegin()
	tx.Execl(query, f.FileID, f.UserID, f.IP, f.Active, f.Completed, f.Announced, f.Uploaded, f.Downloaded, f.Left)
	return tx.Commit()
}

// --- scrapeLog.go ---

func (db *dbw) LoadScrapeLog(id interface{}, col string) (ScrapeLog, error) {
	query := "SELECT * FROM announce_log WHERE `" + col + "`=?"
	data := ScrapeLog{}
	if err := db.Get(&data, query, id); err != nil {
		return ScrapeLog{}, err
	}
	return data, nil
}

func (db *dbw) SaveScrapeLog(s ScrapeLog) error {
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"
	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.IP)
	return tx.Commit()
}

// --- userRecord.go ---

func (db *dbw) LoadUserRecord(id interface{}, col string) (UserRecord, error) {
	query := "SELECT * FROM users WHERE `" + col + "`=?"
	data := UserRecord{}
	if err := db.Get(&data, query, id); err != nil {
		return UserRecord{}, err
	}
	return data, nil
}

func (db *dbw) SaveUserRecord(u UserRecord) error {
	query := "INSERT INTO users " +
		"(`username`, `passkey`, `torrent_limit`) " +
		"VALUES (?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"`username`=values(`username`), `passkey`=values(`passkey`), `torrent_limit`=values(`torrent_limit`);"
	tx := db.MustBegin()
	tx.Execl(query, u.Username, u.Passkey, u.TorrentLimit)
	return tx.Commit()
}

func (db *dbw) GetUserUploaded(uid int) (int64, error) {
	// Calculate sum of this user's upload via their file/user relationship records
	query := "SELECT SUM(uploaded) AS uploaded FROM files_users WHERE user_id=?"
	result := struct{ Uploaded int64 }{0}
	if err := db.Get(&result, query, uid); err != nil {
		return -1, err
	}
	return result.Uploaded, nil
}

func (db *dbw) GetUserDownloaded(uid int) (int64, error) {
	// Calculate sum of this user's download via their file/user relationship records
	query := "SELECT SUM(downloaded) AS downloaded FROM files_users WHERE user_id=?"
	result := struct{ Downloaded int64 }{0}
	if err := db.Get(&result, query, uid); err != nil {
		return -1, err
	}
	return result.Downloaded, nil
}

// --- whitelistRecord.go ---

func (db *dbw) LoadWhitelistRecord(id interface{}, col string) (WhitelistRecord, error) {
	query := "SELECT * FROM whitelist WHERE `" + col + "`=?"
	result := WhitelistRecord{}
	if err := db.Get(&result, query, id); err != nil {
		return WhitelistRecord{}, err
	}
	return result, nil
}

func (db *dbw) SaveWhitelistRecord(w WhitelistRecord) error {
	// NOTE: Not using INSERT IGNORE because it ignores all errors
	// Thanks: http://stackoverflow.com/questions/2366813/on-duplicate-key-ignore
	query := "INSERT INTO whitelist " +
		"(`client`, `approved`) " +
		"VALUES (?, ?) " +
		"ON DUPLICATE KEY UPDATE `client`=`client`;"
	tx := db.MustBegin()
	tx.Execl(query, w.Client, w.Approved)
	return tx.Commit()
}
