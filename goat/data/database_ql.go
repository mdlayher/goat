// +build ql

package data

import (
	"io"
	"log"
	"os"
	"os/user"
	ospath "path"
	"time"

	"github.com/mdlayher/goat/goat/common"

	// Bring in the ql driver
	"github.com/cznic/ql"
)

// ql backend functions, courtesy of Tim Jurcka (sdgoij)
// https://github.com/mdlayher/goat/pull/16

var (
	qlOptions = ql.Options{CanCreate: true}
	qlwdb     *qlw

	// Map of compiled ql queries
	qlc = map[string]ql.List{}

	// Map of all queries available to ql
	qlq = map[string]string{
		// AnnounceLog
		"announcelog_delete_id":       "DELETE FROM announce_log WHERE id()==$1",
		"announcelog_load_id":         "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE id()==$1 ORDER BY id()",
		"announcelog_load_info_hash":  "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE info_hash==$1 ORDER BY id()",
		"announcelog_load_passkey":    "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE passkey==$1 ORDER BY id()",
		"announcelog_load_key":        "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE key==$1 ORDER BY id()",
		"announcelog_load_ip":         "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE ip==$1 ORDER BY id()",
		"announcelog_load_port":       "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE port==$1 ORDER BY id()",
		"announcelog_load_udp":        "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE udp==$1 ORDER BY id()",
		"announcelog_load_uploaded":   "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE uploaded==$1 ORDER BY id()",
		"announcelog_load_downloaded": "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE downloaded==$1 ORDER BY id()",
		"announcelog_load_left":       "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE left==$1 ORDER BY id()",
		"announcelog_load_event":      "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE event==$1 ORDER BY id()",
		"announcelog_load_client":     "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE client==$1 ORDER BY id()",
		"announcelog_load_time":       "SELECT id(),info_hash,passkey,key,ip,port,udp,uploaded,downloaded,left,event,client,ts FROM announce_log WHERE time==$1 ORDER BY id()",
		"announcelog_save":            "INSERT INTO announce_log VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,now());",

		// APIKey
		"apikey_delete_id":    "DELETE FROM api_keys WHERE id()==$1",
		"apikey_delete_key":   "DELETE FROM api_keys WHERE key==$1",
		"apikey_load_id":      "SELECT id(),user_id,key,expire FROM api_keys WHERE id()==$1",
		"apikey_load_user_id": "SELECT id(),user_id,key,expire FROM api_keys WHERE user_id==$1",
		"apikey_load_key":     "SELECT id(),user_id,key,expire FROM api_keys WHERE key==$1",
		"apikey_insert":       "INSERT INTO api_keys VALUES ($1, $2, $3)",
		"apikey_update":       "UPDATE api_keys key=$2,expire=$3 WHERE id()==$1",

		// FileRecord
		"filerecord_delete_id":          "DELETE FROM files WHERE id()==$1",
		"filerecord_delete_info_hash":   "DELETE FROM files WHERE info_hash==$1",
		"filerecord_find_peerlist_http": "SELECT DISTINCT a.ip, a.port FROM announce_log AS a, (SELECT id() AS id, info_hash FROM files) AS f, (SELECT file_id, ip FROM files_users) AS u WHERE a.ip==u.ip && (now()-$1) <= a.time && f.info_hash==$2",
		"filerecord_find_peerlist_udp":  "SELECT DISTINCT a.ip, a.port FROM announce_log AS a, (SELECT id() AS id, info_hash FROM files) AS f, WHERE (now()-$1) <= a.time && f.info_hash==$2",
		"filerecord_load_all":           "SELECT id(),info_hash,verified,create_time,update_time FROM files",
		"filerecord_load_id":            "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE id()==$1 ORDER BY id()",
		"filerecord_load_info_hash":     "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE info_hash==$1 ORDER BY id()",
		"filerecord_load_verified":      "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE verified==$1 ORDER BY id()",
		"filerecord_load_create_time":   "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE create_time==$1 ORDER BY id()",
		"filerecord_load_update_time":   "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE update_time==$1 ORDER BY id()",
		"filerecord_insert":             "INSERT INTO files VALUES ($1,$2,now(),now())",
		"filerecord_update":             "UPDATE files verified=$2,update_time=now() WHERE id()==$1",

		// fileUser
		"fileuser_delete":          "DELETE FROM files_users WHERE file_id==$1 && user_id==$2 && ip==$3",
		"fileuser_load":            "SELECT * FROM files_users WHERE file_id==$1 && user_id==$2 && ip==$3",
		"fileuser_load_file_id":    "SELECT * FROM files_users WHERE file_id==$1",
		"fileuser_count_completed": "SELECT count(user_id) FROM files_users WHERE file_id==$1 && completed==true && left==0",
		"fileuser_count_seeders":   "SELECT count(user_id) FROM files_users WHERE file_id==$1 && active==true && completed==true && left==0",
		"fileuser_count_leechers":  "SELECT count(user_id) FROM files_users WHERE file_id==$1 && active==true && completed==false && left>0",
		"fileuser_find_inactive":   "SELECT user_id, ip FROM files_users WHERE (ts<(now()-$2)) && active==true && file_id==$1",
		"fileuser_mark_inactive":   "UPDATE files_users active=false WHERE file_id==$1 && user_id==$2 && ip==$3",
		"fileuser_insert":          "INSERT INTO files_users VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,now())",
		"fileuser_update":          "UPDATE files_users active=$4,completed=$5,announced=$6,uploaded=$7,downloaded=$8,left=$9,ts=now() WHERE file_id==$1 && user_id==$2 && ip==$3",

		// ScrapeLog
		"scrapelog_delete_id":      "DELETE FROM scrape_log WHERE id()==$1",
		"scrapelog_load_id":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE id()==$1",
		"scrapelog_load_info_hash": "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE info_hash==$1",
		"scrapelog_load_passkey":   "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE passkey==$1",
		"scrapelog_load_ip":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE ip==$1",
		"scrapelog_insert":         "INSERT INTO scrape_log VALUES ($1, $2, $3, now())",

		// UserRecord
		"user_delete_username":    "DELETE FROM users WHERE username==$1",
		"user_load_all":           "SELECT id(),username,password,passkey,torrent_limit FROM users",
		"user_load_id":            "SELECT id(),username,password,passkey,torrent_limit FROM users WHERE id()==$1",
		"user_load_username":      "SELECT id(),username,password,passkey,torrent_limit FROM users WHERE username==$1",
		"user_load_password":      "SELECT id(),username,password,passkey,torrent_limit FROM users WHERE password==$1",
		"user_load_passkey":       "SELECT id(),username,password,passkey,torrent_limit FROM users WHERE passkey==$1",
		"user_load_torrent_limit": "SELECT id(),username,password,passkey,torrent_limit FROM users WHERE torrent_limit==$1",
		"user_insert":             "INSERT INTO users VALUES($1, $2, $3, $4)",
		"user_update":             "UPDATE users username=$2, password=$3, passkey=$4, torrent_limit=$5 WHERE id()==$1",
		"user_uploaded":           "SELECT sum(uploaded) AS uploaded FROM files_users WHERE user_id==$1",
		"user_downloaded":         "SELECT sum(downloaded) AS downloaded FROM files_users WHERE user_id==$1",
		"user_seeding":            "SELECT count(user_id) AS seeding FROM files_users WHERE user_id==$1 && active==true && completed==true && left==0",
		"user_leeching":           "SELECT count(user_id) AS leeching FROM files_users WHERE user_id==$1 && active==true && completed==false && left>0",

		// WhitelistRecord
		"whitelist_delete_client": "DELETE FROM whitelist WHERE client==$1",
		"whitelist_load_id":       "SELECT id(),client,approved FROM whitelist WHERE id()==$1",
		"whitelist_load_client":   "SELECT id(),client,approved FROM whitelist WHERE client==$1",
		"whitelist_load_approved": "SELECT id(),client,approved FROM whitelist WHERE approved==$1",
		"whitelist_insert":        "INSERT INTO whitelist VALUES ($1, $2)",
		"whitelist_update":        "UPDATE whitelist client=$2, approved=$3 WHERE id()==$1",
	}
)

