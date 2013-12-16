package goat

import (
	"time"
)

const APP = "goat"

func Manager(killChan chan bool, exitChan chan int) {
	// Set up logger
	logChan := make(chan string)
	Static.LogChan = logChan
	go LogManager()

	// Print startup status banner
	go PrintStatusBanner()

	// Set up database manager
	requestChan := make(chan Request, 100)
	Static.RequestChan = requestChan
	go DbManager()

	// Load configuration
	Static.Config = LoadConfig()

	// Set up graceful shutdown channel
	shutdownChan := make(chan bool)
	Static.ShutdownChan = shutdownChan

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
			// Trigger a graceful shutdown
			Static.LogChan <- "triggering graceful shutdown, press Ctrl+C again to force halt"
			Static.ShutdownChan <- true

			// Allow 1 second for graceful shutdown of all goroutines
			time.Sleep(1 * time.Second)

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
