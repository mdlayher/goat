package common

// static is a struct containing values which should be shared globally
var Static struct {
	// Stats about API server
	API TimedStats

	// Configuration object
	Config Conf

	// Stats about HTTP server
	HTTP TimedStats

	// Maintenance
	Maintenance bool

	// Startup time
	StartTime int64

	// Status message
	StatusMessage string

	// Stats about UDP server
	UDP TimedStats
}

// TimedStats represents statistics over a period of time for a given listener
type TimedStats struct {
	Minute   int64 `json:"minute"`
	HalfHour int64 `json:"halfHour"`
	Hour     int64 `json:"hour"`
	Total    int64 `json:"total"`
}
