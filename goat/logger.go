package goat

import (
	"fmt"
	"time"
)

// LogManager is responsible for creating the main loggers
func LogManager() {
	// Start the system status logger
	go StatusLogger()

	// Wait for error to be passed on the logChan channel, or termination signal
	for {
		select {
		case msg := <-Static.LogChan:
			now := time.Now()
			out := fmt.Sprintf("%s : [%4d-%02d-%02d %02d:%02d:%02d] %s\n", App, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), msg)
			fmt.Print(out)
		}
	}
}

// StatusLogger logs and displays system status at regular intervals
func StatusLogger() {
	ticker := time.NewTicker(5 * time.Minute)

	// Loop infinitely, trigger events via ticker
	for {
		select {
		case <-ticker.C:
			// Fetch status, log it
			go PrintCurrentStatus()
		}
	}
}
