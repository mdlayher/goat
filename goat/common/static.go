package common

// static is a struct containing values which should be shared globally
var Static struct {
	// Configuration object
	Config Conf

	// Stats about HTTP server
	HTTP HTTPStats

	// Maintenance
	Maintenance bool

	// Startup time
	StartTime int64

	// Status message
	StatusMessage string

	// Stats about UDP server
	UDP UDPStats
}

// HTTPStats represents statistics regarding HTTP server
type HTTPStats struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

// UDPStats represents statistics regarding UDP server
type UDPStats struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}
