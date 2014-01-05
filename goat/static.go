package goat

// HTTPStats represents statistics regarding HTTP server
type HTTPStats struct {
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
	HTTP HTTPStats
}
