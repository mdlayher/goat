package goat

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int) {
	// Set up logger
	logChan := make(chan string)
	doneChan := make(chan bool)
	go LogManager(doneChan, logChan)

	// Print startup status banner
	go PrintStatusBanner(logChan)

	// Load configuration
	config := LoadConfig()

	// Launch listeners as configured
	if config.Http {
		go new(HttpListener).Listen(config.Port, logChan)
		logChan <- "HTTP listener launched on port " + config.Port
	}
	if config.Udp {
		go new(UdpListener).Listen(config.Port, logChan)
		logChan <- "UDP listener launched on port " + config.Port
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			exitChan <- 0
		}
	}
}
