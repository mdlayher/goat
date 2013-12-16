package goat

import (
	"crypto/sha1"
	"fmt"
)

// Generates a unique hash for a given struct with its ID, to be used as ID for storage
type Hashable interface {
	Hash() string
}

type PersistentWriter interface {
	PersistentWrite()
}

// Struct representing an announce, to be logged to storage
type AnnounceLog struct {
	InfoHash   string
	PeerId     string
	Ip         string
	Port       int
	Uploaded   int
	Downloaded int
	Left       int
	Event      string
	Time       int64
}

// Generate a SHA1 hash of the form: announce_log_InfoHash
func (log AnnounceLog) Hash() string {
	hash := sha1.New()
	hash.Write([]byte("announce_log"+log.InfoHash))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

//sends the annoucelog to the persistant write handler
func (log AnnounceLog) PersistentWrite() {
	re := make(chan Response)
	var request Request
	request.Data = log
	request.Id = log.InfoHash
	request.ResponseChan = re
	Static.PersistentChan <- request
	<-re
	close(re)
}

// Struct representing a file tracked by tracker
type FileRecord struct {
	InfoHash   string
	Leechers   int
	Seeders    int
	Completed  int
	Time       int64
	CreateTime int64
}

// Struct representing a scrapelog, to be logged to storage
type ScrapeLog struct {
	InfoHash string
	UserId   string
	Ip       string
	Time     int64
}

func (log ScrapeLog) PersistentWrite() {
	re := make(chan Response)
	var request Request
	request.Data = log
	request.Id = log.InfoHash
	request.ResponseChan = re
	Static.PersistentChan <- request
	<-re
	close(re)
}