// init performs startup routines for database_ql
func init() {
	// DBConnectFunc connects to ql database file
	DBConnectFunc = func() (dbModel, error) {
		if nil == qlwdb {
			// Database name
			name := "goat.db"

			// Load current user from OS, to get home directory
			var path string
			user, err := user.Current()
			if err != nil {
				log.Println(err.Error())
				path = "./"
			} else {
				// Store config in standard location
				path = user.HomeDir + "/.config/goat/"
			}

			// Allow manul override of db path, if flag is set
			if QLDBPath != nil && *QLDBPath != "" {
				// Split db path into path and filename
				path = ospath.Dir(*QLDBPath) + "/"
				name = ospath.Base(*QLDBPath)
			}

			log.Println("Loading ql database: " + path + name)

			// Check file existence
			_, err = os.Stat(path + name)
			if err != nil {
				if os.IsNotExist(err) {
					log.Println("Could not find ql database, attempting to create it...")

					err = os.MkdirAll(path, 0775)
					if err != nil {
						log.Println("Failed to create directory: " + path)
					}

					// Attempt to copy database file to home directory
					source, err := os.Open("./res/ql/" + name)
					if err != nil {
						log.Println("Failed to read source file: " + name)
					}

					// Open destination file
					dest, err := os.Create(path + name)
					if err != nil {
						log.Println("Failed to create destination file: " + path + name)
					}

					// Copy contents
					_, err = io.Copy(dest, source)
					if err != nil {
						log.Println("Failed to copy to database file: " + path + name)
					}

					// Close files
					source.Close()
					dest.Close()
				}
			}

			db, err := ql.OpenFile(path+name, &qlOptions)
			if err != nil {
				return nil, err
			}

			qlwdb = &qlw{db}
		}

		return qlwdb, nil
	}

	// Generate connection string using configuration
	DBCloseFunc = func() {
		if nil != qlwdb {
			qlwdb.Close()
		}
	}

	// DBNameFunc returns the name of this backend
	DBNameFunc = func() string {
		return "ql"
	}
}

