package goat

// "Static" struct containing values which should be shared globally
var Static struct {
	// Trigger a graceful shutdown
	ShutdownChan chan bool

	// Configuration object
	Config Conf

	// Stats about HTTP server
	HTTP HTTPStats

	// Stats about UDP server
	UDP UDPStats
}
