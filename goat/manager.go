package goat

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"
)

// Application name
const App = "goat"

// Application version
const Version = "git-master"

// Manager is responsible for coordinating the application
func Manager(killChan chan bool, exitChan chan int) {
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
		panic(fmt.Errorf("Cannot connect to database %s; exiting now", dbName()))
	}
	log.Println("Database", dbName(), ": OK")

	db, err := dbConnectFunc()
	if err != nil {
	}
	for _, schema := range mysql_schemas {
		_ = (db).(*dbw).Execf(schema)
		if err != nil {
			panic(err)
		}
	}

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
	httpSendChan := make(chan bool)
	httpRecvChan := make(chan bool)
	httpsSendChan := make(chan bool)
	httpsRecvChan := make(chan bool)
	udpSendChan := make(chan bool)
	udpRecvChan := make(chan bool)

	// Set up HTTP(S) route
	http.HandleFunc("/", parseHTTP)

	// Launch listeners as configured
	if static.Config.HTTP {
		go listenHTTP(httpSendChan, httpRecvChan)
		log.Println("HTTP listener launched on port " + strconv.Itoa(static.Config.Port))
	}
	if static.Config.HTTPS {
		go listenHTTPS(httpsSendChan, httpsRecvChan)
		log.Println("HTTPS listener launched on port " + strconv.Itoa(static.Config.SSL.Port))
	}
	if static.Config.UDP {
		go listenUDP(udpSendChan, udpRecvChan)
		log.Println("UDP listener launched on port " + strconv.Itoa(static.Config.Port))
	}

	// Wait for shutdown signal
	for {
		select {
		case <-killChan:
			// Trigger a graceful shutdown
			log.Println("Triggering graceful shutdown, press Ctrl+C again to force halt")

			// If program hangs for more than 10 seconds, trigger a force halt
			go func() {
				<-time.After(10 * time.Second)
				log.Println("Timeout reached, triggering force halt")
				if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
					log.Println(err.Error())
				}
			}()

			// Stop listeners
			if static.Config.HTTP {
				log.Println("Stopping HTTP listener")
				httpSendChan <- true
				<-httpRecvChan
			}
			if static.Config.HTTPS {
				log.Println("Stopping HTTPS listener")
				httpsSendChan <- true
				<-httpsRecvChan
			}
			if static.Config.UDP {
				log.Println("Stopping UDP listener")
				udpSendChan <- true
				<-udpRecvChan
			}

			log.Println("Closing database connection")
			dbCloseFunc()

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
