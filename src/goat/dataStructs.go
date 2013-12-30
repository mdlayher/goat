package goat

// Generates a unique hash for a given struct with its ID, to be used as ID for storage
type Hashable interface {
	Hash() string
}

// Defines methods to save, load, and delete data for an object
type Persistent interface {
	Save() bool
	Load(interface{}, string) interface{}
	Delete() bool
}

type ErrorRes struct {
	ErrLocation string
	Time        int64
	Error       string
}

// Holds information for request from database
type Request struct {
	Id           string
	Data         interface{}
	ResponseChan chan Response
}

// Holds information for response from database
type Response struct {
	Id   string
	Db   string
	Data interface{}
}

// Specifies pointers for a given relationship of data
type Relation struct {
	Index []*interface{}
}

type MapStorage map[string]interface{}