// qlw contains a pointer to the ql database
type qlw struct {
	*ql.DB
}

// Close closes the ql database
func (db *qlw) Close() error {
	return nil
}

// NewTransaction starts a new ql transaction
func (db *qlw) NewTransaction() qltx {
	tx := qltx{ql.NewRWCtx(), db}
	tx.Execute(qlBeginTransaction)

	return tx
}

// --- AnnounceLog.go ---

// DeleteAnnounceLog deletes an AnnounceLog using a defined ID and column for query
func (db *qlw) DeleteAnnounceLog(id interface{}, col string) (err error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok {
		id = int64(value)
	}
	_, _, err = qlQuery(db, "announcelog_delete_"+col, true, id)
	return
}

// LoadAnnounceLog loads an AnnounceLog using a defined ID and column for query
func (db *qlw) LoadAnnounceLog(id interface{}, col string) (AnnounceLog, error) {
	rs, _, err := qlQuery(db, "announcelog_load_"+col, true, id)

	result := AnnounceLog{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = AnnounceLog{
			ID:         int(data[0].(int64)),
			InfoHash:   data[1].(string),
			Passkey:    data[2].(string),
			Key:        data[3].(string),
			IP:         data[4].(string),
			Port:       int(data[5].(int32)),
			UDP:        data[6].(bool),
			Uploaded:   data[7].(int64),
			Downloaded: data[8].(int64),
			Left:       data[9].(int64),
			Event:      data[10].(string),
			Client:     data[11].(string),
			Time:       data[12].(time.Time).Unix(),
		}

		return false, nil
	})

	return result, err
}

// SaveAnnounceLog saves an AnnounceLog to database
func (db *qlw) SaveAnnounceLog(a AnnounceLog) (err error) {
	_, _, err = qlQuery(db, "announcelog_save", true,
		a.InfoHash, a.Passkey, a.Key,
		a.IP, int32(a.Port), a.UDP,
		a.Uploaded, a.Downloaded,
		a.Left, a.Event, a.Client,
		time.Unix(a.Time, 0))

	return
}

