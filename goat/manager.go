package goat

import (
	"log"
	"net/http"
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

	// Set up logging flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting " + App + " " + Version)

	// Grab initial server status
	stat := getServerStatus()
	if stat == (serverStatus{}) {
		log.Println("Could not print get startup status")
	} else {
		log.Printf("%s - %s_%s (%d CPU) [pid: %d]", stat.Hostname, stat.Platform, stat.Architecture, stat.NumCPU, stat.PID)
	}

	// Load configuration
	config := loadConfig()
	if config == (conf{}) {
		log.Println("Cannot load configuration, exiting now.")
		os.Exit(1)
	}
	static.Config = config

	// Check for sane announce interval (10 minutes or more)
	if static.Config.Interval <= 600 {
		log.Println("Announce interval must be at least 600 seconds.")
		os.Exit(1)
	}

	// Attempt database connection
	if !dbPing() {
		log.Println("Cannot connect to database", dbName(), ", exiting now.")
		os.Exit(1)
	}
	log.Println("Database", dbName(), ": OK")

	// If configured, attempt redis connection
	if static.Config.Redis {
		if !redisPing() {
			log.Println("Cannot connect to Redis, exiting now.")
			os.Exit(1)
		}
		log.Println("Redis : OK")
	}

	// Start cron manager
	go cronManager()

	// Set up graceful shutdown channels
	httpDoneChan := make(chan bool)
	httpsDoneChan := make(chan bool)
	udpDoneChan := make(chan bool)

	// Set up HTTP(S) route
	http.HandleFunc("/", parseHTTP)

	// Launch listeners as configured
	if static.Config.HTTP {
		go listenHTTP(httpDoneChan)
		log.Println("HTTP listener launched on port " + strconv.Itoa(static.Config.Port))
	}
	if static.Config.HTTPS {
		go listenHTTPS(httpsDoneChan)
		log.Println("HTTPS listener launched on port " + strconv.Itoa(static.Config.SSL.Port))
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
			log.Println("Triggering graceful shutdown, press Ctrl+C again to force halt")

			// Stop listeners
			if static.Config.HTTP {
				log.Println("Stopping HTTP listener")
				static.ShutdownChan <- true
				<-httpDoneChan
			}
			if static.Config.HTTPS {
				log.Println("Stopping HTTPS listener")
				static.ShutdownChan <- true
				<-httpsDoneChan
			}
			if static.Config.UDP {
				log.Println("Stopping UDP listener")
				static.ShutdownChan <- true
				<-udpDoneChan
			}

			log.Println("Closing database connection")
			dbCloseFunc()

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
