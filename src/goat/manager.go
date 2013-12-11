package goat

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int, config Config) {
	// Set up logger
	logChan := make(chan string)
	doneChan := make(chan bool)
	go LogMng(doneChan, logChan)

	// Launch listeners as configured
	if config.Http {
		go new(HttpListener).Listen(config.Port, logChan)
		logChan <- "HTTP listener launched on port " + config.Port
	}
	if config.Udp {
		go new(UdpListener).Listen(config.Port, logChan)
		logChan <- "UDP listener launched on port " + config.Port
	}

	for {
		select {
		case <-killChan:
			// Change this to kill workers gracefully and exit
			exitChan <- 0
		}
	}
}