// --- APIKey.go ---

// DeleteAPIKey deletes an AnnounceLog using a defined ID and column for query
func (db *qlw) DeleteAPIKey(id interface{}, col string) (err error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok && col == "id" {
		id = int64(value)
	}
	_, _, err = qlQuery(db, "apikey_delete_"+col, true, id)
	return
}

// LoadAPIKey loads an APIKey using a defined ID and column for query
func (db *qlw) LoadAPIKey(id interface{}, col string) (APIKey, error) {
	rs, _, err := qlQuery(db, "apikey_load_"+col, true, id)

	result := APIKey{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = APIKey{
			ID:     int(data[0].(int64)),
			UserID: int(data[1].(int64)),
			Key:    data[2].(string),
			Expire: data[3].(int64),
		}

		return false, nil
	})

	return result, err
}

// SaveApiKey saves an apiKey to the database
func (db *qlw) SaveAPIKey(key APIKey) (err error) {
	if k, _ := db.LoadAPIKey(key.ID, "id"); (k == APIKey{}) && err == nil {
		_, _, err = qlQuery(db, "apikey_insert", true, int64(key.UserID), key.Key, key.Expire)
	} else {
		_, _, err = qlQuery(db, "apikey_update", true, int64(k.ID), key.Key, key.Expire)
	}

	return
}

// --- FileRecord.go ---

// DeleteFileRecord deletes an AnnounceLog using a defined ID and column for query
func (db *qlw) DeleteFileRecord(id interface{}, col string) (err error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok && col == "id" {
		id = int64(value)
	}
	_, _, err = qlQuery(db, "filerecord_delete_"+col, true, id)
	return
}

// LoadFileRecord loads a FileRecord using a defined ID and column for query
func (db *qlw) LoadFileRecord(id interface{}, col string) (FileRecord, error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok && col == "id" {
		id = int64(value)
	}
	rs, _, err := qlQuery(db, "filerecord_load_"+col, true, id)

	result := FileRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = FileRecord{
			ID:         int(data[0].(int64)),
			InfoHash:   data[1].(string),
			Verified:   data[2].(bool),
			CreateTime: data[3].(time.Time).Unix(),
			UpdateTime: data[4].(time.Time).Unix(),
		}

		return false, nil
	})

	return result, err
}

// SaveFileRecord saves a fileRecord to the database
func (db *qlw) SaveFileRecord(f FileRecord) (err error) {
	if fr, _ := db.LoadFileRecord(f.ID, "id"); (fr == FileRecord{}) && err == nil {
		_, _, err = qlQuery(db, "filerecord_insert", true, f.InfoHash, f.Verified)
	} else {
		_, _, err = qlQuery(db, "filerecord_update", true, int64(f.ID), f.Verified)
	}

	return
}

// CountFileRecordCompleted counts the number of peers who have completed this file
func (db *qlw) CountFileRecordCompleted(id int) (int, error) {
	completed, err := qlQueryI64(db, "fileuser_count_completed", int64(id))
	return int(completed), err
}

// CountFileRecordSeeders counts the number of peers who are actively seeding this file
func (db *qlw) CountFileRecordSeeders(id int) (int, error) {
	seeders, err := qlQueryI64(db, "fileuser_count_seeders", int64(id))
	return int(seeders), err
}

// CountFileRecordLeechers counts the number of peers who are actively leeching this file
func (db *qlw) CountFileRecordLeechers(id int) (int, error) {
	leechers, err := qlQueryI64(db, "fileuser_count_leechers", int64(id))
	return int(leechers), err
}

// GetFileRecordPeerList returns a list of Peers
func (db *qlw) GetFileRecordPeerList(infoHash string, limit int, http bool) ([]Peer, error) {
	// Select query using HTTP bool
	var query string
	if http {
		query = "filerecord_find_peerlist_http"
	} else {
		query = "filerecord_find_peerlist_udp"
	}

	rs, _, err := qlQuery(db, query, true, common.Static.Config.Interval, infoHash)

	// Generate peer list
	peers := make([]Peer, 0)

	if err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			peer := Peer{
				IP:   data[0].(string),
				Port: uint16(data[1].(int32)),
			}

			peers = append(peers[:], peer)

			return len(peers) < limit, nil
		})
	}

	return peers, err
}

