package goat

// HTTPStats represents statistics regarding HTTP server
type HTTPStats struct {
	Current int64
	Total   int64
}

// "Static" struct containing values which should be shared globally
var Static struct {
	// Trigger a graceful shutdown
	ShutdownChan chan bool

	// Configuration object
	Config Conf

	// Stats about HTTP server
	HTTP HTTPStats
}
