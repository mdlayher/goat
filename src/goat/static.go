package goat

// Server configuration
type Conf struct {
	Port    string
	Passkey bool
	Http    bool
	Udp     bool
	Map     bool
	Sql     bool
	// Must be an int between 1-8
	// Higher values result in faster performance and larger memory requirements
	CacheSize int
}

// Statistics regarding HTTP server
type HttpStats struct {
	Current int
	Total   int
}

// "Static" struct containing values which should be shared globally
var Static struct {
	// Server logging
	LogChan     chan string
	RequestChan chan Request

	// Configuration object
	Config Conf
	Http   HttpStats
}
