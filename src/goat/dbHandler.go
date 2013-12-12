package goat

//holds information for request from database
type Request struct {
	Id           string
	ResponseChan chan Response
}

//holds information for response from databaseS
type Response struct {
	Id   string
	Data interface{}
}

// DbManager interface method ManageDb defines a database Manager which Manages requests
type DbManager interface {
	ManageDb(logChan chan string)
}

// MapDb is a key value storage database
// Id will be an identification for sharding
type MapDb struct {
	DbId string
	//this will change once we have a data structure, but for now is a generic interface
	Db map[string]interface{}
}

// Manage data MapDb requests
func (m MapDb) ManageDb(logChan chan string, requestChan chan Request) {
	var shards [100]MapDb

	for {
		select {
		case hold <- requestChan:
			go IoWorker(Request, shards[Request.Id[0:2]])
		}
	}


}

// SqlDb is a sql based database
type SqlDb struct {
}

// Manage Sql based requests
func (s SqlDb) ManageDb(logChan chan string) {

}
