package goat

// Server configuration
type Conf struct {
	Port      string
	Passkey   bool
	Whitelist bool
	Http      bool
	Udp       bool
	Map       bool
	Sql       bool
	// Must be an int between 1-8
	// Higher values result in faster performance and larger memory requirements
	CacheSize int
}

// Statistics regarding HTTP server
type HttpStats struct {
	Current int64
	Total   int64
}

// "Static" struct containing values which should be shared globally
var Static struct {
	// Server logging
	LogChan chan string

	// Trigger a graceful shutdown
	ShutdownChan chan bool

	// Buffered request channel for DB manager
	RequestChan chan Request

	// Backup request channel for persistant storage
	PersistentChan chan Request

	// Error handling chan for various errors in goroutines
	ErrChan chan Response

	// Configuration object
	Config Conf

	// Stats about HTTP server
	Http HttpStats
}
