package goat

// static is a struct containing values which should be shared globally
var static struct {
	// Trigger a graceful shutdown
	ShutdownChan chan bool

	// Configuration object
	Config conf

	// Stats about HTTP server
	HTTP httpStats

	// Stats about UDP server
	UDP udpStats
}
