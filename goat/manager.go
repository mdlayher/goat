package goat

import (
	"log"
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
	go LogManager()

	log.Println("Starting " + App + " " + Version)

	// Print startup status banner
	go PrintStatusBanner()

	// Load configuration
	config := LoadConfig()
	if config == (Conf{}) {
		log.Println("Cannot load configuration, exiting now.")
		os.Exit(1)
	}
	Static.Config = config

	// Attempt database connection
	if !DBPing() {
		log.Println("Cannot connect to MySQL database, exiting now.")
		os.Exit(1)
	}

	// Set up graceful shutdown channels
	httpDoneChan := make(chan bool)
	udpDoneChan := make(chan bool)

	// Launch listeners as configured
	if Static.Config.HTTP {
		go new(HTTPListener).Listen(httpDoneChan)
		log.Println("HTTP listener launched on port " + strconv.Itoa(Static.Config.Port))
	}
	if Static.Config.UDP {
		go new(UDPListener).Listen(udpDoneChan)
		log.Println("UDP listener launched on port " + strconv.Itoa(Static.Config.Port))
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			// Trigger a graceful shutdown
			log.Println("triggering graceful shutdown, press Ctrl+C again to force halt")
			Static.ShutdownChan <- true

			// Stop listeners
			if Static.Config.HTTP {
				log.Println("stopping HTTP listener")
				//<-httpDoneChan
			}
			if Static.Config.UDP {
				log.Println("stopping UDP listener")
				//<-udpDoneChan
			}

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
