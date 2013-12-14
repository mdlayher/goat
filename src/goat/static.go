package goat

// "Static" struct containing values which should be shared globally
var Static struct {
	// Server logging
	LogChan     chan string
	RequestChan chan Request

	// Server configuration
	Config struct {
		Port      string
		Passkey   bool
		Http      bool
		Udp       bool
		Map       bool
		CacheSize int //must be an int between 1-8. Higher values result in
		//faster performance and larger memory requirements
		Sql bool
	}

	// HTTP requests in last stat period, and total overall
	Http struct {
		Current int
		Total   int
	}
}
