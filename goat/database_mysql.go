// +build !ql

package goat

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Bring in the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var MysqlDSN *string

// init performs startup routines for database_mysql
func init() {
	// dbConnectFunc connects to MySQL database
	dbConnectFunc = func() (dbModel, error) {
		var conn string
		// Generate connection string using configuration
		if MysqlDSN == nil || *MysqlDSN == "" {
			conn = fmt.Sprintf("%s:%s@/%s", static.Config.DB.Username, static.Config.DB.Password, static.Config.DB.Database)
		} else {
			conn = *MysqlDSN
		}

		log.Printf("dialing %s", conn)
		// Return connection and associated errors
		db, err := sqlx.Connect("mysql", conn)
		return &dbw{db}, err
	}

	// dbNameFunc returns the name of this backend
	dbNameFunc = func() string {
		return "MySQL"
	}

	// dbPingFunc tests the connection to the MySQL database
	dbPingFunc = func() bool {
		db, err := dbConnect()
		if err != nil {
			log.Println(err.Error())
			return false
		}

		if err = db.(*dbw).Close(); err != nil {
			log.Println(err.Error())
			return false
		}

		return true
	}
}

// dbw contains a sqlx MySQL database connection
type dbw struct {
	*sqlx.DB
}

// Close database connection
func (db *dbw) Close() error {
	return db.DB.Close()
}

// --- announceLog.go ---

