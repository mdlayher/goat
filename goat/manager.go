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
	static.ShutdownChan = shutdownChan

	// Set up logger
	go logManager()

	log.Println("Starting " + App + " " + Version)

	// Print startup status banner
	go printStatusBanner()

	// Load configuration
	config := loadConfig()
	if config == (conf{}) {
		log.Println("Cannot load configuration, exiting now.")
		os.Exit(1)
	}
	static.Config = config

	// Attempt database connection
	if !dbPing() {
		log.Println("Cannot connect to MySQL database, exiting now.")
		os.Exit(1)
	}
	log.Println("MySQL: OK")

	// Attempt redis connection
	if !redisPing() {
		log.Println("Cannot connect to Redis, exiting now.")
		os.Exit(1)
	}
	log.Println("Redis: OK")

	// Set up graceful shutdown channels
	httpDoneChan := make(chan bool)
	udpDoneChan := make(chan bool)

	// Launch listeners as configured
	if static.Config.HTTP {
		go listenHTTP(httpDoneChan)
		log.Println("HTTP listener launched on port " + strconv.Itoa(static.Config.Port))
	}
	if static.Config.UDP {
		go listenUDP(udpDoneChan)
		log.Println("UDP listener launched on port " + strconv.Itoa(static.Config.Port))
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			// Trigger a graceful shutdown
			log.Println("triggering graceful shutdown, press Ctrl+C again to force halt")

			// Stop listeners
			if static.Config.HTTP {
				log.Println("stopping HTTP listener")
				static.ShutdownChan <- true
				<-httpDoneChan
			}
			if static.Config.UDP {
				log.Println("stopping UDP listener")
				static.ShutdownChan <- true
				<-udpDoneChan
			}

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
