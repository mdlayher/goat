package goat

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int) {
	// Set up graceful shutdown channel
	shutdownChan := make(chan bool)
	Static.ShutdownChan = shutdownChan

	// Set up logger
	logChan := make(chan string)
	Static.LogChan = logChan
	go LogManager()

	// Print startup status banner
	go PrintStatusBanner()

	// Set up database manager
	dbDoneChan := make(chan bool)
	requestChan := make(chan Request, 100)
	Static.RequestChan = requestChan
	go DbManager(dbDoneChan)

	// Load configuration
	Static.Config = LoadConfig()

	// Set up graceful shutdown channels
	httpDoneChan := make(chan bool)
	udpDoneChan := make(chan bool)

	// Launch listeners as configured
	if Static.Config.Http {
		go new(HttpListener).Listen(httpDoneChan)
		Static.LogChan <- "HTTP listener launched on port " + Static.Config.Port
	}
	if Static.Config.Udp {
		go new(UdpListener).Listen(udpDoneChan)
		Static.LogChan <- "UDP listener launched on port " + Static.Config.Port
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			// Trigger a graceful shutdown
			Static.LogChan <- "triggering graceful shutdown, press Ctrl+C again to force halt"
			Static.ShutdownChan <- true

			// Stop listeners
			if Static.Config.Http {
				Static.LogChan <- "stopping HTTP listener"
				// TODO: figure out why this keeps hanging
				// <-httpDoneChan
			}
			if Static.Config.Udp {
				Static.LogChan <- "stopping UDP listener"
				<-udpDoneChan
			}

			// Stop storage manager
			Static.LogChan <- "stopping database manager"
			<-dbDoneChan

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