// DeleteAnnounceLog deletes an announceLog using a defined ID and column
func (db *dbw) DeleteAnnounceLog(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM announce_log WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadAnnounceLog loads an announceLog using a defined ID and column for query
func (db *dbw) LoadAnnounceLog(id interface{}, col string) (announceLog, error) {
	data := announceLog{}

	if err := db.Get(&data, "SELECT * FROM announce_log WHERE `"+col+"`=?", id); err != nil && err != sql.ErrNoRows {
		return announceLog{}, err
	}

	return data, nil
}

// SaveAnnounceLog saves an announceLog to database
func (db *dbw) SaveAnnounceLog(a announceLog) error {
	query := "INSERT INTO announce_log " +
		"(`info_hash`, `passkey`, `key`, `ip`, `port`, `udp`, `uploaded`, `downloaded`, `left`, `event`, `client`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.Passkey, a.Key, a.IP, a.Port, a.UDP, a.Uploaded, a.Downloaded, a.Left, a.Event, a.Client)

	return tx.Commit()
}

// --- apiKey.go ---

// DeleteAPIKey deletes an apiKey using a defined ID and column
func (db *dbw) DeleteAPIKey(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM api_keys WHERE `"+col+"` = ?", id)

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return tx.Commit()
}

// LoadAPIKey loads an apiKey using a defined ID and column for query
func (db *dbw) LoadAPIKey(id interface{}, col string) (apiKey, error) {
	key := apiKey{}

	err := db.Get(&key, "SELECT * FROM api_keys WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		return apiKey{}, err
	}

	return key, nil
}

// SaveAPIKey saves an apiKey to the database
func (db *dbw) SaveAPIKey(key apiKey) error {
	query := "INSERT INTO api_keys (`user_id`, `key`, `salt`) " +
		"VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE " +
		"`key`=values(`key`), `salt`=values(`salt`);"

	tx := db.MustBegin()
	tx.Execl(query, key.UserID, key.Key, key.Salt)

	return tx.Commit()
}

// --- fileRecord.go ---

// DeleteFileRecord deletes an announceLog using a defined ID and column
func (db *dbw) DeleteFileRecord(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM files WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadFileRecord loads a fileRecord using a defined ID and column for query
func (db *dbw) LoadFileRecord(id interface{}, col string) (fileRecord, error) {
	data := fileRecord{}

	if err := db.Get(&data, "SELECT * FROM files WHERE `"+col+"`=?", id); err != nil && err != sql.ErrNoRows {
		return fileRecord{}, err
	}

	return data, nil
}

// SaveFileRecord saves a fileRecord to the database
func (db *dbw) SaveFileRecord(f fileRecord) error {
	query := "INSERT INTO files " +
		"(`info_hash`, `verified`, `create_time`, `update_time`) " +
		"VALUES (?, ?, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()) " +
		"ON DUPLICATE KEY UPDATE " +
		"`verified`=values(`verified`), `update_time`=UNIX_TIMESTAMP();"

	tx := db.MustBegin()
	tx.Execl(query, f.InfoHash, f.Verified)

	return tx.Commit()
}

// CountFileRecordCompleted counts the number of peers who have completed this file
func (db *dbw) CountFileRecordCompleted(id int) (int, error) {
	// Calculate number of completions on this file, defined as users who are completed, and 0 left
	query := "SELECT COUNT(user_id) AS completed FROM files_users WHERE file_id = ? AND completed = 1 AND `left` = 0;"
	result := struct{ Completed int }{0}

	if err := db.Get(&result, query, id); err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return result.Completed, nil
}

// CountFileRecordCompleted counts the number of peers who are actively seeding this file
func (db *dbw) CountFileRecordSeeders(id int) (int, error) {
	// Calculate number of seeders on this file, defined as users who are active, completed, and 0 left
	query := "SELECT COUNT(user_id) AS seeders FROM files_users WHERE file_id = ? AND active = 1 AND completed = 1 AND `left` = 0;"
	result := struct{ Seeders int }{0}

	if err := db.Get(&result, query, id); err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return result.Seeders, nil
}

// CountFileRecordCompleted counts the number of peers who are actively leeching this file
func (db *dbw) CountFileRecordLeechers(id int) (int, error) {
	// Calculate number of leechers on this file, defined as users who are active, completed, and 0 left
	query := "SELECT COUNT(user_id) AS leechers FROM files_users WHERE file_id = ? AND active = 1 AND completed = 0 AND `left` > 0;"
	result := struct{ Leechers int }{0}

	if err := db.Get(&result, query, id); err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return result.Leechers, nil
}

// GetFileRecordPeerList returns a compact peer list containing IP/port pairs
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
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Buffer for compact list
	buf := make([]byte, 0)

	for rows.Next() {
		if err = rows.StructScan(&result); err != nil {
			log.Println(err.Error())
			break
		}

		buf = append(buf, ip2b(result.IP, result.Port)...)
	}

	return buf, nil
}

// GetInactiveUserInfo returns a list of users who have not been active for the specified time interval
func (db *dbw) GetInactiveUserInfo(fid int, interval time.Duration) (users []peerInfo, err error) {
	query := `SELECT user_id, ip FROM files_users
		WHERE time < (UNIX_TIMESTAMP() - ?)
		AND active = 1
		AND file_id = ?;`

	result := peerInfo{}
	checkInterval := int(interval / time.Second)

	var rows *sqlx.Rows
	if rows, err = db.Queryx(query, checkInterval, fid); err == nil && err != sql.ErrNoRows {
		for rows.Next() {
			if err = rows.StructScan(&result); err == nil {
				users = append(users, result)
			}
		}
	}

	return
}

// MarkFileUsersInactive sets users to be inactive once they have been reaped
func (db *dbw) MarkFileUsersInactive(fid int, users []peerInfo) error {
	query := "UPDATE files_users SET active = 0 WHERE file_id = ? AND user_id = ? AND ip = ?;"

	tx := db.MustBegin()
	for _, u := range users {
		tx.Execl(query, fid, u.UserID, u.IP)
	}

	return tx.Commit()
}

// GetAllFileRecords returns a list of all fileRecords known to the database
func (db *dbw) GetAllFileRecords() ([]fileRecord, error) {
	rows, err := db.Queryx("SELECT * FROM files")
	files, file := []fileRecord{}, fileRecord{}

	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return files, err
	}

	for rows.Next() {
		if err = rows.StructScan(&file); err != nil {
			break
		}

		files = append(files[:], file)
	}

	return files, nil
}

// --- fileUserRecord.go ---

// DeleteFileUserRecord deletes a fileUserRecord using using a file ID, user ID, and IP triple
func (db *dbw) DeleteFileUserRecord(fid, uid int, ip string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?", fid, uid, ip)

	return tx.Commit()
}

// LoadFileUserRecord loads a fileUserRecord using a file ID, user ID, and IP triple
func (db *dbw) LoadFileUserRecord(fid, uid int, ip string) (fileUserRecord, error) {
	query := "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?;"

	data := fileUserRecord{}
	if err := db.Get(&data, query, fid, uid, ip); err != nil && err != sql.ErrNoRows {
		return fileUserRecord{}, err
	}

	return data, nil
}

