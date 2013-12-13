package goat

func DbManager(RequestChan chan Request) {
	if Static.Config.Map {
		go new(MapDb).HandleDb()
		Static.LogChan <- "MapDb instance launched"
	}
	if Static.Config.Sql {
		go new(SqlDb).HandleDb()
		Static.LogChan <- "SqlDb instance launched"
	}
	for {

	}
}

// Holds information for request from database
type Request struct {
	Id           string
	Read, Write  bool
	Data         interface{}
	ResponseChan chan Response
}

// Holds information for response from database
type Response struct {
	Id, Db string
	Data   interface{}
}

// DbHandler interface method HandleDb defines a database handler which handles requests
type DbHandler interface {
	HandleDb(RequestChan chan Request)
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	Id string
	Db map[string]map[string]interface{}
}

// Handle data MapDb requests
func (m MapDb) HandleDb(RequestChan chan Request) {

}

// SqlDb is a sql based database
type SqlDb struct {
}

// Handle Sql based requests
func (s SqlDb) HandleDb(RequestChan chan Request) {

}
