package goat

type Request struct {
	Query        string
	ResponseChan chan Response
}
type Response struct {
	Data, Id string
}

type DbHandler interface {
	HandleDb(logChan chan string)
}

type MapDb struct {
}

func (m MapDb) HandleDb(logChan chan string) {

}

type SqlDb struct {
}

func (s SqlDb) HandleDb(logChan chan string) {

}
