package goat

// "Static" struct containing values which should be shared globally
var Static struct {
	// Server logging
	LogChan chan string

	// Server configuration
	Config struct {
		Port    string
		Passkey bool
		Http    bool
		Udp     bool
		Map     bool
		Sql     bool
	}

	// HTTP requests in last stat period, and total overall
	Http struct {
		Current int
		Total   int
	}
}
