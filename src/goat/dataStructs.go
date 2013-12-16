package goat

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

type ScrapeLogs struct {
}
