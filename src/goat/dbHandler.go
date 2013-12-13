package goat

// Holds information for request from database
type Request struct {
	Query        string
	ResponseChan chan Response
}

// Holds information for response from database
type Response struct {
	Data, Id string
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	HandleDb()
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id string
	Db map[string]interface{}
}

// Handle data MapDb requests
func (m MapDb) HandleDb() {

}

// SqlDb is a sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb() {

}
