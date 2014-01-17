// +build ql

package goat

import (
	"io"
	"log"
	"os"
	"os/user"
	"time"

	// Bring in the ql driver
	"github.com/cznic/ql"
)

// ql backend functions, courtesy of Tim Jurcka (sdgoij)
// https://github.com/mdlayher/goat/pull/16

var (
	qlOptions = ql.Options{CanCreate: true}
	qlwdb     *qlw

	qlc = map[string]ql.List{}

	// Map of all queries available to ql
	qlq = map[string]string{
		// announceLog
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

		// apiKey
		"apikey_load_id":      "SELECT id(),user_id,key,salt FROM api_keys WHERE id()==$1",
		"apikey_load_user_id": "SELECT id(),user_id,key,salt FROM api_keys WHERE user_id==$1",
		"apikey_load_key":     "SELECT id(),user_id,key,salt FROM api_keys WHERE key==$1",
		"apikey_insert":       "INSERT INTO api_keys VALUES ($1, $2, $3)",
		"apikey_update":       "UPDATE api_keys key=$2,salt=$3 WHERE id()==$1",

		// fileRecord
		"filerecord_load_all":         "SELECT id(),info_hash,verified,create_time,update_time FROM files",
		"filerecord_load_id":          "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE id()==$1 ORDER BY id()",
		"filerecord_load_info_hash":   "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE info_hash==$1 ORDER BY id()",
		"filerecord_load_verified":    "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE verified==$1 ORDER BY id()",
		"filerecord_load_create_time": "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE create_time==$1 ORDER BY id()",
		"filerecord_load_update_time": "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE update_time==$1 ORDER BY id()",
		"filerecord_insert":           "INSERT INTO files VALUES ($1,$2,now(),now())",
		"filerecord_update":           "UPDATE files verified=$2,update_time=now() WHERE id()==$1",

		// fileUser
		"fileuser_load":            "SELECT * FROM files_users WHERE file_id==$1 && user_id==$2 && ip==$3",
		"fileuser_load_file_id":    "SELECT * FROM files_users WHERE file_id==$1",
		"fileuser_count_completed": "SELECT count(user_id) FROM files_users WHERE file_id==$1 && completed==true && left==0",
		"fileuser_count_seeders":   "SELECT count(user_id) FROM files_users WHERE file_id==$1 && active==true && completed==true && left==0",
		"fileuser_count_leechers":  "SELECT count(user_id) FROM files_users WHERE file_id==$1 && active==true && completed==false && left>0",
		"fileuser_find_peerlist":   "SELECT DISTINCT a.ip, a.port FROM announce_log AS a, (SELECT id() AS id, info_hash FROM files) AS f, (SELECT file_id, ip FROM files_users) AS u WHERE a.ip==u.ip && a.ip!=$2 && f.info_hash==$1",
		"fileuser_find_inactive":   "SELECT user_id, ip FROM files_users WHERE (ts<(now()-$2)) && active==true && file_id==$1",
		"fileuser_mark_inactive":   "UPDATE files_users active=false WHERE file_id==$1 && user_id==$2 && ip==$3",
		"fileuser_insert":          "INSERT INTO files_users VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,now())",
		"fileuser_update":          "UPDATE files_users active=$4,completed=$5,announced=$6,uploaded=$7,downloaded=$8,left=$9,ts=now() WHERE file_id==$1 && user_id==$2 && ip==$3",

		// scrapeLog
		"scrapelog_load_id":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE id()==$1",
		"scrapelog_load_info_hash": "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE info_hash==$1",
		"scrapelog_load_passkey":   "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE passkey==$1",
		"scrapelog_load_ip":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE ip==$1",
		"scrapelog_insert":         "INSERT INTO scrape_log VALUES ($1, $2, $3, now())",

		// userRecord
		"user_load_id":            "SELECT id(),username,passkey,torrent_limit FROM users WHERE id()==$1",
		"user_load_username":      "SELECT id(),username,passkey,torrent_limit FROM users WHERE username==$1",
		"user_load_passkey":       "SELECT id(),username,passkey,torrent_limit FROM users WHERE passkey==$1",
		"user_load_torrent_limit": "SELECT id(),username,passkey,torrent_limit FROM users WHERE torrent_limit==$1",
		"user_insert":             "INSERT INTO users VALUES($1, $2, $3)",
		"user_update":             "UPDATE users username=$2, passkey=$3, torrent_limit=$4 WHERE id()==$1",
		"user_uploaded":           "SELECT sum(uploaded) AS uploaded FROM files_users WHERE user_id==$1",
		"user_downloaded":         "SELECT sum(downloaded) AS downloaded FROM files_users WHERE user_id==$1",
		"user_seeding":            "SELECT count(user_id) AS seeding FROM files_users WHERE user_id==$1 && active==true && completed==true && left==0",
		"user_leeching":           "SELECT count(user_id) AS leeching FROM files_users WHERE user_id==$1 && active==true && completed==false && left>0",

		// whitelistRecord
		"whitelist_load_id":       "SELECT id(),client,approved FROM whitelist WHERE id()==$1",
		"whitelist_load_client":   "SELECT id(),client,approved FROM whitelist WHERE client==$1",
		"whitelist_load_approved": "SELECT id(),client,approved FROM whitelist WHERE approved==$1",
		"whitelist_insert":        "INSERT INTO whitelist VALUES ($1, $2)",
		"whitelist_update":        "UPDATE whitelist client=$2, approved=$3 WHERE id()==$1",
	}
)

