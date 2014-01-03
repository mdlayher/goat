package goat

// Server configuration
type Conf struct {
	Port      int
	Passkey   bool
	Whitelist bool
	Interval  int
	Http      bool
	Udp       bool
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

	// Configuration object
	Config Conf

	// Stats about HTTP server
	Http HttpStats
}
