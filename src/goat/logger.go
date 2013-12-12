package goat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogManager(doneChan chan bool, logChan chan string) {
	// Create log directory and file, and pull current date to add to logfile name
	now := time.Now()
	os.Mkdir("logs", os.ModeDir|os.ModePerm)
	logFile, err := os.Create(fmt.Sprintf("logs/goat-%d-%d-%d.log", now.Year(), now.Month(), now.Day()))
	if err != nil {
		fmt.Println(err)
	}

	// create a logger that will use the writer created above
	logger := log.New(bufio.NewWriter(logFile), "", log.Lmicroseconds|log.Lshortfile)
	amIDone := false
	msg := ""

	// Start the system status logger
	go StatusLogger(logChan)

	// Wait for error to be passed on the logChan channel, or termination signal on doneChan
	for !amIDone {
		select {
		case amIDone = <-doneChan:
			logFile.Close()
		case msg = <-logChan:
			now := time.Now()
			logger.Printf("%s : [%2d:%2d:%2d] %s\n", APP, now.Hour(), now.Minute(), now.Second(), msg)
			fmt.Printf("%s : [%2d:%2d:%2d] %s\n", APP, now.Hour(), now.Minute(), now.Second(), msg)
		}
	}
}

// Log and display system status every minute
func StatusLogger(logChan chan string) {
	// Poll status every 30 seconds
	ticker := time.NewTicker(30 * time.Second)

	// Loop infinitely, trigger events via ticker
	for {
		select {
		case <- ticker.C:
			// Fetch status, log it
			go PrintCurrentStatus(logChan)
		}
	}
}