// GetInactiveUserInfo returns a list of users who have not been active for the specified time interval
func (db *qlw) GetInactiveUserInfo(fid int, interval time.Duration) (users []peerInfo, err error) {
	if rs, _, err := qlQuery(db, "fileuser_find_inactive", true, int64(fid), interval); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			users = append(users, peerInfo{int(data[0].(int64)), data[1].(string)})

			return true, nil
		})
	}

	return
}

// MarkFileUsersInactive sets users to be inactive once they have been reaped
func (db *qlw) MarkFileUsersInactive(fid int, users []peerInfo) (err error) {
	if list, err := qlCompile("fileuser_mark_inactive", false); err == nil {
		tx := db.NewTransaction()

		for _, user := range users {
			if _, _, err = tx.Execute(list, int64(fid), int64(user.UserID), user.IP); err != nil {
				tx.Rollback()

				return err
			}
		}
		err = tx.Commit()
	}

	return
}

// GetAllFileRecords returns a list of all FileRecords known to the database
func (db *qlw) GetAllFileRecords() (files []FileRecord, err error) {
	if rs, _, err := qlQuery(db, "filerecord_load_all", false); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			files = append(files, FileRecord{
				ID:         int(data[0].(int64)),
				InfoHash:   data[1].(string),
				Verified:   data[2].(bool),
				CreateTime: data[3].(time.Time).Unix(),
				UpdateTime: data[4].(time.Time).Unix(),
			})

			return true, nil
		})
	}

	return
}

// --- FileUserRecord.go ---

// DeleteFileUserRecord deletes an AnnounceLog using a file ID, user ID, and IP triple
func (db *qlw) DeleteFileUserRecord(fid, uid int, ip string) (err error) {
	_, _, err = qlQuery(db, "fileuser_delete", true, int64(fid), int64(uid), ip)
	return
}

// LoadFileUserRecord loads a FileUserRecord using a file ID, user ID, and IP triple
func (db *qlw) LoadFileUserRecord(fid, uid int, ip string) (FileUserRecord, error) {
	rs, _, err := qlQuery(db, "fileuser_load", true, int64(fid), int64(uid), ip)

	result := FileUserRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = FileUserRecord{
			FileID:     int(data[0].(int64)),
			UserID:     int(data[1].(int64)),
			IP:         data[2].(string),
			Active:     data[3].(bool),
			Completed:  data[4].(bool),
			Announced:  int(data[5].(int64)),
			Uploaded:   data[6].(int64),
			Downloaded: data[7].(int64),
			Left:       data[8].(int64),
			Time:       data[9].(time.Time).Unix(),
		}

		return false, nil
	})

	return result, err
}

// SaveFileUserRecord saves a FileUserRecord to the database
func (db *qlw) SaveFileUserRecord(f FileUserRecord) (err error) {
	if fr, e := db.LoadFileUserRecord(f.FileID, f.UserID, f.IP); (fr == FileUserRecord{}) {
		if nil == e {
			_, _, err = qlQuery(db, "fileuser_insert", true,
				int64(f.FileID), int64(f.UserID), f.IP,
				f.Active, f.Completed, int64(f.Announced),
				f.Uploaded, f.Downloaded, f.Left,
				time.Unix(f.Time, 0))
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "fileuser_update", true,
			int64(f.FileID), int64(f.UserID), f.IP,
			f.Active, f.Completed, int64(f.Announced),
			f.Uploaded, f.Downloaded, f.Left)
	}

	return
}