// init performs startup routines for database_ql
func init() {
	// dbConnectFunc connects to ql database file
	dbConnectFunc = func() (dbModel, error) {
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
	dbCloseFunc = func() {
		if nil != qlwdb {
			log.Println("Closing ql database")
			qlwdb.Close()
		}
	}

	// dbNameFunc returns the name of this backend
	dbNameFunc = func() string {
		return "ql"
	}
}

// qlw contains a pointer to the ql database
type qlw struct {
	*ql.DB
}

// NewTransaction starts a new ql transaction
func (db *qlw) NewTransaction() qltx {
	tx := qltx{ql.NewRWCtx(), db}
	tx.Run("BEGIN TRANSACTION;")

	return tx
}

// --- announceLog.go ---

// LoadAnnounceLog loads an announceLog using a defined ID and column for query
func (db *qlw) LoadAnnounceLog(id interface{}, col string) (announceLog, error) {
	rs, _, err := qlQuery(db, "announcelog_load_"+col, true, id)

	result := announceLog{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = announceLog{
			ID:         data[0].(int),
			InfoHash:   data[1].(string),
			Passkey:    data[2].(string),
			Key:        data[3].(string),
			IP:         data[4].(string),
			Port:       data[5].(int),
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

// SaveAnnounceLog saves an announceLog to database
func (db *qlw) SaveAnnounceLog(a announceLog) (err error) {
	_, _, err = qlQuery(db, "announcelog_save", true,
		a.InfoHash, a.Passkey, a.Key,
		a.IP, int32(a.Port), a.UDP,
		a.Uploaded, a.Downloaded,
		a.Left, a.Event, a.Client,
		time.Unix(a.Time, 0))

	return
}

// --- apiKey.go ---

// LoadApiKey loads an apiKey using a defined ID and column for query
func (db *qlw) LoadApiKey(id interface{}, col string) (apiKey, error) {
	rs, _, err := qlQuery(db, "apikey_load_"+col, true, id)

	result := apiKey{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = apiKey{
			ID:     int(data[0].(int64)),
			UserID: int(data[1].(int64)),
			Key:    data[2].(string),
			Salt:   data[3].(string),
		}

		return false, nil
	})

	return result, err
}

// SaveApiKey saves an apiKey to the database
func (db *qlw) SaveApiKey(key apiKey) (err error) {
	if k, err := db.LoadApiKey(key.ID, "id"); (k == apiKey{}) && err == nil {
		_, _, err = qlQuery(db, "apikey_insert", true, int64(key.UserID), key.Key, key.Salt)
	} else {
		_, _, err = qlQuery(db, "apikey_update", true, k.ID, key.Key, key.Salt)
	}

	return
}

// --- fileRecord.go ---

// LoadFileRecord loads a fileRecord using a defined ID and column for query
func (db *qlw) LoadFileRecord(id interface{}, col string) (fileRecord, error) {
	// Prevent error cannot convert 1 (type int) to type int64
	if value, ok := id.(int); ok {
		id = int64(value)
	}
	rs, _, err := qlQuery(db, "filerecord_load_"+col, true, id)

	result := fileRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = fileRecord{
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
func (db *qlw) SaveFileRecord(f fileRecord) (err error) {
	if fr, err := db.LoadFileRecord(f.ID, "id"); (fr == fileRecord{}) && err == nil {
		_, _, err = qlQuery(db, "filerecord_insert", true, f.InfoHash, f.Verified)
	} else {
		_, _, err = qlQuery(db, "filerecord_update", true, f.ID, f.Verified)
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

// GetFileRecordPeerList returns a compact peer list containing IP/port pairs
func (db *qlw) GetFileRecordPeerList(infohash, exclude string, limit int) ([]byte, error) {
	rs, _, err := qlQuery(db, "fileuser_find_peerlist", true, infohash, exclude)
	buf := []byte{}

	if err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			buf = append(buf, ip2b(data[0].(string), uint16(data[1].(int32)))...)

			return len(buf)/6 < limit, nil
		})
	}

	return buf, err
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

// GetAllFileRecords returns a list of all fileRecords known to the database
func (db *qlw) GetAllFileRecords() (files []fileRecord, err error) {
	if rs, _, err := qlQuery(db, "filerecord_load_all", false); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			files = append(files, fileRecord{
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

// --- fileUserRecord.go ---

// LoadFileUserRecord loads a fileUserRecord using a defined ID and column for query
func (db *qlw) LoadFileUserRecord(fid, uid int, ip string) (fileUserRecord, error) {
	rs, _, err := qlQuery(db, "fileuser_load", true, int64(fid), int64(uid), ip)

	result := fileUserRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = fileUserRecord{
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

// SaveFileUserRecord saves a fileUserRecord to the database
func (db *qlw) SaveFileUserRecord(f fileUserRecord) (err error) {
	if fr, e := db.LoadFileUserRecord(f.FileID, f.UserID, f.IP); (fr == fileUserRecord{}) {
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

// LoadFileUserRepository loads all fileUserRecords matching a defined ID and column for query
func (db *qlw) LoadFileUserRepository(id interface{}, col string) (files []fileUserRecord, err error) {
	if rs, _, err := qlQuery(db, "fileuser_load_"+col, true, id); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			files = append(files, fileUserRecord{
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

// --- scrapeLog.go ---

// LoadScrapeLog loads a scrapeLog using a defined ID and column for query
func (db *qlw) LoadScrapeLog(id interface{}, col string) (scrape scrapeLog, err error) {
	if rs, _, err := qlQuery(db, "scrapelog_load_"+col, true, id); err == nil && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			scrape = scrapeLog{
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

// SaveScrapeLog saves a scrapeLog to the database
func (db *qlw) SaveScrapeLog(s scrapeLog) (err error) {
	_, _, err = qlQuery(db, "scrapelog_insert", true, s.InfoHash, s.Passkey, s.IP)
	return
}

// --- userRecord.go ---

// LoadUserRecord loads a userRecord using a defined ID and column for query
func (db *qlw) LoadUserRecord(id interface{}, col string) (userRecord, error) {
	rs, _, err := qlQuery(db, "user_load_"+col, true, id)

	result := userRecord{}
	if err != nil {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = userRecord{
			ID:           int(data[0].(int64)),
			Username:     data[1].(string),
			Passkey:      data[2].(string),
			TorrentLimit: int(data[3].(int64)),
		}

		return false, nil
	})

	return result, err
}

// SaveUserRecord saves a userRecord to the database
func (db *qlw) SaveUserRecord(u userRecord) (err error) {
	if user, e := db.LoadUserRecord(u.ID, "id"); (user == userRecord{}) {
		if nil == e {
			_, _, err = qlQuery(db, "user_insert", true,
				u.Username, u.Passkey, u.TorrentLimit)
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "user_update", true,
			user.ID, u.Username, u.Passkey, u.TorrentLimit)
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

// --- whitelistRecord.go ---

// LoadWhitelistRecord loads a whitelistRecord using a defined ID and column for query
func (db *qlw) LoadWhitelistRecord(id interface{}, col string) (whitelistRecord, error) {
	rs, _, err := qlQuery(db, "whitelist_load_"+col, true, id)

	result := whitelistRecord{}
	if err != nil || len(rs) < 1 {
		return result, err
	}

	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = whitelistRecord{
			ID:       int(data[0].(int64)),
			Client:   data[1].(string),
			Approved: data[2].(bool),
		}

		return false, nil
	})

	return result, err
}

// SaveWhitelistRecord saves a whitelistRecord to the database
func (db *qlw) SaveWhitelistRecord(w whitelistRecord) (err error) {
	if wl, e := db.LoadWhitelistRecord(w.ID, "id"); (wl == whitelistRecord{}) {
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
	if list, err := qlCompile(key, wraptx); err == nil {
		return db.Execute(ql.NewRWCtx(), list, arg...)
	} else {
		return []ql.Recordset(nil), 0, err
	}
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
	} else {
		list = l
	}

	return
}

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
	_, _, err = t.Run("COMMIT;")
	return
}

// Rollback performs a database rollback on failed end of transaction
func (t *qltx) Rollback() (err error) {
	_, _, err = t.Run("ROLLBACK;")
	return
}
