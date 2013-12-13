package goat

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int) {
	// Set up logger
	logChan := make(chan string)
	Static.LogChan = logChan
	doneChan := make(chan bool)
	go LogManager(doneChan)

	// Print startup status banner
	go PrintStatusBanner()

	// Load configuration
	config := LoadConfig()
	Static.Config.Port = config.Port
	Static.Config.Passkey = config.Passkey
	Static.Config.Http = config.Http
	Static.Config.Udp = config.Udp
	Static.Config.Map = config.Map
	Static.Config.Sql = config.Sql

	// Launch listeners as configured
	if Static.Config.Http {
		go new(HttpListener).Listen()
		Static.LogChan <- "HTTP listener launched on port " + Static.Config.Port
	}
	if Static.Config.Udp {
		go new(UdpListener).Listen()
		Static.LogChan <- "UDP listener launched on port " + Static.Config.Port
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			exitChan <- 0
		}
	}
}
