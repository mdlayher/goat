package data

import (
	"encoding/json"
	"time"

	"github.com/mdlayher/goat/goat/common"
)

// FileRecord represents a file tracked by tracker
type FileRecord struct {
	ID         int    `json:"id"`
	InfoHash   string `db:"info_hash" json:"infoHash"`
	Verified   bool   `json:"verified"`
	CreateTime int64  `db:"create_time" json:"createTime"`
	UpdateTime int64  `db:"update_time" json:"updateTime"`
}

// jsonFileRecord represents output FileRecord JSON for API
type jsonFileRecord struct {
	ID         int              `json:"id"`
	InfoHash   string           `json:"infoHash"`
	Verified   bool             `json:"verified"`
	CreateTime int64            `json:"createTime"`
	UpdateTime int64            `json:"updateTime"`
	Completed  int              `json:"completed"`
	Seeders    int              `json:"seeders"`
	Leechers   int              `json:"leechers"`
	FileUsers  []FileUserRecord `json:"fileUsers"`
}

// peerInfo represents a peer which will be marked as active or not
type peerInfo struct {
	UserID int `db:"user_id"`
	IP     string
}

// ToJSON converts a FileRecord to a jsonFileRecord struct
func (f FileRecord) ToJSON() ([]byte, error) {
	// Convert all standard fields to the JSON equivalent struct
	j := jsonFileRecord{}
	j.ID = f.ID
	j.InfoHash = f.InfoHash
	j.Verified = f.Verified
	j.CreateTime = f.CreateTime
	j.UpdateTime = f.UpdateTime

	// Load in FileUserRecords associated with this file
	j.FileUsers = f.Users()

	// Load counts for completions, seeding, leeching
	var err error
	j.Completed, err = f.Completed()
	if err != nil {
		return nil, err
	}

	j.Seeders, err = f.Seeders()
	if err != nil {
		return nil, err
	}

	j.Leechers, err = f.Leechers()
	if err != nil {
		return nil, err
	}

	// Marshal into JSON
	out, err := json.Marshal(j)
	if err != nil {
		return nil, err
	}

	return out, nil
}

// FileRecordRepository is used to contain methods to load multiple FileRecord structs
type FileRecordRepository struct {
}

// Delete FileRecord from storage
func (f FileRecord) Delete() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Delete FileRecord
	if err = db.DeleteFileRecord(f.InfoHash, "info_hash"); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Save FileRecord to storage
func (f FileRecord) Save() error {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return err
	}

	// Save FileRecord
	if err := db.SaveFileRecord(f); err != nil {
		return err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return err
	}

	return nil
}

// Load FileRecord from storage
func (f FileRecord) Load(id interface{}, col string) (FileRecord, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return FileRecord{}, err
	}

	// Load FileRecord by column
	if f, err = db.LoadFileRecord(id, col); err != nil {
		return FileRecord{}, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return FileRecord{}, err
	}

	return f, nil
}

// CompactPeerList returns a packed byte array of peers who are active on this file
func (f FileRecord) CompactPeerList(numwant int, http bool) ([]byte, error) {
	// Retrieve list of peers
	peers, err := f.PeerList(numwant, http)
	if err != nil {
		return nil, err
	}

	// Create buffer to store compact peers
	compactPeers := make([]byte, 0)

	// Iterate peers
	for _, peer := range peers {
		// Marshal each peer to binary
		peerBuf, err := peer.MarshalBinary()
		if err != nil {
			return nil, err
		}

		// Append peer to compact list
		compactPeers = append(compactPeers[:], peerBuf...)
	}

	// Return compact peer list
	return compactPeers, nil
}

// Completed returns the number of completions, active or not, on this file
func (f FileRecord) Completed() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve number of file completions
	completed, err := db.CountFileRecordCompleted(f.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return completed, nil
}

// Seeders returns the number of seeders on this file
func (f FileRecord) Seeders() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Return number of active seeders
	seeders, err := db.CountFileRecordSeeders(f.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return seeders, nil
}

// Leechers returns the number of leechers on this file
func (f FileRecord) Leechers() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Return number of active leechers
	leechers, err := db.CountFileRecordLeechers(f.ID)
	if err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return leechers, nil
}

// PeerList returns a list of peers on this torrent, for tracker announce
func (f FileRecord) PeerList(numwant int, http bool) ([]Peer, error) {
	// List of peers
	peers := make([]Peer, 0)

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return peers, err
	}

	// Return list of peers, up to numwant
	if peers, err = db.GetFileRecordPeerList(f.InfoHash, numwant, http); err != nil {
		return peers, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return peers, err
	}

	return peers, nil
}

// PeerReaper reaps peers who have not recently announced on this torrent, and mark them inactive
func (f FileRecord) PeerReaper() (int, error) {
	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return 0, err
	}

	// Retrieve list of inactive users (have not announced in just above maximum interval)
	users, err := db.GetInactiveUserInfo(f.ID, time.Duration(int64(common.Static.Config.Interval))*time.Second+60)
	if err != nil {
		return 0, err
	}

	// Mark those users inactive on this file
	if err := db.MarkFileUsersInactive(f.ID, users); err != nil {
		return 0, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return 0, err
	}

	return len(users), nil
}

// Users loads all FileUserRecord structs associated with this FileRecord struct
func (f FileRecord) Users() []FileUserRecord {
	return new(FileUserRecordRepository).Select(f.ID, "file_id")
}

// All loads all FileRecord structs from storage
func (f FileRecordRepository) All() ([]FileRecord, error) {
	files := make([]FileRecord, 0)

	// Open database connection
	db, err := DBConnect()
	if err != nil {
		return files, err
	}

	// Retrieve all files
	files, err = db.GetAllFileRecords()
	if err != nil {
		return files, err
	}

	// Close database connection
	if err := db.Close(); err != nil {
		return files, err
	}

	return files, nil
}
