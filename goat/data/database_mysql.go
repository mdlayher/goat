// +build !ql

package data

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/mdlayher/goat/goat/common"

	// Bring in the MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// init performs startup routines for database_mysql
func init() {
	// DBConnectFunc connects to MySQL database
	DBConnectFunc = func() (dbModel, error) {
		var conn string
		// Generate connection string using configuration file
		if MySQLDSN == nil || *MySQLDSN == "" {
			conn = fmt.Sprintf("%s:%s@/%s", common.Static.Config.DB.Username, common.Static.Config.DB.Password, common.Static.Config.DB.Database)
		} else {
			// Use connection string passed by command line flag
			conn = *MySQLDSN
		}

		// Return connection and associated errors
		db, err := sqlx.Connect("mysql", conn)
		return &dbw{db}, err
	}

	// DBNameFunc returns the name of this backend
	DBNameFunc = func() string {
		return "MySQL"
	}

	// DBPingFunc tests the connection to the MySQL database
	DBPingFunc = func() bool {
		db, err := DBConnect()
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

// --- AnnounceLog.go ---

// DeleteAnnounceLog deletes an AnnounceLog using a defined ID and column
func (db *dbw) DeleteAnnounceLog(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM announce_log WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadAnnounceLog loads an AnnounceLog using a defined ID and column for query
func (db *dbw) LoadAnnounceLog(id interface{}, col string) (AnnounceLog, error) {
	data := AnnounceLog{}

	if err := db.Get(&data, "SELECT * FROM announce_log WHERE `"+col+"`=?", id); err != nil && err != sql.ErrNoRows {
		return AnnounceLog{}, err
	}

	return data, nil
}

// SaveAnnounceLog saves an AnnounceLog to database
func (db *dbw) SaveAnnounceLog(a AnnounceLog) error {
	query := "INSERT INTO announce_log " +
		"(`info_hash`, `passkey`, `key`, `ip`, `port`, `udp`, `uploaded`, `downloaded`, `left`, `event`, `client`, `time`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, UNIX_TIMESTAMP());"

	tx := db.MustBegin()
	tx.Execl(query, a.InfoHash, a.Passkey, a.Key, a.IP, a.Port, a.UDP, a.Uploaded, a.Downloaded, a.Left, a.Event, a.Client)

	return tx.Commit()
}

// --- APIKey.go ---

// DeleteAPIKey deletes an APIKey using a defined ID and column
func (db *dbw) DeleteAPIKey(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM api_keys WHERE `"+col+"` = ?", id)

	if err := db.Close(); err != nil {
		log.Println(err.Error())
	}

	return tx.Commit()
}

// LoadAPIKey loads an APIKey using a defined ID and column for query
func (db *dbw) LoadAPIKey(id interface{}, col string) (APIKey, error) {
	key := APIKey{}

	err := db.Get(&key, "SELECT * FROM api_keys WHERE `"+col+"`=?", id)
	if err != nil && err != sql.ErrNoRows {
		return APIKey{}, err
	}

	return key, nil
}

// SaveAPIKey saves an APIKey to the database
func (db *dbw) SaveAPIKey(key APIKey) error {
	query := "INSERT INTO api_keys (`user_id`, `key`, `salt`) " +
		"VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE " +
		"`key`=values(`key`), `salt`=values(`salt`);"

	tx := db.MustBegin()
	tx.Execl(query, key.UserID, key.Key, key.Salt)

	return tx.Commit()
}

// --- FileRecord.go ---

// DeleteFileRecord deletes an AnnounceLog using a defined ID and column
func (db *dbw) DeleteFileRecord(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM files WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadFileRecord loads a FileRecord using a defined ID and column for query
func (db *dbw) LoadFileRecord(id interface{}, col string) (FileRecord, error) {
	data := FileRecord{}

	if err := db.Get(&data, "SELECT * FROM files WHERE `"+col+"`=?", id); err != nil && err != sql.ErrNoRows {
		return FileRecord{}, err
	}

	return data, nil
}

// SaveFileRecord saves a FileRecord to the database
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

// GetFileRecordPeerList returns a list of Peers, containing IP/port pairs
func (db *dbw) GetFileRecordPeerList(infoHash string, limit int) ([]Peer, error) {
	// Get IP and port of all peers who are active and seeding this file
	query := `SELECT DISTINCT announce_log.ip,announce_log.port FROM announce_log
		JOIN files ON announce_log.info_hash = files.info_hash
		JOIN files_users ON files.id = files_users.file_id
		AND announce_log.ip = files_users.ip
		WHERE files_users.active=1
		AND files.info_hash=?
		LIMIT ?;`

	// Create output list
	peers := make([]Peer, 0)

	// Perform query
	rows, err := db.Queryx(query, infoHash, limit)
	if err != nil && err != sql.ErrNoRows {
		return peers, err
	}

	// Scan each peer from database
	peer := Peer{}
	for rows.Next() {
		// Check for error
		if err = rows.StructScan(&peer); err != nil {
			return peers, nil
		}

		log.Println(peer)

		// Append peer to list
		peers = append(peers[:], peer)
	}

	return peers, nil
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

// GetAllFileRecords returns a list of all FileRecords known to the database
func (db *dbw) GetAllFileRecords() ([]FileRecord, error) {
	rows, err := db.Queryx("SELECT * FROM files")
	files, file := []FileRecord{}, FileRecord{}

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

// --- FileUserRecord.go ---

// DeleteFileUserRecord deletes a FileUserRecord using using a file ID, user ID, and IP triple
func (db *dbw) DeleteFileUserRecord(fid, uid int, ip string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?", fid, uid, ip)

	return tx.Commit()
}

// LoadFileUserRecord loads a FileUserRecord using a file ID, user ID, and IP triple
func (db *dbw) LoadFileUserRecord(fid, uid int, ip string) (FileUserRecord, error) {
	query := "SELECT * FROM files_users WHERE `file_id`=? AND `user_id`=? AND `ip`=?;"

	data := FileUserRecord{}
	if err := db.Get(&data, query, fid, uid, ip); err != nil && err != sql.ErrNoRows {
		return FileUserRecord{}, err
	}

	return data, nil
}

// SaveFileUserRecord saves a FileUserRecord to the database
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

// LoadFileUserRepository loads all FileUserRecords matching a defined ID and column for query
func (db *dbw) LoadFileUserRepository(id interface{}, col string) ([]FileUserRecord, error) {
	rows, err := db.Queryx("SELECT * FROM files_users WHERE `"+col+"`=?", id)
	files, user := []FileUserRecord{}, FileUserRecord{}

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

// --- ScrapeLog.go ---

// DeleteScrapeLog deletes a ScrapeLog using a defined ID and column
func (db *dbw) DeleteScrapeLog(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM scrape_log WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadScrapeLog loads a ScrapeLog using a defined ID and column for query
func (db *dbw) LoadScrapeLog(id interface{}, col string) (ScrapeLog, error) {
	query := "SELECT * FROM scrape_log WHERE `" + col + "`=?;"

	data := ScrapeLog{}
	if err := db.Get(&data, query, id); err != nil && err != sql.ErrNoRows {
		return ScrapeLog{}, err
	}

	return data, nil
}

// SaveScrapeLog saves a ScrapeLog to the database
func (db *dbw) SaveScrapeLog(s ScrapeLog) error {
	query := "INSERT INTO scrape_log " +
		"(`info_hash`, `passkey`, `ip`, `time`) " +
		"VALUES (?, ?, ?, UNIX_TIMESTAMP());"

	tx := db.MustBegin()
	tx.Execl(query, s.InfoHash, s.Passkey, s.IP)

	return tx.Commit()
}

// --- UserRecord.go ---

// DeleteUserRecord deletes a UserRecord using a defined ID and column
func (db *dbw) DeleteUserRecord(id interface{}, col string) error {
	tx := db.MustBegin()
	tx.Execl("DELETE FROM users WHERE `"+col+"` = ?", id)

	return tx.Commit()
}

// LoadUserRecord loads a UserRecord using a defined ID and column for query
func (db *dbw) LoadUserRecord(id interface{}, col string) (UserRecord, error) {
	query := "SELECT * FROM users WHERE `" + col + "`=?;"

	data := UserRecord{}
	if err := db.Get(&data, query, id); err != nil && err != sql.ErrNoRows {
		return UserRecord{}, err
	}

	return data, nil
}

// SaveUserRecord saves a UserRecord to the database
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

// --- WhitelistRecord.go ---

// LoadWhitelistRecord loads a WhitelistRecord using a defined ID and column for query
func (db *dbw) LoadWhitelistRecord(id interface{}, col string) (WhitelistRecord, error) {
	query := "SELECT * FROM whitelist WHERE `" + col + "`=?;"

	result := WhitelistRecord{}
	if err := db.Get(&result, query, id); err != nil && err != sql.ErrNoRows {
		return WhitelistRecord{}, err
	}

	return result, nil
}

// SaveWhitelistRecord saves a WhitelistRecord to the database
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
