package goat

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/mdlayher/goat/goat/common"
	"github.com/mdlayher/goat/goat/data"
)

// Application name
const App = "goat"

// Application version
const Version = "git-master"

// Manager is responsible for coordinating the application
func Manager(killChan chan bool, exitChan chan int) {
	// Capture startup time
	common.Static.StartTime = time.Now().Unix()

	// Set up logging flags
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Println("Starting " + App + " " + Version)

	// Grab initial server status
	stat, err := common.GetServerStatus()
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Printf("%s - %s_%s (%d CPU) [pid: %d]", stat.Hostname, stat.Platform, stat.Architecture, stat.NumCPU, stat.PID)
	}

	// Load configuration
	config, err := common.LoadConfig()
	if config == (common.Conf{}) || err != nil {
		log.Println(err.Error())
		panic("Cannot load configuration, panicking")
	}
	common.Static.Config = config

	// Check for sane announce interval (10 minutes or more)
	if common.Static.Config.Interval <= 600 {
		panic("Announce interval must be at least 600 seconds, panicking")
	}

	// Attempt database connection
	if !data.DBPing() {
		panic(fmt.Errorf("cannot connect to database %s; panicking", data.DBName()))
	}
	log.Println("Database", data.DBName(), ": OK")

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
	if common.Static.Config.HTTP {
		go listenHTTP(httpSendChan, httpRecvChan)
		log.Println("HTTP listener launched on port " + strconv.Itoa(common.Static.Config.Port))
	}
	if common.Static.Config.SSL.Enabled {
		go listenHTTPS(httpsSendChan, httpsRecvChan)
		log.Println("HTTPS listener launched on port " + strconv.Itoa(common.Static.Config.SSL.Port))
	}
	if common.Static.Config.UDP {
		go listenUDP(udpSendChan, udpRecvChan)
		log.Println("UDP listener launched on port " + strconv.Itoa(common.Static.Config.Port))
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
			if common.Static.Config.HTTP {
				log.Println("Stopping HTTP listener")
				httpSendChan <- true
				<-httpRecvChan
			}
			if common.Static.Config.SSL.Enabled {
				log.Println("Stopping HTTPS listener")
				httpsSendChan <- true
				<-httpsRecvChan
			}
			if common.Static.Config.UDP {
				log.Println("Stopping UDP listener")
				udpSendChan <- true
				<-udpRecvChan
			}

			log.Println("Closing database:", data.DBName())
			data.DBCloseFunc()

			// Report that program should exit gracefully
			exitChan <- 0
		}
	}
}
