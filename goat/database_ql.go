// +build ql

package goat

import (
	"log"
	"time"

	"github.com/cznic/ql"
)

type qlqcc struct {
	c map[string]ql.List
}

var (
	qlOptions = ql.Options{CanCreate: true}
	qlwdb     *qlw

	qlc = map[string]ql.List{}
	qlq = map[string]string{
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

		"apikey_load_id":      "SELECT id(),user_id,key FROM api_keys WHERE id()==$1",
		"apikey_load_user_id": "SELECT id(),user_id,key FROM api_keys WHERE user_id==$1",
		"apikey_load_key":     "SELECT id(),user_id,key FROM api_keys WHERE key==$1",
		"apikey_insert":       "INSERT INTO api_keys VALUES ($1, $2)",
		"apikey_update":       "UPDATE api_keys key=$1 WHERE id()==$1",

		"filerecord_load_all":         "SELECT id(),info_hash,verified,create_time,update_time FROM files",
		"filerecord_load_info_hash":   "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE info_hash==$1 ORDER BY id()",
		"filerecord_load_verified":    "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE verified==$1 ORDER BY id()",
		"filerecord_load_create_time": "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE create_time==$1 ORDER BY id()",
		"filerecord_load_update_time": "SELECT id(),info_hash,verified,create_time,update_time FROM files WHERE update_time==$1 ORDER BY id()",
		"filerecord_insert":           "INSERT INTO files VALUES ($1,$2,now(),now())",
		"filerecord_update":           "UPDATE files verified=$2,update_time=now() WHERE id()==$1",

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

		"scrapelog_load_id":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE id()==$1",
		"scrapelog_load_info_hash": "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE info_hash==$1",
		"scrapelog_load_passkey":   "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE passkey==$1",
		"scrapelog_load_ip":        "SELECT id(),info_hash,passkey,ip,ts FROM scrape_log WHERE ip==$1",
		"scrapelog_insert":         "INSERT INTO scrape_log VALUES ($1, $2, $3, now())",

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

		"whitelist_load_id":       "SELECT id(),client,approved FROM whitelist WHERE id()==$1",
		"whitelist_load_client":   "SELECT id(),client,approved FROM whitelist WHERE client==$1",
		"whitelist_load_approved": "SELECT id(),client,approved FROM whitelist WHERE approved==$1",
		"whitelist_insert":        "INSERT INTO whitelist VALUES ($1, $2)",
		"whitelist_update":        "UPDATE whitelist client=$2, approved=$3 WHERE id()==$1",
	}
)

func init() {
	dbConnectFunc = func() (dbmodel, error) {
		if nil == qlwdb {
			name := static.Config.DB.Database + ".db"
			db, err := ql.OpenFile(name, &qlOptions)
			if nil != err {
				return nil, err
			}
			log.Println("Opened ql database '" + name + "'")
			qlwdb = &qlw{db}
		}
		return qlwdb, nil
	}
	dbCloseFunc = func() {
		if nil != qlwdb {
			log.Println("closing ql database")
			qlwdb.Close()
		}
	}
}

type qlw struct {
	*ql.DB
}

func (db *qlw) NewTransaction() qltx {
	tx := qltx{ql.NewRWCtx(), db}
	tx.Run("BEGIN TRANSACTION;")
	return tx
}

// --- announceLog.go ---

func (db *qlw) LoadAnnounceLog(id interface{}, col string) (announceLog, error) {
	rs, _, err := qlQuery(db, "announcelog_load_"+col, true, id)
	result := announceLog{}
	if err != nil {
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

func (db *qlw) LoadApiKey(id interface{}, col string) (apiKey, error) {
	rs, _, err := qlQuery(db, "apikey_load_"+col, true, id)
	result := apiKey{}
	if err != nil {
		return result, err
	}
	err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
		result = apiKey{
			ID:     int(data[0].(int64)),
			UserID: int(data[1].(int64)),
			Key:    data[2].(string),
		}
		return false, nil
	})
	return result, err
}

func (db *qlw) SaveApiKey(key apiKey) (err error) {
	if k, e := db.LoadApiKey(key.ID, "id"); (k == apiKey{}) {
		if nil == e {
			_, _, err = qlQuery(db, "apikey_insert", true, int64(key.UserID), key.Key)
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "apikey_update", true, k.ID, key.Key)
	}
	return
}

// --- fileRecord.go ---

func (db *qlw) LoadFileRecord(id interface{}, col string) (fileRecord, error) {
	rs, _, err := qlQuery(db, "filerecord_load_"+col, true, id)
	result := fileRecord{}
	if err != nil {
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

func (db *qlw) SaveFileRecord(f fileRecord) (err error) {
	if fr, e := db.LoadFileRecord(f.ID, "id"); (fr == fileRecord{}) {
		if nil == e {
			_, _, err = qlQuery(db, "filerecord_insert", true, f.InfoHash, f.Verified)
		} else {
			err = e
		}
	} else {
		_, _, err = qlQuery(db, "filerecord_update", true, f.ID, f.Verified)
	}
	return
}

func (db *qlw) CountFileRecordCompleted(id int) (int, error) {
	rs, _, err := qlQuery(db, "fileuser_count_completed", true, id)
	completed := int(0)

	if nil == err && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			completed = int(data[0].(int64))
			return false, err
		})
	}
	return completed, err
}