// LoadFileUserRepository loads all FileUserRecords matching a defined ID and column for query
func (db *qlw) LoadFileUserRepository(id interface{}, col string) (files []FileUserRecord, err error) {
	if rs, _, err := qlQuery(db, "fileuser_load_"+col, true, id); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			files = append(files, FileUserRecord{
				FileID:     int(data[0].(int64)),
				UserID:     int(data[1].(int64)),
				IP:         data[2].(string),
				Active:     data[3].(bool),
				Completed:  data[4].(bool),
				Announced:  data[5].(int),
				Uploaded:   data[6].(int64),
				Downloaded: data[7].(int64),
				Left:       data[8].(int64),
				Time:       data[9].(time.Time).Unix(),
			})

			return false, nil
		})
	}

	return
}

// --- ScrapeLog.go ---

// DeleteScrapeLog deletes an ScrapeLog using a defined ID and column for query
func (db *qlw) DeleteScrapeLog(id interface{}, col string) (err error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok {
		id = int64(value)
	}
	_, _, err = qlQuery(db, "scrapelog_delete_"+col, true, id)
	return
}

// LoadScrapeLog loads a ScrapeLog using a defined ID and column for query
func (db *qlw) LoadScrapeLog(id interface{}, col string) (scrape ScrapeLog, err error) {
	if rs, _, err := qlQuery(db, "scrapelog_load_"+col, true, id); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			scrape = ScrapeLog{
				ID:       int(data[0].(int64)),
				InfoHash: data[1].(string),
				Passkey:  data[2].(string),
				IP:       data[3].(string),
				Time:     data[4].(time.Time).Unix(),
			}

			return false, nil
		})
	}

	return
}

// SaveScrapeLog saves a ScrapeLog to the database
func (db *qlw) SaveScrapeLog(s ScrapeLog) (err error) {
	_, _, err = qlQuery(db, "scrapelog_insert", true, s.InfoHash, s.Passkey, s.IP)
	return
}

// --- UserRecord.go ---

// DeleteUserRecord deletes an AnnounceLog using a defined ID and column for query
func (db *qlw) DeleteUserRecord(id interface{}, col string) (err error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok {
		id = int64(value)
	}

	_, _, err = qlQuery(db, "user_delete_"+col, true, id)
	return
}

// LoadUserRecord loads a UserRecord using a defined ID and column for query
func (db *qlw) LoadUserRecord(id interface{}, col string) (UserRecord, error) {
	rs, _, err := qlQuery(db, "user_load_"+col, true, id)

	result := UserRecord{}
	if err != nil {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = UserRecord{
			ID:           int(data[0].(int64)),
			Username:     data[1].(string),
			Password:     data[2].(string),
			Passkey:      data[3].(string),
			TorrentLimit: int(data[4].(int64)),
		}

		return false, nil
	})

	return result, err
}

// SaveUserRecord saves a userRecord to the database
func (db *qlw) SaveUserRecord(u UserRecord) (err error) {
	if user, e := db.LoadUserRecord(int64(u.ID), "id"); (user == UserRecord{}) {
		if nil == e {
			_, _, err = qlQuery(db, "user_insert", true,
				u.Username, u.Password, u.Passkey, int64(u.TorrentLimit))
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "user_update", true,
			int64(user.ID), u.Username, u.Password, u.Passkey, int64(u.TorrentLimit))
	}

	return
}

// GetUserUploaded calculates the total number of bytes this user has uploaded
func (db *qlw) GetUserUploaded(uid int) (int64, error) {
	return qlQueryI64(db, "user_uploaded", uid)
}

// GetUserDownloaded calculates the total number of bytes this user has downloaded
func (db *qlw) GetUserDownloaded(uid int) (int64, error) {
	return qlQueryI64(db, "user_downloaded", uid)
}

// GetUserSeeding calculates the total number of files this user is actively seeding
func (db *qlw) GetUserSeeding(uid int) (int, error) {
	i, err := qlQueryI64(db, "user_seeding", uid)
	return int(i), err
}

// GetUserLeeching calculates the total number of files this user is actively leeching
func (db *qlw) GetUserLeeching(uid int) (int, error) {
	i, err := qlQueryI64(db, "user_leeching", uid)
	return int(i), err
}

