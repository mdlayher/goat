package goat

// Persistent defines methods to save, load, and delete data for an object
type Persistent interface {
	Save() bool
	Load(interface{}, string) interface{}
	Delete() bool
}