func (db *qlw) CountFileRecordSeeders(id int) (int, error) {
	rs, _, err := qlQuery(db, "fileuser_count_seeders", true, int64(id))
	seeders := int(0)

	if nil == err && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			seeders = int(data[0].(int64))
			return false, nil
		})
	}
	return seeders, err
}

func (db *qlw) CountFileRecordLeechers(id int) (int, error) {
	rs, _, err := qlQuery(db, "fileuser_count_leechers", true, int64(id))
	leechers := int(0)

	if nil == err && len(rs) > 0 {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			leechers = int(data[0].(int64))
			return false, nil
		})
	}
	return leechers, err
}

func (db *qlw) GetFileRecordPeerList(infohash, exclude string, limit int) ([]byte, error) {
	rs, _, err := qlQuery(db, "fileuser_find_peerlist", true, infohash, exclude)
	buf := []byte{}

	if nil == err {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			buf = append(buf, ip2b(data[0].(string), uint16(data[1].(int32)))...)
			return len(buf)/6 < limit, nil
		})
	}
	return buf, err
}

func (db *qlw) GetInactiveUserInfo(fid int, interval time.Duration) (users []userinfo, err error) {
	if rs, _, e := qlQuery(db, "fileuser_find_inactive", true, int64(fid), interval); nil == e {
		err = rs[0].Do(false, func(data []interface{}) (bool, error) {
			users = append(users, userinfo{int(data[0].(int64)), data[1].(string)})
			return true, nil
		})
	} else {
		err = e
	}
	return
}

func (db *qlw) MarkFileUsersInactive(fid int, users []userinfo) (err error) {
	if list, e := qlCompile("fileuser_mark_inactive", false); nil == err {
		tx := db.NewTransaction()
		for _, user := range users {
			if _, _, err = tx.Execute(list, int64(fid), int64(user.UserID), user.IP); nil != err {
				tx.Rollback()
				return
			}
		}
		err = tx.Commit()
	} else {
		err = e
	}
	return
}

func (db *qlw) GetAllFileRecords() (files []fileRecord, err error) {
	if rs, _, err := qlQuery(db, "fileuser_load_all", false); nil == err {
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

func (db *qlw) LoadFileUserRecord(fid, uid int, ip string) (fileUserRecord, error) {
	rs, _, err := qlQuery(db, "fileuser_load", true, int64(fid), int64(uid), ip)
	result := fileUserRecord{}
	if err != nil {
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

func (db *qlw) LoadFileUserRepository(id interface{}, col string) (files []fileUserRecord, err error) {
	if rs, _, err := qlQuery(db, "fileuser_load_"+col, true, id); nil == err {
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

func (db *qlw) LoadScrapeLog(id interface{}, col string) (scrape scrapeLog, err error) {
	if rs, _, err := qlQuery(db, "scrapelog_load_"+col, true, id); nil == err {
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

func (db *qlw) SaveScrapeLog(s scrapeLog) (err error) {
	_, _, err = qlQuery(db, "scrapelog_insert", true, s.InfoHash, s.Passkey, s.IP)
	return
}

// --- userRecord.go ---

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
			TorrentLimit: data[3].(int),
		}
		return false, nil
	})
	return result, err
}

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

func (db *qlw) GetUserUploaded(uid int) (int64, error) {
	return qlQueryI64(db, "user_uploaded", uid)
}

func (db *qlw) GetUserDownloaded(uid int) (int64, error) {
	return qlQueryI64(db, "user_downloaded", uid)
}

func (db *qlw) GetUserSeeding(uid int) (int, error) {
	i, err := qlQueryI64(db, "user_seeding", uid)
	return int(i), err
}

func (db *qlw) GetUserLeeching(uid int) (int, error) {
	i, err := qlQueryI64(db, "user_leeching", uid)
	return int(i), err
}

// --- whitelistRecord.go ---

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

func qlQuery(db *qlw, key string, wraptx bool, arg ...interface{}) ([]ql.Recordset, int, error) {
	if list, err := qlCompile(key, wraptx); nil == err {
		return db.Execute(ql.NewRWCtx(), list, arg...)
	} else {
		return []ql.Recordset(nil), 0, err
	}
}

func qlQueryI64(db *qlw, key string, arg ...interface{}) (i int64, err error) {
	if rs, _, err := qlQuery(db, key, false, arg...); nil == err {
		err = rs[len(rs)-1].Do(false, func(data []interface{}) (bool, error) {
			i = data[0].(int64)
			return false, nil
		})
	}
	return
}

func qlCompile(key string, wraptx bool) (list ql.List, err error) {
	var src string
	if l, ok := qlc[key]; !ok {
		if src, ok = qlq[key]; !ok {
			src = key
		}
		if wraptx {
			src = "BEGIN TRANSACTION; " + src + "; COMMIT;"
		}
		if l, e := ql.Compile(src); nil != err {
			err = e
		} else {
			list = l
		}
	} else {
		list = l
	}
	return
}

type qltx struct {
	ctx *ql.TCtx
	db  *qlw
}

func (t *qltx) Execute(list ql.List, arg ...interface{}) ([]ql.Recordset, int, error) {
	return t.db.Execute(t.ctx, list, arg...)
}

func (t *qltx) Run(src string, arg ...interface{}) ([]ql.Recordset, int, error) {
	return t.db.Run(t.ctx, src, arg...)
}

func (t *qltx) Commit() (err error) {
	_, _, err = t.Run("COMMIT;")
	return
}

func (t *qltx) Rollback() (err error) {
	_, _, err = t.Run("ROLLBACK;")
	return
}
