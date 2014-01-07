package goat

import (
	"os"
	"strconv"
)

// Application name
const App = "goat"

// Application version
const Version = "git-master"

// Manager is responsible for coordinating the application
func Manager(killChan chan bool, exitChan chan int) {
	// Set up graceful shutdown channel
	shutdownChan := make(chan bool)
	Static.ShutdownChan = shutdownChan

	// Set up logger
	logChan := make(chan string)
	Static.LogChan = logChan
	go LogManager()

	Static.LogChan <- "Starting " + App + " " + Version

	// Print startup status banner
	go PrintStatusBanner()

	// Load configuration
	config := LoadConfig()
	if config == (Conf{}) {
		Static.LogChan <- "Cannot load configuration, exiting now."
		os.Exit(1)
	}
	Static.Config = config

	// Set up graceful shutdown channels
	httpDoneChan := make(chan bool)
	udpDoneChan := make(chan bool)

	// Launch listeners as configured
	if Static.Config.HTTP {
		go new(HTTPListener).Listen(httpDoneChan)
		Static.LogChan <- "HTTP listener launched on port " + strconv.Itoa(Static.Config.Port)
	}
	if Static.Config.UDP {
		go new(UDPListener).Listen(udpDoneChan)
		Static.LogChan <- "UDP listener launched on port " + strconv.Itoa(Static.Config.Port)
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			// Trigger a graceful shutdown
			Static.LogChan <- "triggering graceful shutdown, press Ctrl+C again to force halt"
			Static.ShutdownChan <- true

			// Stop listeners
			if Static.Config.HTTP {
				Static.LogChan <- "stopping HTTP listener"
				//<-httpDoneChan
			}
			if Static.Config.UDP {
				Static.LogChan <- "stopping UDP listener"
				//<-udpDoneChan
			}

			// TODO Is this the right place?
			Static.LogChan <- "calling dbCloseFunc"
			dbCloseFunc()

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
