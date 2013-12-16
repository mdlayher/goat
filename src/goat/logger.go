package goat

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"
)

func LogManager() {
	// Create log directory and file, and pull current date to add to logfile name
	now := time.Now()
	os.Mkdir("logs", os.ModeDir|os.ModePerm)
	logFile, err := os.Create(fmt.Sprintf("logs/goat-%d-%d-%d.log", now.Year(), now.Month(), now.Day()))
	if err != nil {
		fmt.Println(err)
	}

	// create a logger that will use the writer created above
	logger := log.New(bufio.NewWriter(logFile), "", log.Lmicroseconds|log.Lshortfile)

	// Start the system status logger
	go StatusLogger()

	// Wait for error to be passed on the logChan channel, or termination signal
	for {
		select {
		case msg := <-Static.LogChan:
			now := time.Now()
			log := fmt.Sprintf("%s : [%4d-%02d-%02d %02d:%02d:%02d] %s\n", APP, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), msg)
			logger.Print(log)
			fmt.Print(log)
		}
	}
}

// Log and display system status every minute
func StatusLogger() {
	// Poll status every 30 seconds
	ticker := time.NewTicker(30 * time.Second)

	// Loop infinitely, trigger events via ticker
	for {
		select {
		case <-ticker.C:
			// Fetch status, log it
			go PrintCurrentStatus()
		}
	}
}