// GetAllUserRecords returns a list of all UserRecords known to the database
func (db *qlw) GetAllUserRecords() (users []UserRecord, err error) {
	if rs, _, err := qlQuery(db, "user_load_all", false); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			users = append(users, UserRecord{
				ID:           int(data[0].(int64)),
				Username:     data[1].(string),
				Password:     data[2].(string),
				Passkey:      data[3].(string),
				TorrentLimit: int(data[4].(int64)),
			})

			return true, nil
		})
	}

	return
}

// --- WhitelistRecord.go ---

// DeleteWhitelistRecord deletes a WhitelistRecord using a defined column and ID
func (db *qlw) DeleteWhitelistRecord(id interface{}, col string) (err error) {
	_, _, err = qlQuery(db, "whitelist_delete_"+col, true, id)
	return
}

// LoadWhitelistRecord loads a WhitelistRecord using a defined ID and column for query
func (db *qlw) LoadWhitelistRecord(id interface{}, col string) (WhitelistRecord, error) {
	rs, _, err := qlQuery(db, "whitelist_load_"+col, true, id)

	result := WhitelistRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = WhitelistRecord{
			ID:       int(data[0].(int64)),
			Client:   data[1].(string),
			Approved: data[2].(bool),
		}

		return false, nil
	})

	return result, err
}

// SaveWhitelistRecord saves a WhitelistRecord to the database
func (db *qlw) SaveWhitelistRecord(w WhitelistRecord) (err error) {
	if wl, e := db.LoadWhitelistRecord(w.ID, "id"); (wl == WhitelistRecord{}) {
		if nil == e {
			_, _, err = qlQuery(db, "whitelist_insert", true,
				w.Client, w.Approved)
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "whitelist_update", true,
			w.ID, w.Client, w.Approved)
	}

	return
}

// qlQuery provides a wrapper to compile a ql query
func qlQuery(db *qlw, key string, wraptx bool, arg ...interface{}) ([]ql.Recordset, int, error) {
	var err error
	if list, err := qlCompile(key, wraptx); err == nil {
		return db.Execute(ql.NewRWCtx(), list, arg...)
	}

	return []ql.Recordset(nil), 0, err
}

// qlQueryI64 provides a wrapper to return int64 values from ql
func qlQueryI64(db *qlw, key string, arg ...interface{}) (i int64, err error) {
	if rs, _, err := qlQuery(db, key, false, arg...); err == nil && len(rs) > 0 {
		err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
			i = data[0].(int64)

			return false, nil
		})
	}

	return
}

// qlCompile provides a wrapper to create safe, transaction-encased queries
func qlCompile(key string, wraptx bool) (list ql.List, err error) {
	var src string
	if l, ok := qlc[key]; !ok {
		if src, ok = qlq[key]; !ok {
			src = key
		}
		if wraptx {
			src = "BEGIN TRANSACTION; " + src + "; COMMIT;"
		}
		if l, e := ql.Compile(src); err != nil {
			err = e
		} else {
			list = l
		}
		qlc[key] = list
	} else {
		list = l
	}

	return
}

// Pre-compiled begin transaction, commit and rollback statements
var (
	qlBeginTransaction, _ = ql.Compile("BEGIN TRANSACTION;")
	qlCommit, _           = ql.Compile("COMMIT;")
	qlRollback, _         = ql.Compile("ROLLBACK;")
)

// qltx contains a ql context and database connection
type qltx struct {
	ctx *ql.TCtx
	db  *qlw
}

// Execute allows for execution via ql context in a transaction
func (t *qltx) Execute(list ql.List, arg ...interface{}) ([]ql.Recordset, int, error) {
	return t.db.Execute(t.ctx, list, arg...)
}

// Run allows for execution without a transaction
func (t *qltx) Run(src string, arg ...interface{}) ([]ql.Recordset, int, error) {
	return t.db.Run(t.ctx, src, arg...)
}

// Commit performs a database commit at the end of a transaction
func (t *qltx) Commit() (err error) {
	_, _, err = t.Execute(qlCommit)
	return
}

// Rollback performs a database rollback on failed end of transaction
func (t *qltx) Rollback() (err error) {
	_, _, err = t.Execute(qlRollback)
	return
}