// SaveFileUserRecord saves a fileUserRecord to the database
func (db *dbw) SaveFileUserRecord(f fileUserRecord) error {
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

// LoadFileUserRepository loads all fileUserRecords matching a defined ID and column for query
func (db *dbw) LoadFileUserRepository(id interface{}, col string) ([]fileUserRecord, error) {
	rows, err := db.Queryx("SELECT * FROM files_users WHERE `"+col+"`=?", id)
	files, user := []fileUserRecord{}, fileUserRecord{}

	if err != nil && err != sql.ErrNoRows {
		return files, err
	}

	for rows.Next() {
		if err = rows.StructScan(&user); err != nil {
			log.Println(err.Error())
			break
		}

		files = append(files[:], user)
	}

	return files, nil
}

// --- scrapeLog.go ---

// DeleteScrapeLog deletes a scrapeLog using a defined ID and column
func (db *dbw) DeleteScrapeLog(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM scrape_log WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadScrapeLog loads a scrapeLog using a defined ID and column for query
func (db *dbw) LoadScrapeLog(id interface{}, col string) (scrapeLog, error) {
	query := "SELECT * FROM scrape_log WHERE `" + col + "`=?;"

	data := scrapeLog{}
	if err := db.Get(&data, query, id); err != nil && err != sql.ErrNoRows {
		return scrapeLog{}, err
	}

	return data, nil
}

// SaveScrapeLog saves a scrapeLog to the database
func (db *dbw) SaveScrapeLog(s scrapeLog) error {
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"

	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.IP)

	return tx.Commit()
}

// --- userRecord.go ---

// DeleteUserRecord deletes a userRecord using a defined ID and column
func (db *dbw) DeleteUserRecord(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM users WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadUserRecord loads a userRecord using a defined ID and column for query
func (db *dbw) LoadUserRecord(id interface{}, col string) (userRecord, error) {
	query := "SELECT * FROM users WHERE `" + col + "`=?;"

	data := userRecord{}
	if err := db.Get(&data, query, id); err != nil && err != sql.ErrNoRows {
		return userRecord{}, err
	}

	return data, nil
}

// SaveUserRecord saves a userRecord to the database
func (db *dbw) SaveUserRecord(u userRecord) error {
	query := "INSERT INTO users " +
		"(`username`, `passkey`, `torrent_limit`) " +
		"VALUES (?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE " +
		"`username`=values(`username`), `passkey`=values(`passkey`), `torrent_limit`=values(`torrent_limit`);"

	tx := db.MustBegin()
	tx.Execl(query, u.Username, u.Passkey, u.TorrentLimit)

	return tx.Commit()
}

// GetUserUploaded calculates the total number of bytes this user has uploaded
func (db *dbw) GetUserUploaded(uid int) (int64, error) {
	// Calculate sum of this user's upload via their file/user relationship records
	query := "SELECT SUM(uploaded) AS uploaded FROM files_users WHERE user_id=?;"

	result := struct{ Uploaded int64 }{0}
	if err := db.Get(&result, query, uid); err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return result.Uploaded, nil
}

// GetUserDownloaded calculates the total number of bytes this user has downloaded
func (db *dbw) GetUserDownloaded(uid int) (int64, error) {
	// Calculate sum of this user's download via their file/user relationship records
	query := "SELECT SUM(downloaded) AS downloaded FROM files_users WHERE user_id=?;"

	result := struct{ Downloaded int64 }{0}
	if err := db.Get(&result, query, uid); err != nil && err != sql.ErrNoRows {
		return -1, err
	}

	return result.Downloaded, nil
}

// GetUserSeeding calculates the total number of files this user is actively seeding
func (db *dbw) GetUserSeeding(uid int) (int, error) {
	// Calculate sum of this user's seeding torrents via their file/user relationship records
	query := "SELECT COUNT(user_id) AS seeding FROM files_users WHERE user_id = ? AND active = 1 AND completed = 1 AND `left` = 0;"

	result := struct{ Seeding int }{0}
	if err := db.Get(&result, query, uid); err != nil {
		return -1, err
	}

	return result.Seeding, nil
}

// GetUserLeeching calculates the total number of files this user is actively leeching
func (db *dbw) GetUserLeeching(uid int) (int, error) {
	// Calculate sum of this user's leeching torrents via their file/user relationship records
	query := "SELECT COUNT(user_id) AS leeching FROM files_users WHERE user_id = ? AND active = 1 AND completed = 0 AND `left` > 0;"

	result := struct{ Leeching int }{0}
	if err := db.Get(&result, query, uid); err != nil {
		return -1, err
	}

	return result.Leeching, nil
}

// --- whitelistRecord.go ---

// LoadWhitelistRecord loads a whitelistRecord using a defined ID and column for query
func (db *dbw) LoadWhitelistRecord(id interface{}, col string) (whitelistRecord, error) {
	query := "SELECT * FROM whitelist WHERE `" + col + "`=?;"

	result := whitelistRecord{}
	if err := db.Get(&result, query, id); err != nil && err != sql.ErrNoRows {
		return whitelistRecord{}, err
	}

	return result, nil
}

// SaveWhitelistRecord saves a whitelistRecord to the database
func (db *dbw) SaveWhitelistRecord(w whitelistRecord) error {
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
