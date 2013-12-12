package goat

//holds information for request from database
type Request struct {
	Query        string
	ResponseChan chan Response
}

//holds information for response from databaseS
type Response struct {
	Data, Id string
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	HandleDb(logChan chan string)
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id int
	Db map[string]string
}

// Handle data MapDb requests
func (m MapDb) HandleDb(logChan chan string) {

}

// SqlDb is a sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb(logChan chan string) {

}
